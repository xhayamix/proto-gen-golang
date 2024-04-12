package api

import (
	_ "embed"
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/gen-proto/server/api/game"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/admin/handler"
	adminregistry "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/admin/registry"
	interactorgameapi "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/admin/usecase/interactor/gameapi"
	adminiamrole "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/admin/usecase/interactor/iamrole"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/api/middleware/check"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/api/middleware/responsecache"
	apiregistry "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/api/registry"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/gateway"
	loadtestregistry "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/cmd/loadtest/registry"
	portgameapi "github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/domain/port/gameapi"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/domain/proto/client/api"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/pkg/infra/campusserverapi"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output/web/src/apps/admin/requests"
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

var eachGameCreators = []output.EachTemplateCreator{
	&game.Creator{},
	&api.Creator{},
}

var eachAdminCreators = []output.EachTemplateCreator{
	handler.New(),
	requests.New(),
}

var allGameCreators = []output.TemplateCreator{
	&interactorgameapi.Creator{},
	&portgameapi.Creator{},
	&responsecache.Creator{},
	apiregistry.New(),
	gateway.New(),
	loadtestregistry.New(),
}

var allAdminCreators = []output.TemplateCreator{
	adminregistry.New(),
	adminiamrole.New(),
}

var allGameTemplatesCreators = []output.TemplatesCreator{
	&check.Creator{},
	&campusserverapi.Creator{},
}

func (g *generator) Build() ([]core.GenFile, error) {
	gameFiles := make([]*input.File, 0)
	adminFiles := make([]*input.File, 0)

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

		if strings.HasPrefix(f.Proto.GetPackage(), "server.api.game") {
			gameFiles = append(gameFiles, file)
		} else if strings.HasPrefix(f.Proto.GetPackage(), "server.api.admin") {
			adminFiles = append(adminFiles, file)
		}
	}

	// 入力ファイルの順番に左右されないようソートする
	sort.SliceStable(gameFiles, func(i, j int) bool {
		return gameFiles[i].SnakeName < gameFiles[j].SnakeName
	})
	sort.SliceStable(adminFiles, func(i, j int) bool {
		return adminFiles[i].SnakeName < adminFiles[j].SnakeName
	})

	genFiles := make([]core.GenFile, 0)

	for _, creator := range eachGameCreators {
		for _, file := range gameFiles {
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

	for _, creator := range eachAdminCreators {
		for _, file := range adminFiles {
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

	for _, creator := range allGameCreators {
		info, err := creator.Create(gameFiles)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		if info == nil {
			continue
		}
		genFiles = append(genFiles, core.NewGenFile(info.FilePath, info.Data))
	}

	for _, creator := range allAdminCreators {
		info, err := creator.Create(adminFiles)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		if info == nil {
			continue
		}
		genFiles = append(genFiles, core.NewGenFile(info.FilePath, info.Data))
	}

	for _, creator := range allGameTemplatesCreators {
		infos, err := creator.Create(gameFiles)
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
