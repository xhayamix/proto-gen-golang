SET FOREIGN_KEY_CHECKS = 0;

{{ range . -}}
{{- $tableName := .SnakeName -}}
CREATE TABLE `{{ $tableName }}` (
  {{- range .Columns }}
  `{{ .SnakeName }}` {{ .Type }}
    {{- if not .Nullable }} NOT NULL{{ end -}}
    {{- if .IsAutoIncrement }} AUTO_INCREMENT{{ end -}}
    {{- if ne .DefaultValue "" }} DEFAULT {{ .DefaultValue }}{{ end -}}
    {{- if .Comment }} COMMENT '{{ .Comment }}'{{ end -}}
    ,
  {{- end }}
  PRIMARY KEY (
    {{- range $i, $pk := .PKColumns -}}
      {{- if $i }}, {{ end }}`{{ $pk.SnakeName }}`
    {{- end -}}
  )
  {{- range .Indexes -}}
  ,
  INDEX `idx{{ range .Keys }}_{{ .SnakeName }}{{ end }}` ({{ range $i, $col := .Keys }}{{ if $i }},{{ end }}`{{ $col.SnakeName }}`{{ end }})
  {{- end }}
  {{- range .FKs -}}
  ,
  CONSTRAINT `fk_{{ $tableName }}_{{ .TargetTableSnakeName }}` FOREIGN KEY (`{{ .FromColumnSnakeName }}`) REFERENCES `{{ .TargetTableSnakeName }}` (`{{ .TargetColumnSnakeName }}`) ON DELETE {{ .OnDelete }} ON UPDATE {{ .OnUpdate }}
  {{- end }}
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4
COMMENT='{{ .Comment }}';

{{ end -}}
SET FOREIGN_KEY_CHECKS = 1;
