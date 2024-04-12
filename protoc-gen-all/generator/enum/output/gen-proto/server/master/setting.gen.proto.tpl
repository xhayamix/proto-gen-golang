{{ template "autogen_comment" }}
syntax = "proto3";

package server.master;
option go_package = "github.com/xhayamix/proto-gen-golang/pkg/domain/entity/master";

import "server/enums/enums_gen.proto";
import "server/options/master/master.proto";

// {{ .Comment }}
message {{ .PascalPrefix }}Setting {
  {{ if len .PascalPrefix -}}
  // 設定ID
  string id = 9999 [(server.options.master.field) = {
    ddl: { pk: true },
    validations: [ { key: "required" } ]
  }];
  {{ end }}
  // 設定種別
  enums.{{ .PascalPrefix }}SettingType settingType = 9998 [(server.options.master.field) = {
    accessorType: AdminAndServer,
    ddl: { pk: true },
    validations: [ { key: "required" } ]
  }];

  // 設定値
  string value = 9997 [(server.options.master.field) = {
    accessorType: AdminAndServer,
    validations: [ { key: "required" } ]
  }];
  {{ range $i, $e := .Elements }}
  // {{ $e.Comment }}
  {{ if .IsList }}repeated {{ end }}{{ $e.SettingType }} {{ $e.CamelName }} = {{ $e.Value }} [(server.options.master.field) = {
    accessorType: OnlyServer,
  }];
  {{ end -}}
}
