package input

type AccessorType int32

const (
	AccessorType_OnlyServer AccessorType = iota + 1
	AccessorType_ServerAndClient
)

type SettingType int32

const (
	SettingType_Bool SettingType = iota + 1
	SettingType_Int32
	SettingType_Int64
	SettingType_String
	SettingType_Int32List
	SettingType_Int64List
	SettingType_StringList
)

type SettingAccessorType int32

const (
	SettingAccessorType_All SettingAccessorType = iota + 1
	SettingAccessorType_OnlyServer
	SettingAccessorType_OnlyClient
)

type Element struct {
	// Protoに定義してある名前
	RawName             string
	Value               int32
	Comment             string
	SettingAccessorType SettingAccessorType
	SettingType         SettingType
	IsServerConstant    bool
}

type Enum struct {
	AccessorType AccessorType
	SnakeName    string
	Comment      string
	Elements     []*Element
}
