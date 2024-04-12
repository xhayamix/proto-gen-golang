package campusserverapi

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed request.gen.go.tpl
var templateFileBytes []byte

type Service struct {
	PascalName string
	Methods    []*Method
}

type Method struct {
	PascalName     string
	IsRequestEmpty bool
}

type Creator struct{}

func (c *Creator) Create(files []*input.File) ([]*output.TemplateInfo, error) {
	tpl := template.Must(core.GetBaseTemplate().Parse(string(templateFileBytes)))
	results := make([]*output.TemplateInfo, 0, len(files))
	for _, file := range files {
		if file.Service == nil {
			continue
		}
		service := file.Service
		if len(service.Methods) == 0 {
			continue
		}

		methods := make([]*Method, 0, len(service.Methods))
		for _, method := range service.Methods {
			methods = append(methods, &Method{
				PascalName:     core.ToPascalCase(method.SnakeName),
				IsRequestEmpty: method.InputMessage.SnakeName == "empty",
			})
		}
		info := &Service{
			PascalName: core.ToPascalCase(service.SnakeName),
			Methods:    methods,
		}

		buf := &bytes.Buffer{}
		if err := tpl.Execute(buf, info); err != nil {
			return nil, perrors.Stack(err)
		}
		results = append(results, &output.TemplateInfo{
			Data:     buf.Bytes(),
			FilePath: core.JoinPath("pkg/infra/campusserverapi", service.SnakeName+".gen.go"),
		})
	}
	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}
