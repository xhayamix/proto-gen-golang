package api

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

var _ output.TemplateCreator = (*Creator)(nil)

//go:embed api.gen.go.tpl
var templateFileBytes []byte

type Service struct {
	PascalName string
	Methods    []*Method
}

type Method struct {
	PascalName     string
	Description    string
	IsRequestEmpty bool
}

type Creator struct{}

func (c *Creator) Create(files []*input.File) (*output.TemplateInfo, error) {
	infos := make([]*Service, 0, len(files))
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
				Description:    method.Comment,
				IsRequestEmpty: method.InputMessage.SnakeName == "empty",
			})
		}
		infos = append(infos, &Service{
			PascalName: core.ToPascalCase(service.SnakeName),
			Methods:    methods,
		})
	}
	if len(infos) == 0 {
		return nil, nil
	}

	buf := &bytes.Buffer{}
	tpl := template.Must(core.GetBaseTemplate().Parse(string(templateFileBytes)))
	if err := tpl.Execute(buf, infos); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/domain/port/api", "api.gen.go"),
	}, nil
}
