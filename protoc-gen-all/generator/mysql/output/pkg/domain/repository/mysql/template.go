package mysql

import (
	"bytes"
	_ "embed"

	"github.com/Masterminds/sprig/v3"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/mysql/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/mysql/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed repository.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(
	message *input.Message,
	fkParentMap map[string]map[string]*output.FK,
	fkChildMap map[string]map[string][]*output.FK,
) (*output.TemplateInfo, error) {
	type Column struct {
		GoName string
		Type   string
	}
	type Index struct {
		Keys []*Column
	}
	type Table struct {
		SnakeName string
		PkgName   string
		GoName    string

		PKColumns []*Column
		Indexes   []*Index
	}

	data := &Table{
		SnakeName: message.SnakeName,
		PkgName:   core.ToPkgName(message.SnakeName),
		GoName:    core.ToGolangPascalCase(message.SnakeName),
		PKColumns: make([]*Column, 0),
		Indexes:   make([]*Index, 0),
	}

	indexColumnMap := make(map[string]*Column)
	for _, field := range message.Fields {
		typeName := field.Type
		if field.TypeKind == input.TypeKind_Enum {
			typeName = "enum." + typeName
		}

		column := &Column{
			GoName: core.ToGolangPascalCase(field.SnakeName),
			Type:   typeName,
		}
		if field.Option.DDL.PK {
			data.PKColumns = append(data.PKColumns, column)
		}
		indexColumnMap[field.SnakeName] = column
	}

	for _, index := range message.Option.DDL.Indexes {
		i := &Index{}
		for _, key := range index.SnakeNameKeys {
			for columnName, column := range indexColumnMap {
				if columnName == key {
					i.Keys = append(i.Keys, column)
				}
			}
		}
		data.Indexes = append(data.Indexes, i)
	}

	tpl, err := core.GetBaseTemplate().Funcs(sprig.TxtFuncMap()).Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/domain/repository/mysql", data.SnakeName+".gen.go"),
	}, nil
}
