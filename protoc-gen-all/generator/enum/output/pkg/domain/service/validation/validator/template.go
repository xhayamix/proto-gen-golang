package mcache

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed validator.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(enum *input.Enum) (*output.TemplateInfo, error) {
	if !strings.HasSuffix(enum.SnakeName, "setting_type") {
		return nil, nil
	}

	type Element struct {
		PascalName        string
		PascalSettingType string
		IsServerConstant  bool
	}
	type Enum struct {
		Name         string
		PascalName   string
		CamelName    string
		LowerName    string
		PascalPrefix string
		Elements     []*Element
	}

	name := strings.ReplaceAll(enum.SnakeName, "_type", "")
	prefix := strings.ReplaceAll(strings.ReplaceAll(enum.SnakeName, "setting_type", ""), "_", "")
	data := &Enum{
		Name:         name,
		PascalName:   core.ToPascalCase(name),
		CamelName:    core.ToCamelCase(name),
		LowerName:    strings.ReplaceAll(name, "_", ""),
		PascalPrefix: core.ToPascalCase(prefix),
		Elements:     make([]*Element, 0, len(enum.Elements)),
	}
	for _, element := range enum.Elements {
		var pascalSettingType string
		switch element.SettingType {
		case input.SettingType_Bool:
			pascalSettingType = "Bool"
		case input.SettingType_Int32:
			pascalSettingType = "Int32"
		case input.SettingType_Int64:
			pascalSettingType = "Int64"
		case input.SettingType_String:
			pascalSettingType = "String"
		case input.SettingType_Int32List:
			pascalSettingType = "Int32Slice"
		case input.SettingType_Int64List:
			pascalSettingType = "Int64Slice"
		case input.SettingType_StringList:
			pascalSettingType = "StringSlice"
		default:
			return nil, perrors.Newf("サポートされていないSettingTypeです。 SettingType = %v", element.SettingType)
		}

		data.Elements = append(data.Elements, &Element{
			PascalName:        element.RawName,
			PascalSettingType: pascalSettingType,
			IsServerConstant:  element.IsServerConstant,
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
		FilePath: core.JoinPath("pkg/domain/service/validation/validator", name+".gen.go"),
	}, nil
}
