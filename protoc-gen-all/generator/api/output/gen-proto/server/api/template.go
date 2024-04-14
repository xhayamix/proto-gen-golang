package api

import (
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed api.gen.proto.tpl
var templateFileBytes []byte

//go:embed message.tpl
var messageTemplateFileBytes []byte

type Creator struct{}

type MethodOption struct {
	Key   string
	Value string
}

type Method struct {
	PascalName string
	Comment    string
	InputType  string
	OutputType string
	Options    []*MethodOption
}

type Service struct {
	PascalName string
	Comment    string
	Methods    []*Method
}

type Field struct {
	PkgName   string
	CamelName string
	Comment   string
	Type      string
	Number    int32
	Option    string
}

type Message struct {
	PascalName        string
	Comment           string
	Messages          []*Message
	Fields            []*Field
	HasCommonResponse bool
}

type Data struct {
	IsCommon    bool
	Service     *Service
	Messages    []*Message
	ImportPaths []string
}

func createMessage(
	message *input.Message,
	responseMessageNameSetNeedCommonResponse *strset.Set,
	importPathSet *strset.Set,
) *Message {
	ret := &Message{
		PascalName:        core.ToPascalCase(message.SnakeName),
		Comment:           message.Comment,
		Fields:            make([]*Field, 0, len(message.Fields)),
		Messages:          make([]*Message, 0, len(message.Messages)),
		HasCommonResponse: responseMessageNameSetNeedCommonResponse.Has(message.SnakeName),
	}

	for _, m := range message.Messages {
		ret.Messages = append(ret.Messages, createMessage(m, responseMessageNameSetNeedCommonResponse, importPathSet))
	}

	for _, field := range message.Fields {
		typeName := field.Type
		switch field.TypeKind {
		case input.TypeKind_Enum:
			importPathSet.Add("client/enums/enums_gen.proto")
			typeName = "enums." + typeName
		case input.TypeKind_Map:
			typeName = fmt.Sprintf("map<%s, %s>", field.MapKeyType, field.MapValueType)
		}
		switch {
		case field.PkgType == input.PkgType_APICommon:
			if field.ImportFileName != "" {
				importPathSet.Add("client/api/common/" + strings.ReplaceAll(field.ImportFileName, ".proto", "_gen.proto"))
			}
			typeName = "api.common." + typeName
		case field.PkgType == input.PkgType_ServerCommon:
			if field.ImportFileName != "" {
				importPathSet.Add("client/common/" + strings.ReplaceAll(field.ImportFileName, ".proto", "_gen.proto"))
			}
			typeName = "client.common." + typeName
		case field.ImportFileName != "":
			importPathSet.Add("client/api/" + strings.ReplaceAll(field.ImportFileName, ".proto", "_gen.proto"))
		}
		if field.IsList {
			typeName = "repeated " + typeName
		}
		options := make([]string, 0)
		if field.FieldOption.MasterRef != nil {
			importPathSet.Add("client/master/common/options.proto")
			masterRef := field.FieldOption.MasterRef

			fkInfoList := make([]string, 0, 3)
			fkInfoList = append(fkInfoList,
				core.ToGolangPascalCase(masterRef.TableSnakeName),
				core.ToCamelCase(masterRef.ColumnSnakeName),
			)
			if len(masterRef.ParentFieldSnakeNames) > 0 {
				fkInfoList = append(fkInfoList, core.ToCamelCase(masterRef.ParentFieldSnakeNames[0]))
			}

			options = append(options, fmt.Sprintf("(master.common.fk) = %q", strings.Join(fkInfoList, ":")))
		}
		if field.ValidateOption != "" {
			importPathSet.Add("validate/validate.proto")
			options = append(options, "(validate.rules) = { "+field.ValidateOption+" }")
		}
		if field.HiddenOption {
			importPathSet.Add("client/options/zap.proto")
			options = append(options, "(options.zap.hidden) = true")
		}

		var option string
		if len(options) == 0 {
			option = ""
		} else {
			option = " [" + strings.Join(options, ",") + "]"
		}

		ret.Fields = append(ret.Fields, &Field{
			CamelName: core.ToCamelCase(field.SnakeName),
			Comment:   field.Comment,
			Type:      typeName,
			Number:    field.Number,
			Option:    option,
		})
	}

	return ret
}

func (c *Creator) Create(file *input.File) (*output.TemplateInfo, error) {
	if file.IsCommon && file.SnakeName == "user_data" {
		// user_dataはtransaction起因で作成する
		return nil, nil
	}
	data := &Data{
		IsCommon:    file.IsCommon,
		Service:     nil,
		Messages:    make([]*Message, 0, len(file.Messages)),
		ImportPaths: make([]string, 0),
	}
	importPathSet := strset.New()
	responseMessageNameSetNeedCommonResponse := strset.New()
	emptyResponses := make([]string, 0)

	if file.Service != nil {
		// grpc-gateway使う場合は追加
		// importPathSet.Add("google/api/annotations.proto")

		inputService := file.Service
		s := &Service{
			PascalName: core.ToPascalCase(inputService.SnakeName),
			Comment:    inputService.Comment,
			Methods:    make([]*Method, 0, len(inputService.Methods)),
		}
		data.Service = s
		for _, method := range inputService.Methods {
			if !method.DisableCommonResponse && strings.HasSuffix(method.OutputMessage.SnakeName, "_response") {
				responseMessageNameSetNeedCommonResponse.Add(method.OutputMessage.SnakeName)
			}

			methodName := core.ToPascalCase(method.SnakeName)

			var inputType string
			if method.InputMessage.SnakeName == "empty" {
				importPathSet.Add("google/protobuf/empty.proto")
				inputType = "google.protobuf.Empty"
			} else {
				inputType = core.ToPascalCase(method.InputMessage.SnakeName)
			}
			var outputType string
			if method.OutputMessage.SnakeName == "empty" {
				outputType = s.PascalName + methodName + "Response"
				emptyResponses = append(emptyResponses, outputType)
			} else {
				outputType = core.ToPascalCase(method.OutputMessage.SnakeName)
			}

			methodOptions := make([]*MethodOption, 0)
			if method.CheckOption != "<nil>" {
				methodOptions = append(methodOptions, &MethodOption{
					Key:   "options.check_option.checkOption",
					Value: "{" + method.CheckOption + "}",
				})
				importPathSet.Add("client/options/check_option.proto")
			}
			if method.ErrorOption != "" {
				methodOptions = append(methodOptions, &MethodOption{
					Key:   "options.check_option.errorOption",
					Value: "{" + method.ErrorOption + "}",
				})
				importPathSet.Add("client/options/check_option.proto")
			}

			s.Methods = append(s.Methods, &Method{
				PascalName: methodName,
				Comment:    method.Comment,
				InputType:  inputType,
				OutputType: outputType,
				Options:    methodOptions,
			})
		}
	}

	for _, message := range file.Messages {
		m := createMessage(message, responseMessageNameSetNeedCommonResponse, importPathSet)

		data.Messages = append(data.Messages, m)
	}

	for _, response := range emptyResponses {
		data.Messages = append(data.Messages, &Message{
			PascalName:        response,
			Comment:           response,
			HasCommonResponse: true,
		})
	}

	/*
		if len(emptyResponses) > 0 || responseMessageNameSetNeedCommonResponse.Size() > 0 {
			importPathSet.Add("client/api/common/response_gen.proto")
		}
	*/

	paths := importPathSet.List()
	sort.Strings(paths)
	data.ImportPaths = paths

	msgTpl, err := core.GetBaseTemplate().Parse(string(messageTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	tpl, err := msgTpl.Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	var pkgName string
	if file.IsCommon {
		pkgName = "common"
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("gen-proto/client/api", pkgName, file.SnakeName+"_gen.proto"),
	}, nil
}
