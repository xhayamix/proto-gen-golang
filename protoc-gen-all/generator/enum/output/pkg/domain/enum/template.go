package enum

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed enum.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(enum *input.Enum) (*output.TemplateInfo, error) {
	type Element struct {
		PascalName string
		LowerName  string
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
	isSetting := strings.HasSuffix(enum.SnakeName, "setting_type")
	for _, element := range enum.Elements {
		comment := element.Comment
		if isSetting && element.IsServerConstant {
			comment += "（設定無効）"
		}

		data.Elements = append(data.Elements, &Element{
			PascalName: element.RawName,
			LowerName:  strings.ToLower(element.RawName),
			Value:      element.Value,
			Comment:    comment,
		})
	}

	tpl, err := core.GetBaseTemplate().Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/domain/enum", enum.SnakeName+".gen.go"),
	}, nil
}

//go:embed enum_map.gen.go.tpl
var mapTemplateFileBytes []byte

type MapCreator struct{}

func (c *MapCreator) Create(enums []*input.Enum) (*output.TemplateInfo, error) {
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

	data := make([]*Enum, 0, len(enums))

	for _, enum := range enums {
		e := &Enum{
			PascalName: core.ToPascalCase(enum.SnakeName),
			Comment:    "",
			Elements:   make([]*Element, 0, len(enum.Elements)),
		}
		for _, element := range enum.Elements {
			e.Elements = append(e.Elements, &Element{
				PascalName: element.RawName,
				Value:      element.Value,
				Comment:    element.Comment,
			})
		}

		data = append(data, e)
	}

	tpl, err := core.GetBaseTemplate().Parse(string(mapTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/domain/enum", "enum.gen.go"),
	}, nil
}
