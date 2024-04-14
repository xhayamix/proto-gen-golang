package registry

import (
	"bytes"
	_ "embed"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/mysql/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/mysql/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed mysql_repository.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(
	messages []*input.Message,
	fkParentMap map[string]map[string]*output.FK,
	fkChildMap map[string]map[string][]*output.FK,
) (*output.TemplateInfo, error) {
	data := make([]string, 0, len(messages))

	for _, message := range messages {
		data = append(data, core.ToGolangPascalCase(message.SnakeName))
	}

	tpl, err := core.GetBaseTemplate().Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/registry", "mysql_repository.gen.go"),
	}, nil
}
