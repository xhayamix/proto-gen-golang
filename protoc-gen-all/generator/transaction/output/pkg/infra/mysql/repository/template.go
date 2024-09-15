package repository

import (
	"bytes"
	_ "embed"

	"github.com/Masterminds/sprig/v3"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/output"
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
		GoName    string
		SnakeName string
		Type      string
	}
	type Index struct {
		Keys []*Column
	}
	type Table struct {
		PkgName   string
		GoName    string
		SnakeName string
		CamelName string

		Columns          []*Column
		PKColumns        []*Column
		Indexes          []*Index
		HasMasterVersion bool
	}

	data := &Table{
		PkgName:          core.ToPkgName(message.SnakeName),
		GoName:           core.ToGolangPascalCase(message.SnakeName),
		SnakeName:        message.SnakeName,
		CamelName:        core.ToCamelCase(message.SnakeName),
		Columns:          make([]*Column, 0, len(message.Fields)),
		PKColumns:        make([]*Column, 0),
		Indexes:          make([]*Index, 0),
		HasMasterVersion: !core.IsMasterTagKind(message.SnakeName),
	}

	indexColumnMap := make(map[string]*Column)
	for _, field := range message.Fields {
		typeName := field.Type
		if field.TypeKind == input.TypeKind_Enum {
			typeName = "enum." + typeName
		}

		column := &Column{
			GoName:    core.ToGolangPascalCase(field.SnakeName),
			SnakeName: field.SnakeName,
			Type:      typeName,
		}
		if field.Option.DDL.PK {
			data.PKColumns = append(data.PKColumns, column)
		}
		data.Columns = append(data.Columns, column)
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
		FilePath: core.JoinPath("pkg/infra/transaction/repository", data.SnakeName+".gen.go"),
	}, nil
}
