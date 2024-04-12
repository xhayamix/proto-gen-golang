package requests

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed request.gen.ts.tpl
var templateFileBytes []byte

type Data struct {
	Types        []*Type
	Methods      []*Method
	Commons      []string
	Enums        []string
	Transactions []string
}

type Type struct {
	Name    string
	Comment string
	Columns []*Column
}

type Column struct {
	Name    string
	Type    string
	Comment string
}

type Method struct {
	Name         string
	Comment      string
	RequestType  string
	ResponseType string
	HttpMethod   string
	HttpPath     string
}

func New() output.EachTemplateCreator {
	tpl := template.Must(core.GetBaseTemplate().Parse(string(templateFileBytes)))
	return &creator{tpl: tpl}
}

type creator struct {
	tpl *template.Template
}

func (c *creator) Create(file *input.File) (*output.TemplateInfo, error) {
	types := make([]*Type, 0, len(file.Messages))
	tMap := make(typeMap)

	commonSet := strset.New()
	enumSet := strset.New()
	transactionSet := strset.New()
	for _, msg := range file.Messages {
		types = append(types, c.convertMessage(msg, tMap, "", commonSet, enumSet, transactionSet)...)
	}
	methods := make([]*Method, 0, len(file.Service.Methods))
	for _, method := range file.Service.Methods {
		inputType := core.ToPascalCase(method.InputMessage.SnakeName)
		switch {
		case inputType == "Empty":
			inputType = ""
		case inputType == "Any":
			inputType = "unknown"
		case !tMap.Has(inputType):
			types = append(types, c.convertMessage(method.InputMessage, tMap, "", commonSet, enumSet, transactionSet)...)
		}

		outputType := core.ToPascalCase(method.OutputMessage.SnakeName)
		switch {
		case outputType == "Empty":
			outputType = "void"
		case outputType == "Any":
			outputType = "unknown"
		case !tMap.Has(outputType):
			types = append(types, c.convertMessage(method.OutputMessage, tMap, "", commonSet, enumSet, transactionSet)...)
		}

		var httpMethod string
		if method.HttpMethod == http.MethodDelete {
			httpMethod = "del"
		} else {
			httpMethod = strings.ToLower(method.HttpMethod)
		}

		name := core.ToCamelCase(method.SnakeName)
		// typescriptでdeleteは使えないので変更
		if name == "delete" {
			name = "del"
		}
		methods = append(methods, &Method{
			Name:         name,
			Comment:      method.Comment,
			RequestType:  inputType,
			ResponseType: outputType,
			HttpMethod:   httpMethod,
			HttpPath:     method.HttpPath,
		})
	}

	commons := commonSet.List()
	enums := enumSet.List()
	transactions := transactionSet.List()
	sort.Strings(commons)
	sort.Strings(enums)
	sort.Strings(transactions)

	data := &Data{
		Types:        types,
		Methods:      methods,
		Commons:      commons,
		Enums:        enums,
		Transactions: transactions,
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("web/src/apps/admin/requests", file.SnakeName+".gen.ts"),
	}, nil
}

func (c *creator) convertMessage(msg *input.Message, tMap typeMap, prefix string, commonSet, enumSet, transactionSet *strset.Set) []*Type {
	pascalName := core.ToPascalCase(msg.SnakeName)
	name := prefix + pascalName
	tMap.Add(pascalName, name)

	types := make([]*Type, 0, 1+len(msg.Messages))
	for _, m := range msg.Messages {
		types = append(types, c.convertMessage(m, tMap, name, commonSet, enumSet, transactionSet)...)
	}

	columns := make([]*Column, 0, len(msg.Fields))
	for _, field := range msg.Fields {
		typ := c.convertTsType(field.Type)
		comment := field.Comment
		switch {
		case field.TypeKind == input.TypeKind_Map:
			valueType := c.convertTsType(field.MapValueType)
			typ = fmt.Sprintf("{ [key: %s]: %s }", field.MapKeyType, valueType)
		case field.TypeKind == input.TypeKind_Message:
			if field.PkgType == input.PkgType_ClientCommon {
				commonSet.Add(field.Type)
			} else if field.PkgType == input.PkgType_ClientTransaction {
				transactionSet.Add(field.Type)
			} else if tMap.Has(field.Type) {
				typ = tMap.Get(field.Type)
			} else if strings.Contains(field.ParentGoPackage, ".") && !strings.HasPrefix(field.ParentGoPackage, "github.com/QualiArts/campus-server/pkg/cmd/admin/handler") {
				typ = "unknown"
				comment += " " + field.ParentGoPackage + "." + field.Type
			}
		case field.ParentGoPackage != "" && !tMap.Has(field.Type):
			typ = "unknown"
			comment += " " + field.ParentGoPackage + "." + field.Type
		}
		if field.IsEnum {
			enumSet.Add(typ)
		}
		if field.IsList {
			typ += "[]"
		}
		columns = append(columns, &Column{
			Name:    core.ToCamelCase(field.SnakeName),
			Type:    typ,
			Comment: comment,
		})
	}
	types = append(types, &Type{
		Name:    name,
		Comment: msg.Comment,
		Columns: columns,
	})
	return types
}

func (c *creator) convertTsType(typ string) string {
	if strings.HasPrefix(typ, "int") || strings.HasPrefix(typ, "float") || strings.HasPrefix(typ, "double") {
		return "number"
	}
	if typ == "bool" {
		return "boolean"
	}
	if typ == "Any" {
		return "unknown"
	}
	return typ
}

type typeMap map[string]string

func (t typeMap) Add(typ, name string) {
	t[typ] = name
}

func (t typeMap) Get(typ string) string {
	return t[typ]
}

func (t typeMap) Has(typ string) bool {
	_, ok := t[typ]
	return ok
}
