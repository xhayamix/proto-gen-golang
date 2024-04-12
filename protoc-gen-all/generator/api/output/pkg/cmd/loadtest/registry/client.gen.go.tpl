{{ template "autogen_comment" }}
package registry

import (
	"github.com/xhayamix/proto-gen-golang/pkg/domain/proto/client/api"
)

var clientSet = []interface{}{
{{- range . }}
	api.New{{ .PascalName }}Client,
{{- end }}
}
