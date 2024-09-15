package input

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	options "github.com/xhayamix/proto-gen-golang/pkg/domain/proto/server/options/mysql"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

func ConvertMessageFromProto(file *protogen.File) (*Message, error) {
	if len(file.Messages) != 1 {
		return nil, perrors.Newf("このprotoファイルはメッセージの定義数を1にする必要があります。 file = %v", file.Desc.FullName())
	}

	return convert(file.Messages[0])
}

func convert(message *protogen.Message) (*Message, error) {
	messageOption, ok := proto.GetExtension(message.Desc.Options(), options.E_Message).(*options.MessageOption)
	if !ok {
		return nil, perrors.Newf("型アサーションに失敗しました")
	}

	var indexes []*Index
	if messageOption.GetDdl() != nil {
		indexes = make([]*Index, 0, len(messageOption.GetDdl().GetIndexes()))
		for _, index := range messageOption.GetDdl().GetIndexes() {
			snakeNameKeys := make([]string, 0, len(index.GetKeys()))
			for _, key := range index.GetKeys() {
				snakeNameKeys = append(snakeNameKeys, core.ToSnakeCase(key))
			}

			indexes = append(indexes, &Index{SnakeNameKeys: snakeNameKeys})
		}
	}
	ret := &Message{
		SnakeName: core.ToSnakeCase(string(message.Desc.FullName().Name())),
		Comment:   core.CommentReplacer.Replace(message.Comments.Leading.String()),
		Fields:    nil,
		Option: &MessageOption{
			DDL: &MessageOptionDDL{Indexes: indexes},
		},
	}

	inputFields := make([]*Field, 0, len(message.Fields)+1)
	for _, field := range message.Fields {
		var typeName string
		var typeKind TypeKind
		switch field.Desc.Kind() {
		case protoreflect.BoolKind:
			typeName = FieldType_Bool
			typeKind = TypeKind_Bool
		case protoreflect.Int32Kind:
			typeName = FieldType_Int32
			typeKind = TypeKind_Int32
		case protoreflect.Int64Kind:
			typeName = FieldType_Int64
			typeKind = TypeKind_Int64
		case protoreflect.StringKind:
			typeName = FieldType_String
			typeKind = TypeKind_String
		case protoreflect.BytesKind:
			typeName = FieldType_Bytes
			typeKind = TypeKind_Bytes
		case protoreflect.EnumKind:
			typeName = string(field.Desc.Enum().Name())
			typeKind = TypeKind_Enum
		case protoreflect.DoubleKind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
			protoreflect.GroupKind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
			protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind,
			protoreflect.FloatKind, protoreflect.MessageKind:
			return nil, perrors.Newf("サポートされていないKindです。 Kind = %v", field.Desc.Kind().String())
		default:
			return nil, perrors.Newf("サポートされていないKindです。 Kind = %v", field.Desc.Kind().String())
		}

		inputField := &Field{
			SnakeName: core.ToSnakeCase(field.Desc.TextName()),
			Comment:   core.CommentReplacer.Replace(field.Comments.Leading.String()),
			Type:      typeName,
			TypeKind:  typeKind,
			IsList:    field.Desc.IsList(),
			Option:    nil,
		}

		fieldOption, ok := proto.GetExtension(field.Desc.Options(), options.E_Field).(*options.FieldOption)
		if !ok {
			return nil, perrors.Newf("型アサーションに失敗しました")
		}
		ddlOption := fieldOption.GetDdl()

		option := &FieldOption{
			DDL: &FieldOptionDDL{
				PK:              ddlOption.GetPk(),
				FK:              nil,
				Size:            ddlOption.GetSize(),
				Nullable:        ddlOption.GetNullable(),
				IsAutoIncrement: ddlOption.GetIsAutoIncrement(),
				HasDefault:      ddlOption.GetHasDefault(),
			},
		}
		if ddlOption.GetFk() != nil {
			onDelete, err := ConvertReferenceOptionFromProto(ddlOption.GetFk().GetOnDelete())
			if err != nil {
				return nil, perrors.Stack(err)
			}
			onUpdate, err := ConvertReferenceOptionFromProto(ddlOption.GetFk().GetOnUpdate())
			if err != nil {
				return nil, perrors.Stack(err)
			}

			option.DDL.FK = &FieldOptionDDLFK{
				TableSnakeName:  core.ToSnakeCase(ddlOption.GetFk().GetTable()),
				ColumnSnakeName: core.ToSnakeCase(ddlOption.GetFk().GetColumn()),
				OnDelete:        onDelete,
				OnUpdate:        onUpdate,
			}
		}
		inputField.Option = option

		inputFields = append(inputFields, inputField)
	}
	ret.Fields = inputFields

	return ret, nil
}

func ConvertReferenceOptionFromProto(in options.FieldOption_DDL_ReferenceOption) (ReferenceOption, error) {
	var out ReferenceOption

	switch in {
	case options.FieldOption_DDL_RESTRICT:
		out = ReferenceOption_RESTRICT
	case options.FieldOption_DDL_CASCADE:
		out = ReferenceOption_CASCADE
	case options.FieldOption_DDL_SET_NULL:
		out = ReferenceOption_SET_NULL
	case options.FieldOption_DDL_NO_ACTION:
		out = ReferenceOption_NO_ACTION
	default:
		return 0, perrors.Newf("サポートされていないReferenceOptionです。 ReferenceOption = %v", in)
	}

	return out, nil
}
