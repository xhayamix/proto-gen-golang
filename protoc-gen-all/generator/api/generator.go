package api

import (
	_ "embed"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/gen-proto/server/api"
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	interactorapi "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/admin/usecase/interactor/api"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/api/middleware/check"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/api/middleware/responsecache"
	apiregistry "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/api/registry"
	portapi "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/domain/port/api"
	clientapi "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/domain/proto/client/api"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/infra/genserverapi"
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
	&api.Creator{},
	&clientapi.Creator{},
}

var allCreators = []output.TemplateCreator{
	&interactorapi.Creator{},
	&portapi.Creator{},
	&responsecache.Creator{},
	apiregistry.New(),
}

var allTemplatesCreators = []output.TemplatesCreator{
	&check.Creator{},
	&genserverapi.Creator{},
}

func (g *generator) Build() ([]core.GenFile, error) {
	files := make([]*input.File, 0)

	for _, f := range g.plugin.Files {
		if !f.Generate {
			continue
		}
		if !strings.HasPrefix(f.Proto.GetPackage(), "server.api") {
			continue
		}

		file, err := input.ConvertFileFromProto(f)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		if file == nil {
			continue
		}

		files = append(files, file)
	}

	// 入力ファイルの順番に左右されないようソートする
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].SnakeName < files[j].SnakeName
	})

	genFiles := make([]core.GenFile, 0)

	for _, creator := range eachCreators {
		for _, file := range files {
			info, err := creator.Create(file)
			if err != nil {
				return nil, perrors.Stack(err)
			}
			if info == nil {
				continue
			}

			genFiles = append(genFiles, core.NewGenFile(info.FilePath, info.Data))
		}
	}

	for _, creator := range allCreators {
		info, err := creator.Create(files)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		if info == nil {
			continue
		}
		genFiles = append(genFiles, core.NewGenFile(info.FilePath, info.Data))
	}

	for _, creator := range allTemplatesCreators {
		infos, err := creator.Create(files)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		for _, info := range infos {
			if infos == nil {
				continue
			}
			genFiles = append(genFiles, core.NewGenFile(info.FilePath, info.Data))
		}
	}

	return genFiles, nil
}
