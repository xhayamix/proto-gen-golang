{{ template "autogen_comment" }}
package registry

import (
	"github.com/QualiArts/campus-server/pkg/domain/proto/client/api"
)

var clientSet = []interface{}{
{{- range . }}
	api.New{{ .PascalName }}Client,
{{- end }}
}
