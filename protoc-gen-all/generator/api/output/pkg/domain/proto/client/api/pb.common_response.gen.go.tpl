{{ template "autogen_comment" }}
package api

import (
	"github.com/xhayamix/proto-gen-golang/pkg/domain/proto/client/api/common"
)

{{ range . }}
func (m *{{ . }}) SetCommonResponse(res *common.Response) {
	m.CommonResponse = res
}
{{ end }}
