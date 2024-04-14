package mysql

import (
	"bytes"
	_ "embed"
	"sort"

	"github.com/Masterminds/sprig/v3"
	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/mysql/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/mysql/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed entity.gen.go.tpl
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
		Comment   string
		Type      string

		PK       bool
		Nullable bool
		Required bool
		CSV      bool
		HasFK    bool
		FKTarget string
		FKKey    string
	}
	type Table struct {
		PkgName   string
		GoName    string
		SnakeName string
		Comment   string
		Columns   []*Column

		PKColumns      []*Column
		FKParentTables []string
		FKChildTables  []string
	}

	data := &Table{
		PkgName:   core.ToPkgName(message.SnakeName),
		GoName:    core.ToGolangPascalCase(message.SnakeName),
		SnakeName: message.SnakeName,
		Comment:   message.Comment,
		Columns:   make([]*Column, 0, len(message.Fields)),

		FKParentTables: make([]string, 0),
		FKChildTables:  make([]string, 0),
	}

	for _, field := range message.Fields {
		/* 型の整形 */
		typeName := field.Type
		isTime := core.IsAdminTimeField(field.SnakeName)
		if isTime {
			typeName = "time.Time"
		}
		nullable := field.Option.DDL.Nullable
		if nullable {
			typeName = "*" + typeName
		}
		isEnum := field.TypeKind == input.TypeKind_Enum
		isList := field.IsList
		if isEnum {
			typeName = "enum." + typeName
			if isList {
				typeName += "CommaSeparated"
			}
		} else if isList {
			typeName = "string"
		}

		var (
			required bool
			hasFK    bool
			fkTarget string
			fkKey    string
		)
		if field.Option.DDL.FK != nil {
			hasFK = true
			fk := field.Option.DDL.FK
			fkTarget = fk.TableSnakeName
			fkKey = fk.ColumnSnakeName
		}

		column := &Column{
			GoName:    core.ToGolangPascalCase(field.SnakeName),
			SnakeName: field.SnakeName,
			Comment:   field.Comment,
			Type:      typeName,

			PK:       field.Option.DDL.PK,
			Nullable: nullable,
			Required: required,
			CSV:      isList,
			HasFK:    hasFK,
			FKTarget: fkTarget,
			FKKey:    fkKey,
		}

		data.Columns = append(data.Columns, column)
		if column.PK {
			data.PKColumns = append(data.PKColumns, column)
		}
	}

	fKParentTableSet := strset.New()
	fKChildTableSet := strset.New()
	columns, ok := fkParentMap[data.SnakeName]
	if ok {
		for _, fk := range columns {
			fKParentTableSet.Add(core.ToGolangPascalCase(fk.TableSnakeName))
		}
	}
	columnMap, ok := fkChildMap[data.SnakeName]
	if ok {
		for _, fks := range columnMap {
			for _, fk := range fks {
				fKChildTableSet.Add(core.ToGolangPascalCase(fk.TableSnakeName))
			}
		}
	}
	data.FKParentTables = fKParentTableSet.List()
	sort.Strings(data.FKParentTables)
	data.FKChildTables = fKChildTableSet.List()
	sort.Strings(data.FKChildTables)

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
		FilePath: core.JoinPath("pkg/domain/entity/mysql", data.SnakeName+".gen.go"),
	}, nil
}
