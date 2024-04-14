{{ template "autogen_comment" }}
{{- $pkColumns := .PKColumns }}
{{- $goName := .GoName }}
package mysql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/constant"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/dto/column"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/entity"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/enum"
)

const (
	{{ .GoName }}TableName   = "{{ .SnakeName }}"
	{{ .GoName }}Comment     = "{{ .Comment }}"
)

// {{ .Comment }}
type {{ .GoName }} struct {
{{- range .Columns }}
	// {{ .Comment }}
	{{ .GoName }} {{ .Type }} `json:"{{ .SnakeName }},omitempty"`
{{- end }}
}

func (e *{{ .GoName }}) GetPK() *{{ .GoName }}PK {
	return &{{ .GoName }}PK{
	{{- range .PKColumns }}
		{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}

func (e *{{ .GoName }}) PK() string {
	return e.GetPK().Key()
}

func (e *{{ .GoName }}) ToKeyValue() map[string]interface{} {
	return map[string]interface{}{
	{{- range .Columns }}
		"{{ .SnakeName }}": e.{{ .GoName }},
	{{- end }}
	}
}

func (e *{{ .GoName }}) GetTypeMap() map[string]string {
	return map[string]string{
{{- range .Columns }}
		"{{ .SnakeName }}": "{{ .Type }}",
{{- end }}
	}
}

func (e *{{ .GoName }}) SetKeyValue(columns []string, values []string) []string {
	errs := make([]string, 0, len(columns))
	for index, column := range columns {
		if len(values) <= index {
			break
		}
		value := values[index]
		switch column {
		{{- range .Columns }}
		case "{{ .SnakeName }}":
			{{- if eq "string" .Type }}
			e.{{ .GoName }} = value
			{{- else if eq "*string" .Type }}
			if value != "<nil>" {
				e.{{ .GoName }} = &value
			}
			{{- else if eq "[]byte" .Type }}
			e.{{ .GoName }} = []byte(value)
			{{- else if eq "bool" .Type }}
			if value != "" {
				v, err := strconv.ParseBool(value)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				}
				e.{{ .GoName }} = v
			}
			{{- else if eq "*bool" .Type }}
			if value != "" && value != "<nil>" {
				v, err := strconv.ParseBool(value)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				} else {
					e.{{ .GoName }} = &v
				}
			}
			{{- else if eq "int32" .Type }}
			if value != "" {
				v, err := strconv.ParseInt(value, 0, 32)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				}
				e.{{ .GoName }} = int32(v)
			}
			{{- else if eq "*int32" .Type }}
			if value != "" && value != "<nil>" {
				vi, err := strconv.ParseInt(value, 0, 32)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				} else {
					v := int32(vi)
					e.{{ .GoName }} = &v
				}
			}
			{{- else if eq "int64" .Type }}
			if value != "" {
				v, err := strconv.ParseInt(value, 0, 64)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				}
				e.{{ .GoName }} = v
			}
			{{- else if eq "*int64" .Type }}
			if value != "" && value != "<nil>" {
				v, err := strconv.ParseInt(value, 0, 64)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				} else {
					e.{{ .GoName }} = &v
				}
			}
			{{- else if eq "float32" .Type }}
			if value != "" {
				v, err := strconv.ParseFloat(value, 32)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				}
				e.{{ .GoName }} = float32(v)
			}
			{{- else if eq "*float32" .Type }}
			if value != "" {
				vf, err := strconv.ParseFloat(value, 32)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				} else {
					v := float32(vf)
					e.{{ .GoName }} = &v
				}
			}
			{{- else if eq "float64" .Type }}
			if value != "" {
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				}
				e.{{ .GoName }} = v
			}
			{{- else if eq "*float64" .Type }}
			if value != "" {
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				} else {
					e.{{ .GoName }} = &v
				}
			}
			{{- else if eq "time.Time" .Type }}
			if value != "" {
				var v time.Time
				var err error
				switch {
				case constant.NormalDatetimeRegExp.MatchString(value):
					v, err = time.ParseInLocation("2006/01/02 15:04:05", value, time.Local)
				case constant.HyphenDatetimeRegExp.MatchString(value):
					v, err = time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
				default:
					v, err = time.Parse(time.RFC3339, value)
				}
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				}
				e.{{ .GoName }} = v
			}
			{{- else if eq "*time.Time" .Type }}
			if value != "" && value != "<nil>" {
				var v time.Time
				var err error
				switch {
				case constant.NormalDatetimeRegExp.MatchString(value):
					v, err = time.ParseInLocation("2006/01/02 15:04:05", value, time.Local)
				case constant.HyphenDatetimeRegExp.MatchString(value):
					v, err = time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
				default:
					v, err = time.Parse(time.RFC3339, value)
				}
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
				} else {
					e.{{ .GoName }} = &v
				}
			}
			{{- else if hasPrefix "enum." .Type }}
			err := e.{{ .GoName }}.UnmarshalJSON([]byte(fmt.Sprintf("%#v", value)))
			{{- if hasSuffix "CommaSeparated" .Type }}
			if _, errs := e.{{ .GoName }}.Split(); err != nil || len(errs) > 0 {
			{{- else }}
			if err != nil || (e.{{ .GoName }} == 0 && value != "") {
			{{- end }}
				errs = append(errs, fmt.Sprintf("{{ .SnakeName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
			}
			{{- end }}
		{{- end }}
		}
	}
	return errs
}

func (e *{{ .GoName }}) PtrFromMapping(cols []string) []interface{} {
	ptrs := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		switch col {
		{{- range $column := .Columns }}
		case "{{ .SnakeName }}":
			ptrs = append(ptrs, &e.{{ .GoName }})
		{{- end }}
		}
	}

	return ptrs
}

type {{ .GoName }}Slice []*{{ .GoName }}

func (s {{ .GoName }}Slice) CreateMapByPK() {{ .GoName }}MapByPK {
	m := make({{ .GoName }}MapByPK, len(s))
	for _, e := range s {
		{{- range $i, $_ := slice .PKColumns 0 (sub (len .PKColumns) 1) }}
		{{- $cols := slice $pkColumns 0 (add1 $i) }}
		if _, ok := m{{ range $cols }}[e.{{ .GoName }}]{{ end }}; !ok {
			m{{ range slice $cols }}[e.{{ .GoName }}]{{ end }} = make({{ range slice $pkColumns (add1 $i) (len $pkColumns)}}map[{{ .Type }}]{{ end }}*{{ $goName }})
		}
		{{- end }}
		m{{ range .PKColumns }}[e.{{ .GoName }}]{{ end }} = e
	}
	return m
}

func (s {{ .GoName }}Slice) EachRecord(iterator func(entity.Record) bool) {
	for _, e := range s {
		if !iterator(e) {
			break
		}
	}
}

type {{ .GoName }}MapByPK {{ range .PKColumns }}map[{{ .Type }}]{{ end }}*{{ .GoName }}

func (m {{ .GoName }}MapByPK) Has(keys ...interface{}) bool {
	{{- range $i, $col := .PKColumns }}
	{{- if eq $i 0 }}
	m0 := m
	{{- else }}
	var m{{ $i }} {{ range slice $pkColumns $i (len $pkColumns) }}map[{{ .Type }}]{{ end }}*{{ $goName }}
	{{- end }}
	{{- end }}
	for i, key := range keys {
		switch i {
		{{- range $i, $col := .PKColumns }}
		case {{ $i }}:
			k, ok := key.({{ .Type }})
			if !ok {
				return false
			}
			{{- if eq $i (sub (len $pkColumns) 1) }}
			_, ok = m{{ $i }}[k]
			{{- else }}
			m{{ add1 $i }}, ok = m{{ $i }}[k]
			{{- end }}
			if !ok {
				return false
			}
		{{- end }}
		default:
			return false
		}
	}
	return true
}

type {{ .GoName }}PK struct {
	{{ range .PKColumns -}}
		{{ .GoName }} {{ .Type }}
	{{ end -}}
}

func (e *{{ .GoName }}PK) Generate() []interface{} {
	return []interface{}{
	{{- range .PKColumns }}
		e.{{ .GoName }},
	{{- end }}
	}
}

func (e *{{ .GoName }}PK) Key() string {
	return strings.Join([]string{
	{{- range .PKColumns }}
		fmt.Sprint(e.{{ .GoName }}),
	{{- end }}
	}, ".")
}

var {{ .GoName }}ColumnName = struct{
	{{ range .Columns -}}
		{{ .GoName }} string
	{{ end -}}
}{
	{{ range .Columns -}}
		{{ .GoName }}: "{{ .SnakeName }}",
	{{ end -}}
}

{{ $name := .GoName -}}
var {{ .GoName }}ColumnMap = map[string]*column.Column{
	{{- range .Columns }}
	{{ $name }}ColumnName.{{ .GoName }}: {
		Name:     "{{ .SnakeName }}",
		Type:     "{{ .Type }}",
		PK:       {{ .PK }},
		Nullable: {{ .Nullable }},
		Required: {{ .Required }},
		Comment:  "{{ .Comment }}",
		CSV:      {{ .CSV }},
		{{ if .HasFK -}}
		FKTarget: "{{ .FKTarget }}",
		FKKey:    "{{ .FKKey }}",
		{{- end }}
	},
	{{- end }}
}

var {{ .GoName }}Cols = []string{
	{{ range .Columns -}}
		{{ $name }}ColumnName.{{ .GoName }},
	{{ end -}}
}

var {{ .GoName }}Columns = column.Columns{
	{{- range .Columns }}
	{{ $name }}ColumnMap[{{ $name }}ColumnName.{{ .GoName }}],
	{{- end }}
}

var {{ .GoName }}PKCols = []string{
	{{ range .PKColumns -}}
		{{ $name }}ColumnName.{{ .GoName }},
	{{ end -}}
}

var {{ .GoName }}PKColumns = column.Columns{
	{{- range .PKColumns }}
	{{ $name }}ColumnMap[{{ $name }}ColumnName.{{ .GoName }}],
	{{- end }}
}

var {{ .GoName }}FKParentTable = strset.New(
	{{ range .FKParentTables -}}
		{{ . }}TableName,
	{{ end -}}
)

var {{ .GoName }}FKChildTable = strset.New(
	{{ range .FKChildTables -}}
		{{ . }}TableName,
	{{ end -}}
)
