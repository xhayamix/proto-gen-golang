{{ template "autogen_comment" }}
package {{ .PackageName }}

import (
	"github.com/scylladb/go-set/strset"
)

var {{ .Type }}Check{{ .Name }}MethodSet = strset.New(
{{- range .Methods }}
	"{{ . }}",
{{- end}}
)
