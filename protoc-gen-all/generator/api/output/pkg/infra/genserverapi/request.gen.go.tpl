{{ template "autogen_comment" }}
package genserverapi

import (
	"context"

    "google.golang.org/protobuf/types/known/emptypb"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/proto/client/api"
)

{{- $service := . }}
{{ range $method := $service.Methods }}
func (c *client) {{ $service.PascalName }}{{ $method.PascalName }}(ctx context.Context {{ if not $method.IsRequestEmpty }}, req *api.{{ $method.PascalName }}Request {{ end }} ) (*api.{{ $method.PascalName }}Response, error) {
    conn, err := c.getConn()
	if err != nil {
		return nil, cerrors.Stack(err)
	}
    cli := api.New{{ $service.PascalName }}Client(conn)
    {{- if $method.IsRequestEmpty }}
    result, err := cli.{{ $method.PascalName }}(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
    {{ else }}
    result, err := cli.{{ $method.PascalName }}(ctx, req)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
    {{ end -}}
	return result, nil
}
{{ end }}
