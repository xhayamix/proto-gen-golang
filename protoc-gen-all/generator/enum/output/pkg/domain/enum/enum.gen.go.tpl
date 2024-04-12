{{ template "autogen_comment" }}
package enum

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/scylladb/go-set/i32set"

	cstrings "github.com/xhayamix/proto-gen-golang/pkg/util/strings"
)

const {{ .PascalName }}Name = "{{ .PascalName }}"

type {{ .PascalName }} int32

{{ $Name := .PascalName -}}
const (
	{{- range .Elements }}
	{{ if .Comment }} // {{ .Comment }}{{ end }}
	{{ $Name }}_{{ .PascalName }} {{ $Name }} = {{ .Value }}
	{{- end }}
)

var {{ .PascalName }}Map = map[string]int32{
	{{- range .Elements }}
	"{{ .PascalName }}": {{ .Value }},
	{{- end }}
}

func (e {{ .PascalName }}) Int() int {
	return int(e)
}

func (e {{ .PascalName }}) Int32() int32 {
	return int32(e)
}

func (e {{ .PascalName }}) Int64() int64 {
	return int64(e)
}

func (e {{ .PascalName }}) String() string {
	switch e {
	{{- range .Elements }}
	case {{ $Name }}_{{ .PascalName }}:
		return "{{ .PascalName }}"
	{{- end }}
	case 0:
		return ""
	}
	return strconv.FormatInt(int64(e), 10)
}

func (e {{ .PascalName }}) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *{{ .PascalName }}) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	{{- range .Elements }}
	case "{{ .LowerName }}":
		*e = {{ $Name }}_{{ .PascalName }}
	{{- end }}
	default:
		i, _ := strconv.ParseInt(s, 10, 64)
		*e = {{ .PascalName }}(int32(i))
	}
	return nil
}

func (e {{ .PascalName }}) EncodeSpanner() (interface{}, error) {
	return e.Int64(), nil
}

func (e *{{ .PascalName }}) DecodeSpanner(val interface{}) error {
	strVal, ok := val.(string)
	if !ok {
		return errors.New(fmt.Sprintf("{{ .PascalName }}.DecodeSpanner failed. %#v", val))
	}
	i, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("{{ .PascalName }}.DecodeSpanner failed. %#v, %#v", val, err))
	}
	*e = {{ .PascalName }}(i)
	return nil
}

func (e {{ .PascalName }}) Validate() bool {
	switch e {
	{{- range .Elements }}
	case {{ $Name }}_{{ .PascalName }}:
		return true
	{{- end }}
	}
	return false
}

type {{ .PascalName }}Slice []{{ .PascalName }}

func (e {{ .PascalName }}Slice) First() {{ .PascalName }} {
	if len(e) == 0 {
		return 0
	}
	return e[0]
}

func (e {{ .PascalName }}Slice) Last() {{ .PascalName }} {
	if len(e) == 0 {
		return 0
	}
	return e[len(e)-1]
}

func (e {{ .PascalName }}Slice) Set() *i32set.Set {
	set := i32set.New()
	for _, i := range e {
		set.Add(i.Int32())
	}
	return set
}

func (e {{ .PascalName }}Slice) Each(f func(Enum) bool) {
	for _, i := range e {
		if !f(i) {
			break
		}
	}
}

func (e {{ .PascalName }}Slice) Size() int {
	return len(e)
}

func (e {{ .PascalName }}Slice) Includes(typ {{ .PascalName }}) bool {
	for _, i := range e {
		if i == typ {
			return true
		}
	}

	return false
}

func (e {{ .PascalName }}Slice) Strings() []string {
	s := make([]string, 0, len(e))
	for _, i := range e {
		s = append(s, i.String())
	}
	return s
}

func (e {{ .PascalName }}Slice) ToSet() {{ .PascalName }}Set {
	s := make({{ .PascalName }}Set, len(e))
	for _, i := range e {
		s.Add(i)
	}
	return s
}

func (e {{ .PascalName }}Slice) EncodeSpanner() (interface{}, error) {
	ret := make([]int64, 0, e.Size())
	for _, i := range e {
		ret = append(ret, i.Int64())
	}
	return ret, nil
}

func (e {{ .PascalName }}Slice) Validate() bool {
	for _, i := range e {
		if !i.Validate() {
			return false
		}
	}
	return true
}

var {{ .PascalName }}Values = {{ .PascalName }}Slice{
	{{- range .Elements }}
	{{ $Name }}_{{ .PascalName }},
	{{- end }}
}

type {{ .PascalName }}Set map[{{ .PascalName }}]struct{}

func (s {{ .PascalName }}Set) Has(e {{ .PascalName }}) bool {
	_, ok := s[e]
	return ok
}

func (s {{ .PascalName }}Set) Size() int {
	return len(s)
}

func (s {{ .PascalName }}Set) Add(e {{ .PascalName }}) {
	s[e] = struct{}{}
}

func (s {{ .PascalName }}Set) ToSlice() {{ .PascalName }}Slice {
	slice := make({{ .PascalName }}Slice, 0, len(s))
	for _, v := range {{ .PascalName }}Values {
		if s.Has(v) {
			slice = append(slice, v)
		}
	}
	return slice
}

type {{ .PascalName }}CommaSeparated string

func (e {{ .PascalName }}CommaSeparated) Split() (Enums, []string) {
	var errs []string
	list := cstrings.SplitComma(string(e))
	res := make({{ .PascalName }}Slice, 0, len(list))

	for _, str := range list {
		i, err := strconv.Atoi(str)
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to convert enum.{{ .PascalName }}. %v\n", str))
			continue
		}
		res = append(res, {{ .PascalName }}(i))
	}

	return res, errs
}

func (e {{ .PascalName }}CommaSeparated) String() string {
	list, _ := e.Split()
	res := make([]string, 0, list.Size())
	list.Each(func(i Enum) bool {
		res = append(res, i.String())
		return true
	})
	return strings.Join(res, ",")
}

func (e {{ .PascalName }}CommaSeparated) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *{{ .PascalName }}CommaSeparated) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	list := cstrings.SplitComma(s)
	res := make([]string, 0, len(s))
	for _, str := range list {
		var i {{ .PascalName }}
		err := i.UnmarshalJSON([]byte(`"`+str+`"`))
		if err != nil {
			return err
		}
		res = append(res, strconv.FormatInt(i.Int64(), 10))
	}
	*e = {{ .PascalName }}CommaSeparated(strings.Join(res, ","))
	return nil
}

var {{ .PascalName }}ValueDetails = ValueDetails{
	{{- range .Elements }}
	{Type: {{ $Name }}_{{ .PascalName }}, Comment: "{{ .Comment }}"},
	{{- end }}
}
