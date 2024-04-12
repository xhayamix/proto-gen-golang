package input

import (
	"strings"

	"github.com/scylladb/go-set/i32set"
	"github.com/scylladb/go-set/strset"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"

	options "github.com/xhayamix/proto-gen-golang/pkg/domain/proto/definition/options/enums"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

func ConvertMessageFromProto(file *protogen.File) (*Enum, error) {
	if len(file.Messages) != 1 {
		return nil, perrors.Newf("このprotoファイルはメッセージの定義数を1にする必要があります。 file = %v", file.Desc.FullName())
	}

	message := file.Messages[0]

	define, ok := proto.GetExtension(message.Desc.Options(), options.E_Define).(*options.DefineOption)
	if !ok {
		return nil, perrors.Newf("型アサーションに失敗しました")
	}

	var accessorType AccessorType
	t := define.GetAccessorType()
	switch t {
	case options.DefineOption_Unknown:
		return nil, perrors.Newf("サポートされていないAccessorTypeです。 AccessorType = %v", t)
	case options.DefineOption_OnlyServer:
		accessorType = AccessorType_OnlyServer
	case options.DefineOption_ServerAndClient:
		accessorType = AccessorType_ServerAndClient
	default:
		return nil, perrors.Newf("サポートされていないAccessorTypeです。 AccessorType = %v", t)
	}

	elementLength := len(define.GetElements())
	ret := &Enum{
		AccessorType: accessorType,
		SnakeName:    core.ToSnakeCase(string(message.Desc.FullName().Name())),
		Comment:      core.CommentReplacer.Replace(message.Comments.Leading.String()),
		Elements:     make([]*Element, 0, elementLength),
	}

	nameSet := strset.NewWithSize(elementLength)
	numSet := i32set.NewWithSize(elementLength)
	isSetting := strings.Contains(ret.SnakeName, "setting_type")
	for _, element := range define.GetElements() {
		if nameSet.Has(element.GetName()) {
			return nil, perrors.Newf(
				"Enum名が重複しています。 Type = %s, Name = %s",
				string(message.Desc.FullName().Name()), element.GetName(),
			)
		}
		if numSet.Has(element.GetValue()) {
			return nil, perrors.Newf(
				"Enum値が重複しています。 Type = %s, Name = %s, Value = %d",
				string(message.Desc.FullName().Name()), element.GetName(), element.GetValue(),
			)
		}

		var settingAccessorType SettingAccessorType
		var settingType SettingType
		var isServerConstant bool
		if isSetting {
			switch element.GetSettingAccessorType() {
			case options.DefineOption_Element_All:
				settingAccessorType = SettingAccessorType_All
			case options.DefineOption_Element_OnlyServer:
				settingAccessorType = SettingAccessorType_OnlyServer
			case options.DefineOption_Element_OnlyClient:
				settingAccessorType = SettingAccessorType_OnlyClient
			default:
				return nil, perrors.Newf("サポートされていないSettingAccessorTypeです。 SettingAccessorType = %v", element.GetSettingAccessorType())
			}

			switch element.GetSettingType() {
			case options.DefineOption_Element_Bool:
				settingType = SettingType_Bool
			case options.DefineOption_Element_Int32:
				settingType = SettingType_Int32
			case options.DefineOption_Element_Int64:
				settingType = SettingType_Int64
			case options.DefineOption_Element_String:
				settingType = SettingType_String
			case options.DefineOption_Element_Int32List:
				settingType = SettingType_Int32List
			case options.DefineOption_Element_Int64List:
				settingType = SettingType_Int64List
			case options.DefineOption_Element_StringList:
				settingType = SettingType_StringList
			default:
				return nil, perrors.Newf("サポートされていないSettingTypeです。 SettingType = %v", element.GetSettingType())
			}

			isServerConstant = element.GetServerConstant()
		}

		ret.Elements = append(ret.Elements, &Element{
			RawName:             element.GetName(),
			Value:               element.GetValue(),
			Comment:             element.GetComment(),
			SettingAccessorType: settingAccessorType,
			SettingType:         settingType,
			IsServerConstant:    isServerConstant,
		})
		nameSet.Add(element.GetName())
		numSet.Add(element.GetValue())
	}

	return ret, nil
}
