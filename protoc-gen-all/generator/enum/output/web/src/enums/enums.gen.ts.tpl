{{ template "autogen_comment" }}
{{- $name := .PascalName }}
enum {{ .PascalName }} {
  "" = 0,
{{- range .Elements }}
  {{ .PascalName }} = {{ .Value }},
{{- end }}
}

namespace {{ .PascalName }} {
  export function values(): {{ .PascalName }}[] {
    return [
    {{- range .Elements }}
      {{ $name }}.{{ .PascalName }},
    {{- end }}
    ];
  }

  export function strings(): string[] {
    return [
    {{- range .Elements }}
      "{{ .PascalName }}",
    {{- end }}
    ];
  }

  export function toEnum(s: string): {{ .PascalName }} {
    switch (s) {
    {{- range .Elements }}
      case "{{ .PascalName }}": return {{ $name }}.{{ .PascalName }};
    {{- end }}
      default: return {{ $name }}[""];
    }
  }

  export function toComment(v: {{ .PascalName }}): string {
    switch (v) {
    {{- range .Elements }}
      case {{ $name }}.{{ .PascalName }}: return "{{ .Comment }}";
    {{- end }}
      default: return "";
    }
  }
}

export { {{ .PascalName }} };
