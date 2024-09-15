package mysql

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed init.gen.sql.tpl
var templateFileBytes []byte

// 65535byte / (4byte(utf8mb4) * 2)
const mysqlMaxVarcharSize = 8191

type mysqlType = string

const (
	mysqlTypeBool   mysqlType = "TINYINT(1)"
	mysqlTypeInt    mysqlType = "INT"
	mysqlTypeBigInt mysqlType = "BIGINT"
	mysqlTypeString mysqlType = "TEXT"
	mysqlTypeBytes  mysqlType = "MEDIUMBLOB"
	mysqlTypeTime   mysqlType = "DATETIME"
)

type Creator struct{}

func (c *Creator) Create(
	messages []*input.Message,
	fkParentMap map[string]map[string]*output.FK,
	fkChildMap map[string]map[string][]*output.FK,
) (*output.TemplateInfo, error) {
	type Column struct {
		SnakeName       string
		Type            string
		Nullable        bool
		IsAutoIncrement bool
		DefaultValue    string
		Comment         string
	}
	type Index struct {
		Keys []*Column
	}
	type FK struct {
		FromColumnSnakeName   string
		TargetTableSnakeName  string
		TargetColumnSnakeName string
		OnDelete              string
		OnUpdate              string
	}
	type Table struct {
		SnakeName string
		Columns   []*Column
		PKColumns []*Column
		Comment   string
		Indexes   []*Index
		FKs       []*FK
	}

	data := make([]*Table, 0, len(messages))

	for _, message := range messages {
		var fks []*FK
		for _, field := range message.Fields {
			fk := field.Option.DDL.FK
			if fk == nil {
				continue
			}

			fks = append(fks, &FK{
				FromColumnSnakeName:   field.SnakeName,
				TargetTableSnakeName:  fk.TableSnakeName,
				TargetColumnSnakeName: fk.ColumnSnakeName,
				OnDelete:              fk.OnDelete.String(),
				OnUpdate:              fk.OnUpdate.String(),
			})
		}

		table := &Table{
			SnakeName: message.SnakeName,
			Columns:   make([]*Column, 0, len(message.Fields)),
			PKColumns: make([]*Column, 0),
			Comment:   message.Comment,
			Indexes:   make([]*Index, 0, len(message.Option.DDL.Indexes)),
			FKs:       fks,
		}
		indexColumnMap := make(map[string]*Column)

		for _, field := range message.Fields {
			/* 型の整形 */
			var typeName string
			var defaultValue string
			nullable := field.Option.DDL.Nullable
			switch field.TypeKind {
			case input.TypeKind_Bool:
				typeName = mysqlTypeBool
				if field.Option.DDL.HasDefault {
					defaultValue = "0"
				}
			case input.TypeKind_Int32:
				typeName = mysqlTypeInt
				if field.Option.DDL.HasDefault {
					defaultValue = "0"
				}
			case input.TypeKind_Int64:
				typeName = mysqlTypeBigInt
				if field.Option.DDL.HasDefault {
					defaultValue = "0"
				}
			case input.TypeKind_String:
				size := field.Option.DDL.Size
				if size == 0 {
					size = 255
				}
				if size <= mysqlMaxVarcharSize {
					typeName = fmt.Sprintf("VARCHAR(%d)", size)
					if field.Option.DDL.HasDefault {
						defaultValue = "''"
					}
				} else {
					typeName = mysqlTypeString // TODO: マイグレーションで削除されないことを確認してMIDIUMTEXT以上にする
					if field.Option.DDL.HasDefault {
						return nil, perrors.Newf("このTypeKindにはデフォルト値を設定できません。 TypeKind = %v", field.TypeKind)
					}
				}
			case input.TypeKind_Bytes:
				typeName = mysqlTypeBytes
				nullable = true
				if field.Option.DDL.HasDefault {
					return nil, perrors.Newf("このTypeKindにはデフォルト値を設定できません。 TypeKind = %v", field.TypeKind)
				}
			case input.TypeKind_Enum:
				typeName = mysqlTypeInt
				if field.Option.DDL.HasDefault {
					defaultValue = "0"
				}
			default:
				return nil, perrors.Newf("サポートされていないTypeKindです。 TypeKind = %v", field.TypeKind)
			}
			isTime := core.IsAdminTimeField(field.SnakeName)
			if isTime {
				typeName = mysqlTypeTime
			}
			if field.IsList {
				if typeName != mysqlTypeString {
					typeName = "VARCHAR(255)"
				}
			}

			column := &Column{
				SnakeName:       field.SnakeName,
				Type:            typeName,
				Nullable:        nullable,
				IsAutoIncrement: field.Option.DDL.IsAutoIncrement,
				DefaultValue:    defaultValue,
				Comment:         field.Comment,
			}

			table.Columns = append(table.Columns, column)
			if field.Option.DDL.PK {
				table.PKColumns = append(table.PKColumns, column)
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
			table.Indexes = append(table.Indexes, i)
		}

		data = append(data, table)
	}

	tpl, err := template.New("").Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("db/ddl/transaction", "init.sql"),
	}, nil
}
