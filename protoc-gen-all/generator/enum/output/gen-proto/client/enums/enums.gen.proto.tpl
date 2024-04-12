{{ template "autogen_comment" }}
syntax = "proto3";

package client.enums;

option go_package = "github.com/xhayamix/proto-gen-golang/pkg/domain/proto/client/enums";
option csharp_namespace = "Gen.Common.Proto.Client.Enums";
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
