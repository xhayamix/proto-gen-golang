{{ template "autogen_comment" }}

package iamrole

import "github.com/scylladb/go-set/strset"

var onlyAdminUserAPISet = strset.New(
	{{- range $method := .Methods }}
	"{{ $method.HttpMethod }} {{ $method.HttpPath }}",
	{{- end }}
)

func matchOnlyAdminUserAPI(method, path string) bool {
	return onlyAdminUserAPISet.Has(method + " " + path)
}
