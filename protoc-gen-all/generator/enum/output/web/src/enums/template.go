package enums

import (
	"bytes"
	_ "embed"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed index.gen.ts.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(enums []*input.Enum) (*output.TemplateInfo, error) {
	type Enum struct {
		PascalName string
		Comment    string
	}

	data := make([]*Enum, 0, len(enums))

	for _, enum := range enums {
		data = append(data, &Enum{
			PascalName: core.ToPascalCase(enum.SnakeName),
			Comment:    enum.Comment,
		})
	}

	tpl := template.Must(core.GetBaseTemplate().Funcs(sprig.TxtFuncMap()).Parse(string(templateFileBytes)))
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("web/src/enums", "index.gen.ts"),
	}, nil
}

//go:embed enums.gen.ts.tpl
var enumTemplateFileBytes []byte

type EachCreator struct{}

func (c *EachCreator) Create(enum *input.Enum) (*output.TemplateInfo, error) {
	type Element struct {
		PascalName string
		Value      int32
		Comment    string
	}

	type Enum struct {
		PascalName string
		Comment    string
		Elements   []*Element
	}

	data := &Enum{
		PascalName: core.ToPascalCase(enum.SnakeName),
		Comment:    "",
		Elements:   make([]*Element, 0, len(enum.Elements)),
	}

	for _, element := range enum.Elements {
		pascalName := element.RawName
		if _, err := strconv.Atoi(pascalName); err == nil {
			pascalName = "_" + pascalName
		}
		data.Elements = append(data.Elements, &Element{
			PascalName: pascalName,
			Value:      element.Value,
			Comment:    element.Comment,
		})
	}

	tpl, err := core.GetBaseTemplate().Parse(string(enumTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("web/src/enums", data.PascalName+".gen.ts"),
	}, nil
}
