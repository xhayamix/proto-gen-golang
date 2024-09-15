package input

type TypeKind int32

const (
	TypeKind_Bool TypeKind = iota + 1
	TypeKind_Int32
	TypeKind_Int64
	TypeKind_String
	TypeKind_Bytes
	TypeKind_Enum
)

type FieldType = string

const (
	FieldType_Bool   = "bool"
	FieldType_Int32  = "int32"
	FieldType_Int64  = "int64"
	FieldType_String = "string"
	FieldType_Bytes  = "[]byte"
)

type ReferenceOption int32

const (
	ReferenceOption_RESTRICT ReferenceOption = iota + 1
	ReferenceOption_CASCADE
	ReferenceOption_SET_NULL
	ReferenceOption_NO_ACTION
)

func (r ReferenceOption) String() string {
	var str string

	switch r {
	case ReferenceOption_RESTRICT:
		str = "RESTRICT"
	case ReferenceOption_CASCADE:
		str = "CASCADE"
	case ReferenceOption_SET_NULL:
		str = "SET_NULL"
	case ReferenceOption_NO_ACTION:
		str = "NO_ACTION"
	}

	return str
}

type FieldOptionDDLFK struct {
	TableSnakeName  string
	ColumnSnakeName string
	OnDelete        ReferenceOption
	OnUpdate        ReferenceOption
}

type FieldOptionDDL struct {
	PK bool
	// nilチェックが必要
	FK              *FieldOptionDDLFK
	Size            uint32
	Nullable        bool
	IsAutoIncrement bool
	HasDefault      bool
}

type FieldOption struct {
	DDL *FieldOptionDDL
}

type Field struct {
	SnakeName string
	Comment   string
	// TypeKind_Enumの場合はEnum名が入る
	Type     FieldType
	TypeKind TypeKind
	IsList   bool
	Option   *FieldOption
}

type Index struct {
	SnakeNameKeys []string
}

type MessageOptionDDL struct {
	Indexes []*Index
}

type MessageOption struct {
	DDL *MessageOptionDDL
}

type Message struct {
	SnakeName string
	Comment   string
	Fields    []*Field
	Option    *MessageOption
}
