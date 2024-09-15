package api

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed pb.common_response.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(file *input.File) (*output.TemplateInfo, error) {
	if file.Service == nil {
		return nil, nil
	}

	service := file.Service
	if len(service.Methods) == 0 {
		return nil, nil
	}

	data := make([]string, 0, len(service.Methods))
	for _, method := range service.Methods {
		if method.DisableCommonResponse {
			continue
		}
		outputSnakeName := method.OutputMessage.SnakeName

		if outputSnakeName != "empty" && !strings.HasSuffix(method.OutputMessage.SnakeName, "_response") {
			continue
		}

		var name string
		if outputSnakeName == "empty" {
			name = core.ToPascalCase(service.SnakeName) + core.ToPascalCase(method.SnakeName) + "Response"
		} else {
			name = core.ToPascalCase(outputSnakeName)
		}

		data = append(data, name)
	}
	if len(data) == 0 {
		return nil, nil
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
		FilePath: core.JoinPath("pkg/domain/proto/client/api", file.SnakeName+".pb.common_response.gen.go"),
	}, nil
}
