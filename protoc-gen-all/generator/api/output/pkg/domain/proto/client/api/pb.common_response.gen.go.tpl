{{ template "autogen_comment" }}
package api

import (
	"github.com/QualiArts/campus-server/pkg/domain/proto/client/api/common"
)

{{ range . }}
func (m *{{ . }}) SetCommonResponse(res *common.Response) {
	m.CommonResponse = res
}
{{ end }}
