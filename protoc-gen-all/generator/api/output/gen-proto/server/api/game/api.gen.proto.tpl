{{ template "autogen_comment" }}
syntax = "proto3";

package client.api{{ if .IsCommon }}.common{{ end }};
option go_package = "github.com/QualiArts/campus-server/pkg/domain/proto/client/api{{ if .IsCommon }}/common{{ end }}";
option csharp_namespace = "Campus.Common.Proto.Client.Api{{ if .IsCommon }}.Common{{ end }}";

{{ range .ImportPaths -}}
import "{{ . }}";
{{ end }}
{{- if .Service }}
// {{ .Service.Comment }}
{{- $serviceName := .Service.PascalName }}
service {{ $serviceName }} {
{{- range .Service.Methods }}
  // {{ .Comment }}
  rpc {{ .PascalName }}({{ .InputType }}) returns ({{ .OutputType }}) {
    option (google.api.http) = {
      post: "/{{ $serviceName }}/{{ .PascalName }}"
      body: "*"
    };
    {{- range .Options }}
    option ({{ .Key }}) = {{ .Value }};
    {{- end }}
  }
{{- end }}
}
{{ end }}
{{ range .Messages -}}
  {{- template "message" . }}
{{ end -}}
