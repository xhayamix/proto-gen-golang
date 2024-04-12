{{ template "autogen_comment" }}
package {{ .PackageName }}

import (
	"net/http"

	"github.com/labstack/echo/v4"
{{- range .ImportPaths }}
	{{ . }}
{{- end }}
	"github.com/QualiArts/campus-server/pkg/cerrors"
	"github.com/QualiArts/campus-server/pkg/cmd/admin/router"
	"github.com/QualiArts/campus-server/pkg/domain/entity/config"
)

type {{ .PascalName }}Handler interface {
{{- range .Methods }}
	// {{ .Comment }}
	{{ .PascalName }}(c echo.Context, req *{{ .InputType }}) ({{ if not .OutputAny }}*{{ end }}{{ .OutputType }}{{ .OutputTypeSuffix }}, error)
{{- end }}
}

{{- range .Methods }}
{{- if (gt (len .InputFields) 0) }}

type {{ .InputType }}{{ .InputTypeSuffix }} struct{
{{- range .InputFields }}
	// {{ .Comment }}
	{{ .PascalName }} {{ .Type }} `param:"{{ .CamelName }}" query:"{{ .CamelName }}" json:"{{ .CamelName }}"`
{{- end }}
}

func (q *{{ .InputType }}{{ .InputTypeSuffix }}) Convert() *{{ .InputType }} {
	return &{{ .InputType }}{
{{- range .InputFields }}
		{{ .PascalName }}: q.{{ .PascalName }},
{{- end }}
	}
}
{{- end }}

{{- if (gt (len .OutputFields) 0) }}

type {{ .OutputType }}{{ .OutputTypeSuffix }} struct{
	*{{ .OutputType }}
{{- range .OutputFields }}
	// {{ .Comment }}
	{{ .PascalName }} {{ .Type }} `json:"{{ .CamelName }}"`
{{- end }}
}
{{- end }}
{{- end }}

func Register{{ .PascalName }}Handler(
	r router.Router,
	conf *config.AdminConfig,
	authGroup router.AuthGroup,
	apiGroup router.APIGroup,
	googleGroup router.GoogleGroup,
	handler {{ .PascalName }}Handler,
) {
{{- range .Methods }}
	{{ .RouterType }}.{{ .Method }}("{{ .Path }}", func(c echo.Context) error {
		{{- if .DisableOnProduction }}
		if conf.Env.IsProduction() {
			return cerrors.Newf(cerrors.Internal, "本番では使用禁止です")
		}
		{{- end }}
		req := &{{ .InputType }}{{ .InputTypeSuffix }}{}
		{{- if not .DisableRequestBind }}
		if err := c.Bind(req); err != nil {
			return cerrors.Wrap(err, cerrors.InvalidArgument)
		}
		{{- end }}
		{{- if eq .InputTypeSuffix "" }}
		req2 := req
		{{- else }}
		req2 := req.Convert()
		{{- end }}
		{{- if not .DisableRequestBind }}
		if err := c.Validate(req2); err != nil {
			return cerrors.Wrap(err, cerrors.InvalidArgument)
		}
		{{- end }}
		res, err := handler.{{ .PascalName }}(c, req2)
		if err != nil {
			return cerrors.Stack(err)
		}
		if c.Response().Committed {
			return nil
		}
		if res == nil {
			return c.NoContent(http.StatusNoContent)
		}
		return c.JSON(http.StatusOK, res)
	})
{{- end }}
}
