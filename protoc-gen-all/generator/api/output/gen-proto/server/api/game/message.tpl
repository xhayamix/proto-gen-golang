{{- define "message" }}
// {{ .Comment }}
message {{ .PascalName }} {
{{- range .Messages }}
  {{- template "message" . }}
{{- end }}
{{- range .Fields }}
  // {{ .Comment }}
  {{ .PkgName }}{{ .Type }} {{ .CamelName }} = {{ .Number }}{{ .Option }};
{{- end }}
{{- if .HasCommonResponse }}
  // 共通レスポンス
  api.common.Response commonResponse = 9999;
{{- end }}
}
{{ end -}}
