package enum

import (
	_ "embed"
	"sort"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output"
	clientprotoenum "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output/gen-proto/client/enums"
	clientsettingproto "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output/gen-proto/client/master"
	serverprotoenum "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output/gen-proto/server/enums"
	serversettingproto "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output/gen-proto/server/master"
	settingcpackersetter "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output/pkg/domain/component/masterconverter"
	settingmcache "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output/pkg/domain/entity/mcache"
	golangenum "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output/pkg/domain/enum"
	settingvalidator "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output/pkg/domain/service/validation/validator"
	tsenum "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum/output/web/src/enums"
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
	&golangenum.Creator{},
	&settingmcache.Creator{},
	&settingcpackersetter.Creator{},
	&settingvalidator.Creator{},
	&clientsettingproto.Creator{},
	&serversettingproto.Creator{},
	&tsenum.EachCreator{},
}

var bulkCreators = []output.BulkTemplateCreator{
	&clientprotoenum.Creator{},
	&serverprotoenum.Creator{},
	&golangenum.MapCreator{},
	&tsenum.Creator{},
}

func (g *generator) Build() ([]core.GenFile, error) {
	enums := make([]*input.Enum, 0)

	for _, f := range g.plugin.Files {
		if !f.Generate {
			continue
		}
		if f.Proto.GetPackage() != "definition.enums" {
			continue
		}

		enum, err := input.ConvertMessageFromProto(f)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		enums = append(enums, enum)
	}

	// 入力ファイルの順番に左右されないようソートする
	sort.SliceStable(enums, func(i, j int) bool {
		return enums[i].SnakeName < enums[j].SnakeName
	})

	genFiles := make([]core.GenFile, 0)
	for _, creator := range eachCreators {
		for _, enum := range enums {
			info, err := creator.Create(enum)
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
		info, err := creator.Create(enums)
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
