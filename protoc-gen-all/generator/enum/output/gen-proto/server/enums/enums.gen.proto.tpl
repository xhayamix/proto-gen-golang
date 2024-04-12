{{ template "autogen_comment" }}
syntax = "proto3";

package server.enums;

option go_package = "github.com/QualiArts/campus-server/pkg/domain/proto/server/enums";
{{ range . }}
{{ $Name := .PascalName -}}
enum {{ .PascalName }} {
  {{ .PascalName }}_Unknown = 0;
{{- range .Elements }}
  {{ if .Comment }}// {{ .Comment }}{{ end }}
  {{ $Name }}_{{ .PascalName }} = {{ .Value }};
{{- end }}
}
{{ end -}}
