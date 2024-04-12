{{ template "autogen_comment" }}
import request from "@/request";
{{- if len .Commons }}
import {
{{- range .Commons }}
  {{ . }},
{{- end }}
} from "@/common/index.gen";
{{- end }}
{{- if len .Transactions }}
import {
{{- range .Transactions }}
  {{ . }},
{{- end }}
} from "@/entities/transaction/index.gen";
{{- end }}
{{- if len .Enums }}
import {
{{- range .Enums }}
  {{ . }},
{{- end }}
} from "@/enums/index.gen";
{{- end }}

{{- range .Types }}

// {{ .Comment }}
export type {{ .Name }} = {
  {{- range .Columns }}
  // {{ .Comment }}
  {{ .Name }}?: {{ .Type }},
  {{- end }}
};
{{- end }}

{{- range .Methods }}

// {{ .Comment }}
export function {{ .Name }}({{ if ne .RequestType "" }}req: {{ .RequestType }}{{ end }}{{ if eq .HttpMethod "get" }}{{ if ne .RequestType "" }}, {{ end }}async?: boolean, skipRedirect?: boolean{{ end }}): Promise<{{ .ResponseType }}> {
  {{- if eq .HttpMethod "get" }}
  return request.get("{{ .HttpPath }}", {{ if ne .RequestType "" }}req, {{ else }}undefined, {{ end }}async, skipRedirect);
  {{- else }}
  return request.{{ .HttpMethod }}("{{ .HttpPath }}"{{ if ne .RequestType "" }}, undefined, req{{ end }});
  {{- end }}
}
{{- end }}
