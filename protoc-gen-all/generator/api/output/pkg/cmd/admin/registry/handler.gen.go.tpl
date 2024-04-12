{{ template "autogen_comment" }}
package registry

import (
{{- range . }}
	"github.com/xhayamix/proto-gen-golang/pkg/cmd/admin/handler/{{ .PackageName }}"
{{- end }}
	"github.com/xhayamix/proto-gen-golang/pkg/cmd/admin/router"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/entity/config"
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
