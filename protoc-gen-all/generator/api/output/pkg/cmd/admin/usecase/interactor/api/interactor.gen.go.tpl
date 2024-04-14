{{ template "autogen_comment" }}
package api

import (
	"context"
	"encoding/json"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/port/api"
	clientapi "github.com/xhayamix/proto-gen-golang/pkg/domain/proto/client/api"
)

func (i *interactor) list(ctx context.Context) ([]api.Method) {
	return []api.Method{
		{{ range $service := . -}}
		{{ range $method := $service.Methods -}}
		api.{{ $service.PascalName }}{{ $method.PascalName }},
		{{ end -}}
		{{ end -}}
	}
}

func (i *interactor) request(ctx context.Context, method api.Method, param string) (interface{}, error) {
	var result interface{}
	var err error
	switch method {
		{{ range $service := . -}}
		{{ range $method := $service.Methods -}}
		case api.{{ $service.PascalName }}{{ $method.PascalName }}:
			{{- if not $method.IsRequestEmpty }}
			req := &clientapi.{{ $method.PascalName }}Request{}
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
