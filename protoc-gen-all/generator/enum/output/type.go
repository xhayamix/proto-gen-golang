package output

import (
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/input"
)

type TemplateInfo struct {
	Data     []byte
	FilePath string
}

type EachTemplateCreator interface {
	Create(enum *input.Enum) (*TemplateInfo, error)
}

type BulkTemplateCreator interface {
	Create(enums []*input.Enum) (*TemplateInfo, error)
}
