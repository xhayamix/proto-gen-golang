package column

const (
	TypeInt         = "int"
	TypeNullInt     = "*int"
	TypeInt32       = "int32"
	TypeNullInt32   = "*int32"
	TypeFloat32     = "float32"
	TypeNullFloat32 = "*float32"
	TypeBool        = "bool"
	TypeNullBool    = "*bool"
	TypeTime        = "time.Time"
	TypeNullTime    = "*time.Time"
)

type Column struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DatabaseType string `json:"databaseType"`
	PK           bool   `json:"pk"`
	Nullable     bool   `json:"nullable"`
	Required     bool   `json:"required"`
	Comment      string `json:"comment"`
	CSV          bool   `json:"csv"`
	FKTarget     string `json:"fkTarget"`
	FKKey        string `json:"fkKey"`
}

type Columns []*Column

type Index struct {
	Keys    []string `json:"keys"`
	Storing []string `json:"storing"`
	Unique  bool     `json:"unique"`
	ORDER   string   `json:"order"`
}

type Indexes []*Index

type Option struct {
	ID    interface{} `json:"id"`
	Label string      `json:"label"`
}

type Options []*Option

type LogColumn struct {
	Name         string     `json:"name"`
	Comment      string     `json:"comment"`
	Type         string     `json:"type"`
	DatabaseType string     `json:"databaseType"`
	ChildColumns LogColumns `json:"childColumns"`
}

type LogColumns []*LogColumn
