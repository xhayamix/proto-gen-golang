{{ template "autogen_comment" }}
syntax = "proto3";

package server.enums;

option go_package = "github.com/xhayamix/proto-gen-golang/pkg/domain/proto/server/enums";
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
