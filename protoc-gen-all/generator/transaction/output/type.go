package output

import (
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/input"
)

type FK struct {
	TableSnakeName  string
	ColumnSnakeName string
	OnDelete        string
	OnUpdate        string
}

type TemplateInfo struct {
	Data     []byte
	FilePath string
}

type EachTemplateCreator interface {
	Create(
		message *input.Message,
		fkParentMap map[string]map[string]*FK,
		fkChildMap map[string]map[string][]*FK,
	) (*TemplateInfo, error)
}

type BulkTemplateCreator interface {
	Create(
		messages []*input.Message,
		fkParentMap map[string]map[string]*FK,
		fkChildMap map[string]map[string][]*FK,
	) (*TemplateInfo, error)
}
