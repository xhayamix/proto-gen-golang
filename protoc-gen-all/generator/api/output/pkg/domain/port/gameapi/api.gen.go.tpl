{{ template "autogen_comment" }}

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_api.go
//go:generate goimports -w --local "github.com/QualiArts/campus-server" mock_$GOPACKAGE/mock_api.go
package gameapi

import (
	"context"

	"github.com/QualiArts/campus-server/pkg/domain/proto/client/api"
)

type Method string

var (
{{- range $service := . }}
    // {{ $service.PascalName }}
{{ range $method := $service.Methods -}}
	{{ $service.PascalName }}{{ $method.PascalName }} Method = "{{ $service.PascalName }}{{ $method.PascalName }}"
{{ end -}}
{{ end -}}
)

type GameAPI interface {
{{ range $service := . -}}
{{ range $method := $service.Methods -}}
	// {{ $service.PascalName }}{{ $method.PascalName }} {{ $method.Description }}
	{{ $service.PascalName }}{{ $method.PascalName }}(ctx context.Context {{ if not $method.IsRequestEmpty }}, req *api.{{ $service.PascalName }}{{ $method.PascalName }}Request {{ end }} ) (*api.{{ $service.PascalName }}{{ $method.PascalName }}Response, error)
{{ end -}}
{{ end -}}
	// Close クローズ処理
	Close() error
}
