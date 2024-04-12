package main

import (
	"strings"
	"time"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/enum"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/plogging"
)

type flagKind string

const (
	flagKindGenMaster       flagKind = "gen_master"
	flagKindGenTransaction  flagKind = "gen_transaction"
	flagKindGenZap          flagKind = "gen_zap"
	flagKindGenEnum         flagKind = "gen_enum"
	flagKindGenApi          flagKind = "gen_api"
	flagKindGenAdmin        flagKind = "gen_admin"
	flagKindGenShare        flagKind = "gen_share"
	flagKindGenCommon       flagKind = "gen_common"
	flagKindGenLog          flagKind = "gen_log"
	flagKindGenLogActionGen flagKind = "gen_log_action_gen"
	flagKindGenRanking      flagKind = "gen_ranking"
	flagKindGenCache        flagKind = "gen_cache"
	flagKindWritePb         flagKind = "write_pb"
)

func main() {
	locationName := "Asia/Tokyo"
	location, err := time.LoadLocation(locationName)
	if err != nil {
		location = time.FixedZone(locationName, 9*60*60)
	}
	time.Local = location

	startTime := time.Now()
	logger := plogging.GetLogger()
	logger.Infof("protoc-gen-all start\n")

	generatorBuilder := core.NewGeneratorBuilder()

	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		generatorMap, _ := createGeneratorMap(plugin)
		kinds := make([]string, 0, len(generatorMap))
		for kind, generator := range generatorMap {
			kinds = append(kinds, string(kind))
			generatorBuilder.AppendGenerator(generator)
		}
		logger.Infof("flag %s\n", strings.Join(kinds, ","))

		generatedFilenamePrefixList := make([]string, 0, len(plugin.Files))
		/*
			if writePb {
				for _, file := range plugin.Files {
					if !file.Generate {
						continue
					}
					generatedFilenamePrefixList = append(generatedFilenamePrefixList, file.GeneratedFilenamePrefix)
				}
			}
		*/
		if err := generatorBuilder.Generate(generatedFilenamePrefixList); err != nil {
			return perrors.Stack(err)
		}
		return nil
	})

	endTime := time.Now()
	logger.Infof("protoc-gen-campus end, elapsed: %s\n", endTime.Sub(startTime).String())
}

func createGeneratorMap(plugin *protogen.Plugin) (map[flagKind]core.Generator, bool) {
	generatorMap := make(map[flagKind]core.Generator)
	var writePb bool

	for _, param := range strings.Split(plugin.Request.GetParameter(), ",") {
		s := strings.Split(param, "=")

		switch flagKind(s[0]) {
		case flagKindGenApi:
			generatorMap[flagKindGenApi] = api.NewGenerator(plugin)
		case flagKindGenEnum:
			generatorMap[flagKindGenEnum] = enum.NewGenerator(plugin)
			/*
				case flagKindWritePb:
					generatorMap[flagKindWritePb] = core.NewEmptyGenerator()
					writePb = true

			*/
		default:
			continue
		}
	}

	return generatorMap, writePb
}
