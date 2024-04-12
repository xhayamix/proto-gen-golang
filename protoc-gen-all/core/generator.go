package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/parallel"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

type Generator interface {
	Build() ([]GenFile, error)

	SetGenFiles(genFiles []GenFile)
	Format() error
	Generate() error
	GetGeneratedFilePaths() []string
}

type GeneratorBase struct {
	genFiles []GenFile
}

func NewGeneratorBase() *GeneratorBase {
	return &GeneratorBase{
		genFiles: make([]GenFile, 0),
	}
}

type emptyGenerator struct {
	*GeneratorBase
}

func (e *emptyGenerator) Build() ([]GenFile, error) {
	return nil, nil
}

func NewEmptyGenerator() Generator {
	return &emptyGenerator{
		GeneratorBase: NewGeneratorBase(),
	}
}

func (g *GeneratorBase) SetGenFiles(genFiles []GenFile) {
	g.genFiles = genFiles
}

func (g *GeneratorBase) Format() error {
	pg, ctx := parallel.NewGroupWithContext(context.Background(), parallel.DefaultSize)

	for _, file := range g.genFiles {
		file := file

		if !strings.HasSuffix(file.GetFilePath(), ".gen.go") {
			continue
		}

		pg.Go(ctx, func(_ context.Context) error {
			before := time.Now()

			if err := file.Format(); err != nil {
				return perrors.Stack(err)
			}

			after := time.Since(before)
			if after.Seconds() > 3 {
				_, _ = fmt.Fprintf(os.Stderr, "file format trace -> file: %s, time: %d ms\n", file.GetFilePath(), after.Milliseconds())
			}
			return nil
		})
	}
	if err := pg.Wait(); err != nil {
		return perrors.Stack(err)
	}

	return nil
}

func (g *GeneratorBase) Generate() error {
	pg, ctx := parallel.NewGroupWithContext(context.Background(), parallel.DefaultSize)

	for _, file := range g.genFiles {
		file := file

		pg.Go(ctx, func(_ context.Context) error {
			outputDir := filepath.Dir(file.GetFilePath())
			if err := os.MkdirAll(outputDir, 0777); err != nil {
				return perrors.Stack(err)
			}

			if err := file.CreateOrWrite(); err != nil {
				return perrors.Stack(err)
			}

			return nil
		})
	}
	if err := pg.Wait(); err != nil {
		return perrors.Stack(err)
	}

	return nil
}

func (g *GeneratorBase) GetGeneratedFilePaths() []string {
	paths := make([]string, 0, len(g.genFiles))

	for _, file := range g.genFiles {
		paths = append(paths, file.GetFilePath())
	}

	return paths
}
