{{ template "autogen_comment" }}
package maintenance

import (
	"github.com/QualiArts/campus-server/pkg/domain/enum"
)

func getFeatureMaintenanceTypes(method string) enum.FeatureMaintenanceTypeSlice {
	switch method {
	{{- range $v := .DataList }}
	case "{{ $v.Method }}":
		return enum.FeatureMaintenanceTypeSlice{
		{{- range $typ := $v.Types }}
			enum.FeatureMaintenanceType_{{ $typ }},
		{{- end }}
		}
	{{- end }}
	default:
		return nil
	}
}
