{{ template "autogen_comment" }}
package registry

import (
{{- range . }}
	"github.com/QualiArts/campus-server/pkg/cmd/admin/handler/{{ .PackageName }}"
{{- end }}
	"github.com/QualiArts/campus-server/pkg/cmd/admin/router"
	"github.com/QualiArts/campus-server/pkg/domain/entity/config"
)

var handlerSet = []interface{}{
{{- range . }}
	{{ .PackageName }}.New,
{{- end }}
}

func registerHandler(
	r router.Router,
	c *config.AdminConfig,
	authGroup router.AuthGroup,
	apiGroup router.APIGroup,
	googleGroup router.GoogleGroup,
{{- range . }}
	{{ .CamelName }}Handler {{ .PackageName }}.{{ .PascalName }}Handler,
{{- end }}
) {
{{- range . }}
	{{ .PackageName }}.Register{{ .PascalName }}Handler(r, c, authGroup, apiGroup, googleGroup, {{ .CamelName }}Handler)
{{- end }}
}
