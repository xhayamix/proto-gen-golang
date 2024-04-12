package output

import (
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
)

type TemplateInfo struct {
	Data     []byte
	FilePath string
}

type EachTemplateCreator interface {
	Create(message *input.File) (*TemplateInfo, error)
}

type TemplateCreator interface {
	Create(files []*input.File) (*TemplateInfo, error)
}

type TemplatesCreator interface {
	Create(files []*input.File) ([]*TemplateInfo, error)
}
