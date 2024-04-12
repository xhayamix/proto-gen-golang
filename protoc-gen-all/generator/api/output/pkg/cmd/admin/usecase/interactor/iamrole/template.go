package handler

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed api.gen.go.tpl
var templateFileBytes []byte

type Data struct {
	Methods []*Method
}

type Method struct {
	HttpMethod string
	HttpPath   string
}

func New() output.TemplateCreator {
	tpl := template.Must(core.GetBaseTemplate().Parse(string(templateFileBytes)))
	return &creator{tpl: tpl}
}

type creator struct {
	tpl *template.Template
}

func (c *creator) Create(files []*input.File) (*output.TemplateInfo, error) {
	methods := make([]*Method, 0)
	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, &Data{Methods: methods}); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/cmd/admin/usecase/interactor/iamrole", "api.gen.go"),
	}, nil
}
