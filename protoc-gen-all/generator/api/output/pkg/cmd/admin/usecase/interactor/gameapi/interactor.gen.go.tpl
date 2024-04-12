{{ template "autogen_comment" }}
package gameapi

import (
	"context"
	"encoding/json"

	"github.com/QualiArts/campus-server/pkg/cerrors"
	"github.com/QualiArts/campus-server/pkg/domain/port/gameapi"
	"github.com/QualiArts/campus-server/pkg/domain/proto/client/api"
)

func (i *interactor) list(ctx context.Context) ([]gameapi.Method) {
	return []gameapi.Method{
		{{ range $service := . -}}
		{{ range $method := $service.Methods -}}
		gameapi.{{ $service.PascalName }}{{ $method.PascalName }},
		{{ end -}}
		{{ end -}}
	}
}

func (i *interactor) request(ctx context.Context, method gameapi.Method, param string) (interface{}, error) {
	var result interface{}
	var err error
	switch method {
		{{ range $service := . -}}
		{{ range $method := $service.Methods -}}
		case gameapi.{{ $service.PascalName }}{{ $method.PascalName }}:
			{{- if not $method.IsRequestEmpty }}
			req := &api.{{ $service.PascalName }}{{ $method.PascalName }}Request{}
			if err := json.Unmarshal([]byte(param), req); err != nil {
				return "", cerrors.Wrap(err, cerrors.Internal)
			}
			{{ end -}}
			result, err = i.api.{{ $service.PascalName }}{{ $method.PascalName }}(ctx {{ if not $method.IsRequestEmpty }}, req {{ end }})
			if err != nil {
				return "", cerrors.Wrap(err, cerrors.Internal)
			}
		{{ end -}}
		{{ end -}}
		default:
			return "", cerrors.Newf(cerrors.InvalidArgument, "APIが存在しません。 method = %q", method)
	}

	return result, nil
}
