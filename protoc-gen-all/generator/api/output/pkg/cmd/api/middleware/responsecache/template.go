package responsecache

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

var _ output.TemplateCreator = (*Creator)(nil)

//go:embed options.gen.go.tpl
var templateFileBytes []byte

type Data struct {
	EnableResponseCacheMethods []*Method
}

type Method struct {
	Name string
	Type string
}

type Creator struct{}

func (c *Creator) Create(files []*input.File) (*output.TemplateInfo, error) {
	enableResponseCacheMethods := make([]*Method, 0)
	for _, file := range files {
		if file.Service == nil {
			continue
		}
		serviceName := core.ToPascalCase(file.Service.SnakeName)
		for _, method := range file.Service.Methods {
			if !method.DisableResponseCache {
				methodName := core.ToPascalCase(method.SnakeName)
				var typeName string
				if method.OutputMessage.SnakeName == "empty" {
					typeName = serviceName + methodName + "Response"
				} else {
					typeName = core.ToPascalCase(method.OutputMessage.SnakeName)
				}
				enableResponseCacheMethods = append(enableResponseCacheMethods, &Method{
					Name: fmt.Sprintf("/%s.%s/%s", file.PackageName, serviceName, methodName),
					Type: typeName,
				})
			}
		}
	}

	tpl, err := core.GetBaseTemplate().Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, &Data{
		EnableResponseCacheMethods: enableResponseCacheMethods,
	}); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/cmd/api/middleware/responsecache", "options.gen.go"),
	}, nil
}
