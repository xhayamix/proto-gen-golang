package input

type TypeKind int32

const (
	TypeKind_Primitive TypeKind = iota + 1
	TypeKind_Enum
	TypeKind_Map
	TypeKind_Message
)

type PkgType int32

const (
	PkgType_API PkgType = iota + 1
	PkgType_APICommon
	PkgType_ClientCommon
	PkgType_ServerCommon
	PkgType_ClientTransaction
)

type MasterRef struct {
	TableSnakeName        string
	ColumnSnakeName       string
	ParentFieldSnakeNames []string
}

type FieldOption struct {
	// nilチェックが必要
	MasterRef *MasterRef
}

type Field struct {
	PkgType PkgType
	// ImportFileName 外部ファイルのフィールドならそのファイル名(hoge.proto)が入り、
	// このフィールドが定義されているファイルと同じファイルのフィールドなら空文字が入る
	ImportFileName  string
	ParentGoPackage string
	SnakeName       string
	Comment         string
	Type            string
	TypeKind        TypeKind
	IsList          bool
	IsEnum          bool
	MapKeyType      string
	MapValueType    string
	Number          int32
	ValidateOption  string
	HiddenOption    bool
	FieldOption     *FieldOption
}

type Message struct {
	SnakeName string
	Comment   string
	Messages  []*Message
	Fields    []*Field
}

type Method struct {
	SnakeName                 string
	Comment                   string
	InputMessage              *Message
	OutputMessage             *Message
	HttpMethod                string
	HttpPath                  string
	DisableCommonResponse     bool
	DisableResponseCache      bool
	DisableCheckMaintenance   bool
	DisableCheckAppVersion    bool
	DisableCheckLoginToday    bool
	DisableFeatureMaintenance bool
	FeatureMaintenanceTypes   []string
	DisableGameAuthToken      bool
	DisableMasterVersion      bool
	EnableRequestSignature    bool
	CheckOption               string
	ErrorOption               string
}

type Service struct {
	FeatureMaintenanceTypes []string
	SnakeName               string
	Comment                 string
	Methods                 []*Method
}

type File struct {
	IsCommon    bool
	SnakeName   string
	PackageName string
	// Service commonなどはserviceがないのでnilチェックが必要
	Service  *Service
	Messages []*Message
}
