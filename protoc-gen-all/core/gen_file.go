package core

import (
	"go/format"
	"os"
	"sync"

	"golang.org/x/tools/imports"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/plogging"
)

type GenFile interface {
	Format() error
	CreateOrWrite() error
	GetFilePath() string
}

type genFile struct {
	filePath     string
	newData      []byte
	oldDataCache map[string][]byte
	mu           *sync.Mutex
}

func NewGenFile(filePath string, newData []byte) GenFile {
	return &genFile{
		filePath:     filePath,
		newData:      newData,
		oldDataCache: make(map[string][]byte),
		mu:           &sync.Mutex{},
	}
}

func (g *genFile) Format() error {
	importsData, err := imports.Process("", g.newData, &imports.Options{
		Fragment:   true,
		AllErrors:  false,
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: false,
	})
	if err != nil {
		return perrors.Newf("goimportsでエラーが発生しました。 err = %v, filePath = %s, fileData = \n%s", err, g.GetFilePath(), string(g.newData))
	}

	fmtData, err := format.Source(importsData)
	if err != nil {
		return perrors.Newf("gofmtでエラーが発生しました。 err = %v, filePath = %s, fileData = \n%s", err, g.GetFilePath(), string(g.newData))
	}
	g.newData = fmtData

	return nil
}

func (g *genFile) CreateOrWrite() error {
	path := g.filePath
	file, err := os.Create(path)
	if err != nil {
		return perrors.Stack(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			plogging.GetLogger().Infof("ファイルのクローズに失敗しました。 err = %v, filePath = %s\n", err, path)
		}
	}()

	if _, err := file.Write(g.newData); err != nil {
		return perrors.Stack(err)
	}

	return nil
}

func (g *genFile) GetFilePath() string {
	return g.filePath
}
