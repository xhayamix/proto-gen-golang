package generator

/**
package example

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator"
)

type generator struct {
	*generator.GeneratorBase
	plugin *protogen.Plugin
}

func NewGenerator(plugin *protogen.Plugin) generator.Generator {
	return &generator{
		GeneratorBase: generator.NewGeneratorBase(),
		plugin:        plugin,
	}
}

func (g *generator) Build() ([]generator.GenFile, error) {
	genFiles := make([]generator.GenFile, 0)

	_, _ = generator.NewGenFile("", []byte(""))

	return genFiles, nil
}
*/
