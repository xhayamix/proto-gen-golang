{{ template "autogen_comment" }}
{{- range . }}
export * from "./{{ .PascalName }}.gen";
{{- end }}
