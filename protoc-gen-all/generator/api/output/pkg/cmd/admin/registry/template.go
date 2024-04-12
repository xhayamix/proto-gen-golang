package registry

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed handler.gen.go.tpl
var templateFileBytes []byte

type Data struct {
	PackageName string
	PascalName  string
	CamelName   string
}

func New() output.TemplateCreator {
	tpl := template.Must(core.GetBaseTemplate().Parse(string(templateFileBytes)))
	return &creator{tpl: tpl}
}

type creator struct {
	tpl *template.Template
}

func (c *creator) Create(files []*input.File) (*output.TemplateInfo, error) {
	infos := make([]*Data, 0, len(files))
	for _, file := range files {
		infos = append(infos, &Data{
			PackageName: core.ToPkgName(file.SnakeName),
			PascalName:  core.ToPascalCase(file.SnakeName),
			CamelName:   core.ToCamelCase(file.SnakeName),
		})
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, infos); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/cmd/admin/registry", "handler.gen.go"),
	}, nil
}
