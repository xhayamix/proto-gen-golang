package core

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/parallel"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/plogging"
)

type GeneratorBuilder interface {
	AppendGenerator(generator Generator) GeneratorBuilder
	Generate(generatedFilenamePrefixList []string) error
}

type generatorBuilder struct {
	generators []Generator
}

func NewGeneratorBuilder() GeneratorBuilder {
	return &generatorBuilder{
		generators: make([]Generator, 0),
	}
}

func (g *generatorBuilder) AppendGenerator(generator Generator) GeneratorBuilder {
	g.generators = append(g.generators, generator)

	return g
}

func (g *generatorBuilder) Generate(generatedFilenamePrefixList []string) error {
	pg, ctx := parallel.NewGroupWithContext(context.Background(), parallel.DefaultSize)

	mu := &sync.Mutex{}
	pathSet := strset.New()

	for _, generator := range g.generators {
		generator := generator

		pg.Go(ctx, func(_ context.Context) error {
			genFileDirectories, err := generator.Build()
			if err != nil {
				return perrors.Stack(err)
			}
			generator.SetGenFiles(genFileDirectories)

			if err := generator.Format(); err != nil {
				return perrors.Stack(err)
			}

			if err := generator.Generate(); err != nil {
				return perrors.Stack(err)
			}

			mu.Lock()
			pathSet.Add(generator.GetGeneratedFilePaths()...)
			mu.Unlock()

			return nil
		})
	}
	if err := pg.Wait(); err != nil {
		return perrors.Stack(err)
	}

	// pb.go系の差分削除
	/*
		for _, fileNamePrefix := range generatedFilenamePrefixList {
			splitFileNamePrefix := strings.Split(fileNamePrefix, "/")
			fileName := splitFileNamePrefix[len(splitFileNamePrefix)-1]
			pkgName := strings.ReplaceAll(fileName, "_", "")

			pathSet.Add(
				// その他
				"pkg/grpc/codec/protoenc/testdata/test.pb.go",

				// option系
				"pkg/domain/proto/client/master/common/options.pb.go",
				"pkg/domain/proto/client/options/check_option.pb.go",
				"pkg/domain/proto/client/options/zap.pb.go",
				"pkg/domain/proto/client/transaction/common/options.pb.go",
				"pkg/domain/proto/definition/options/enums/enums.pb.go",
				"pkg/domain/proto/server/options/admin/admin.pb.go",
				"pkg/domain/proto/server/options/api/admin/admin.pb.go",
				"pkg/domain/proto/server/options/api/game/game.pb.go",
				"pkg/domain/proto/server/options/cache/cache.pb.go",
				"pkg/domain/proto/server/options/common/common.pb.go",
				"pkg/domain/proto/server/options/log/log.pb.go",
				"pkg/domain/proto/server/options/master/master.pb.go",
				"pkg/domain/proto/server/options/ranking/ranking.pb.go",
				"pkg/domain/proto/server/options/transaction/transaction.pb.go",
				"pkg/domain/proto/server/options/transaction/transaction.pb.go",
				"pkg/domain/proto/server/options/zap/zap.pb.go",

				// pkg/cmd/admin/handler
				"pkg/cmd/admin/handler/"+pkgName+"/"+fileName+".pb.go",
				"pkg/cmd/admin/handler/"+pkgName+"/"+fileName+".pb.validate.go",

				// pkg/domain/proto/client/api/
				"pkg/domain/proto/client/api/"+fileName+".pb.common_response.gen.go",
				"pkg/domain/proto/client/api/"+fileName+".pb.go",
				"pkg/domain/proto/client/api/"+fileName+".pb.gw.go",
				"pkg/domain/proto/client/api/"+fileName+".pb.validate.go",
				"pkg/domain/proto/client/api/"+fileName+".pb.zap.gen.go",
				"pkg/domain/proto/client/api/"+fileName+"_grpc.pb.go",

				// pkg/domain/proto/client/api/common/
				"pkg/domain/proto/client/api/common/"+fileName+".pb.go",
				"pkg/domain/proto/client/api/common/"+fileName+".pb.validate.go",
				"pkg/domain/proto/client/api/common/"+fileName+".pb.zap.gen.go",

				// pkg/domain/proto/client/common/
				"pkg/domain/proto/client/common/"+fileName+".pb.go",
				"pkg/domain/proto/client/common/"+fileName+".pb.zap.gen.go",

				// pkg/domain/proto/client/enums/
				"pkg/domain/proto/client/enums/"+fileName+".pb.go",

				// pkg/domain/proto/client/exam/
				"pkg/domain/proto/client/exam/"+fileName+".pb.go",
				"pkg/domain/proto/client/exam/"+fileName+"_grpc.pb.go",

				// pkg/domain/proto/client/master/
				"pkg/domain/proto/client/master/"+fileName+".pb.go",

				// pkg/domain/proto/client/transaction/
				"pkg/domain/proto/client/transaction/"+fileName+".pb.go",
				"pkg/domain/proto/client/transaction/"+fileName+".pb.zap.gen.go",

				// pkg/domain/proto/server/enums/
				"pkg/domain/proto/server/enums/"+fileName+".pb.go",
			)
		}
	*/

	data := "\n" + strings.Join(pathSet.List(), "\n") + "\n"

	filePath := "/tmp/generated_files.txt"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return perrors.Stack(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			plogging.GetLogger().Infof("ファイルのクローズに失敗しました。 err = %v, filePath = %s\n", err, filePath)
		}
	}()

	if _, err = file.WriteString(data); err != nil {
		return perrors.Stack(err)
	}

	return nil
}
