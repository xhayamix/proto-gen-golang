package enums

import (
	"bytes"
	_ "embed"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed enums.gen.proto.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(enums []*input.Enum) (*output.TemplateInfo, error) {
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
		FilePath: core.JoinPath("proto/server/enums", "enums_gen.proto"),
	}, nil
}
