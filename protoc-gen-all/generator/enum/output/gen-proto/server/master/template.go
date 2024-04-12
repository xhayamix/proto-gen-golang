package enums

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/huandu/xstrings"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed setting.gen.proto.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(enum *input.Enum) (*output.TemplateInfo, error) {
	if !strings.HasSuffix(enum.SnakeName, "setting_type") {
		return nil, nil
	}

	type Element struct {
		CamelName   string
		Value       int32
		Comment     string
		SettingType string
		IsList      bool
	}
	type Enum struct {
		PascalPrefix string
		Comment      string
		Elements     []*Element
	}

	prefix := strings.ReplaceAll(strings.ReplaceAll(enum.SnakeName, "setting_type", ""), "_", "")
	data := &Enum{
		PascalPrefix: core.ToPascalCase(prefix),
		Comment:      enum.Comment,
		Elements:     make([]*Element, 0, len(enum.Elements)),
	}
	for _, element := range enum.Elements {
		if element.SettingAccessorType != input.SettingAccessorType_All && element.SettingAccessorType != input.SettingAccessorType_OnlyServer {
			continue
		}

		var settingType string
		var isList bool
		switch element.SettingType {
		case input.SettingType_Bool:
			settingType = "bool"
		case input.SettingType_Int32:
			settingType = "int32"
		case input.SettingType_Int64:
			settingType = "int64"
		case input.SettingType_Int32List:
			settingType = "int32"
			isList = true
		case input.SettingType_Int64List:
			settingType = "int64"
		case input.SettingType_String:
			settingType = "string"
		case input.SettingType_StringList:
			settingType = "string"
			isList = true
		default:
			return nil, perrors.Newf("サポートされていないSettingTypeです。 SettingType = %v", element.SettingType)
		}

		data.Elements = append(data.Elements, &Element{
			CamelName:   xstrings.FirstRuneToLower(element.RawName),
			Value:       element.Value,
			Comment:     element.Comment,
			SettingType: settingType,
			IsList:      isList,
		})
	}

	tpl, err := core.GetBaseTemplate().Funcs(sprig.TxtFuncMap()).Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	fileName := "setting_gen.proto"
	if prefix != "" {
		fileName = prefix + "_" + fileName
	}
	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("gen-proto/server/master", fileName),
	}, nil
}
