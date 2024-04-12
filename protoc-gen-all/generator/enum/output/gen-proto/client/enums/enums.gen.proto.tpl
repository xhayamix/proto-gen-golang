{{ template "autogen_comment" }}
syntax = "proto3";

package client.enums;

option go_package = "github.com/QualiArts/campus-server/pkg/domain/proto/client/enums";
option csharp_namespace = "Campus.Common.Proto.Client.Enums";
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
