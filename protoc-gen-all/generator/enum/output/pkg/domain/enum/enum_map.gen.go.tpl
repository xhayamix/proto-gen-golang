{{ template "autogen_comment" }}
package enum

import "github.com/scylladb/go-set/i32set"

var EnumTypeSlice = []string{
{{- range . }}
	{{ .PascalName }}Name,
{{- end }}
}

var EnumTypeMap = map[string]map[string]int32{
{{- range . }}
	{{ .PascalName }}Name: {{ .PascalName }}Map,
{{- end }}
}

var EnumValueDetailsMap = map[string]ValueDetails{
{{- range . }}
	{{ .PascalName }}Name: {{ .PascalName }}ValueDetails,
{{- end }}
}

type Enum interface {
	Int() int
	Int32() int32
	Int64() int64
	String() string
	MarshalJSON() ([]byte, error)
	Validate() bool
}

type Enums interface {
	Set() *i32set.Set
	Each(f func(Enum) bool)
	Size() int
	Validate() bool
}

type EnumCommaSeparated interface {
	Split() (Enums, []string)
	String() string
	MarshalJSON() ([]byte, error)
}

type ValueDetail struct {
	Type    Enum
	Comment string
}

type ValueDetails []*ValueDetail
