package handler

import (
	"bytes"
	_ "embed"
	"net/http"
	"text/template"

	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed handler.gen.go.tpl
var templateFileBytes []byte

type Data struct {
	PackageName string
	PascalName  string
	ImportPaths []string
	Methods     []*Method
}

type Method struct {
	PascalName          string
	Comment             string
	InputType           string
	InputTypeSuffix     string
	InputFields         []*Field
	OutputType          string
	OutputTypeSuffix    string
	OutputAny           bool
	OutputFields        []*Field
	DisableRequestBind  bool
	DisableOnProduction bool
	OnlyAdminUser       bool
	Method              string
	Path                string
}

type Field struct {
	PascalName string
	CamelName  string
	Comment    string
	Type       string
}

func New() output.EachTemplateCreator {
	tpl := template.Must(core.GetBaseTemplate().Parse(string(templateFileBytes)))
	return &creator{tpl: tpl}
}

type creator struct {
	tpl *template.Template
}

func (c *creator) Create(file *input.File) (*output.TemplateInfo, error) {
	data := &Data{
		PackageName: core.ToPkgName(file.SnakeName),
		PascalName:  core.ToPascalCase(file.SnakeName),
		Methods:     make([]*Method, 0, len(file.Service.Methods)),
	}
	importPathSet := strset.New()
	inputTypeSet := strset.New()
	for _, method := range file.Service.Methods {
		var inputFields []*Field
		inputType := core.ToPascalCase(method.InputMessage.SnakeName)
		var inputTypeSuffix string
		switch {
		case inputType == "Empty":
			inputType = "emptypb.Empty"
			importPathSet.Add(`emptypb "google.golang.org/protobuf/types/known/emptypb"`)
		case inputType == "Any":
			inputType = "anypb.Any"
			importPathSet.Add(`anypb "google.golang.org/protobuf/types/known/anypb"`)
		case method.HttpMethod == http.MethodGet:
			inputTypeSuffix = "Query"
			if !inputTypeSet.Has(inputType) {
				inputTypeSet.Add(inputType)
				for _, f := range method.InputMessage.Fields {
					typ := f.Type
					if f.IsEnum {
						typ = "enums." + typ
						importPathSet.Add(`enums "github.com/xhayamix/proto-gen-golang/pkg/domain/proto/server/enums"`)
					}
					if f.IsList {
						typ = "[]" + typ
					}
					inputFields = append(inputFields, &Field{
						PascalName: core.ToPascalCase(f.SnakeName),
						CamelName:  core.ToCamelCase(f.SnakeName),
						Comment:    f.Comment,
						Type:       typ,
					})
				}
			}
		}
		outputType := core.ToPascalCase(method.OutputMessage.SnakeName)
		var outputAny bool
		if outputType == "Empty" {
			outputType = "emptypb.Empty"
			importPathSet.Add(`emptypb "google.golang.org/protobuf/types/known/emptypb"`)
		} else if outputType == "Any" {
			outputType = "interface{}"
			outputAny = true
		}
		var outputSuffix string
		var outputFields []*Field
		for _, f := range method.OutputMessage.Fields {
			if f.Type == "Any" {
				outputSuffix = "Extend"
				outputFields = append(outputFields, &Field{
					PascalName: core.ToPascalCase(f.SnakeName),
					CamelName:  core.ToCamelCase(f.SnakeName),
					Comment:    f.Comment,
					Type:       "interface{}",
				})
			} else if f.TypeKind == input.TypeKind_Map && f.MapValueType == "Any" {
				outputSuffix = "Extend"
				outputFields = append(outputFields, &Field{
					PascalName: core.ToPascalCase(f.SnakeName),
					CamelName:  core.ToCamelCase(f.SnakeName),
					Comment:    f.Comment,
					Type:       "map[" + f.MapKeyType + "]interface{}",
				})
			}
		}
		path := method.HttpPath

		data.Methods = append(data.Methods, &Method{
			PascalName:       core.ToPascalCase(method.SnakeName),
			Comment:          method.Comment,
			InputType:        inputType,
			InputTypeSuffix:  inputTypeSuffix,
			InputFields:      inputFields,
			OutputType:       outputType,
			OutputTypeSuffix: outputSuffix,
			OutputAny:        outputAny,
			OutputFields:     outputFields,
			Method:           method.HttpMethod,
			Path:             path,
		})
	}
	data.ImportPaths = importPathSet.List()

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/cmd/admin/handler", data.PackageName, core.ToSnakeCase(data.PascalName)+".gen.go"),
	}, nil
}
