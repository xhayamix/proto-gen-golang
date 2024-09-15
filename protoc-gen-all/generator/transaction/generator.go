package transaction

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/output"
	ddl "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/output/db/ddl/mysql"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/output/pkg/cmd/admin/registry"
	entity "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/output/pkg/domain/entity/transaction"
	domainrepository "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/output/pkg/domain/repository/transaction"
	infrarepository "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/transaction/output/pkg/infra/mysql/repository"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

type generator struct {
	*core.GeneratorBase
	plugin *protogen.Plugin
}

func NewGenerator(plugin *protogen.Plugin) core.Generator {
	return &generator{
		GeneratorBase: core.NewGeneratorBase(),
		plugin:        plugin,
	}
}

var eachCreators = []output.EachTemplateCreator{
	&domainrepository.Creator{},
	&infrarepository.Creator{},
	&entity.Creator{},
}

var bulkCreators = []output.BulkTemplateCreator{
	&ddl.Creator{},
	&registry.Creator{},
}

func (g *generator) Build() ([]core.GenFile, error) {
	messages := make([]*input.Message, 0)

	for _, file := range g.plugin.Files {
		if !file.Generate {
			continue
		}
		if file.Proto.GetPackage() != "server.transaction" {
			continue
		}

		message, err := input.ConvertMessageFromProto(file)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		messages = append(messages, message)
	}

	// 入力ファイルの順番に左右されないようソートする
	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].SnakeName < messages[j].SnakeName
	})

	fkParentMap := make(map[string]map[string]*output.FK)
	fkChildMap := make(map[string]map[string][]*output.FK)
	for _, message := range messages {
		for _, field := range message.Fields {
			if field.Option.DDL.FK == nil {
				continue
			}

			fk := field.Option.DDL.FK
			tableName := core.ToSnakeCase(fk.TableSnakeName)
			columnName := core.ToSnakeCase(fk.ColumnSnakeName)

			if _, ok := fkParentMap[message.SnakeName]; !ok {
				fkParentMap[message.SnakeName] = make(map[string]*output.FK)
			}
			fkParentMap[message.SnakeName][field.SnakeName] = &output.FK{
				TableSnakeName:  tableName,
				ColumnSnakeName: columnName,
				OnDelete:        fk.OnDelete.String(),
				OnUpdate:        fk.OnUpdate.String(),
			}

			if _, ok := fkChildMap[tableName]; !ok {
				fkChildMap[tableName] = make(map[string][]*output.FK)
			}

			fkChildMap[tableName][columnName] = append(fkChildMap[tableName][columnName], &output.FK{
				TableSnakeName:  message.SnakeName,
				ColumnSnakeName: field.SnakeName,
			})
		}
	}

	genFiles := make([]core.GenFile, 0)
	for _, creator := range eachCreators {
		for _, message := range messages {
			info, err := creator.Create(message, fkParentMap, fkChildMap)
			if err != nil {
				return nil, perrors.Stack(err)
			}
			if info == nil {
				continue
			}

			genFiles = append(genFiles, core.NewGenFile(info.FilePath, info.Data))
		}
	}
	for _, creator := range bulkCreators {
		info, err := creator.Create(messages, fkParentMap, fkChildMap)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		if info == nil {
			continue
		}

		genFiles = append(genFiles, core.NewGenFile(info.FilePath, info.Data))
	}

	return genFiles, nil
}
