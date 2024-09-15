{{ template "autogen_comment" }}
package registry

import (
	"github.com/xhayamix/proto-gen-golang/pkg/infra/mysql/repository"
)

var MysqlRepositorySet = []interface{}{
	{{ range . -}}
		repository.New{{ . }}Repository,
	{{ end -}}
}
