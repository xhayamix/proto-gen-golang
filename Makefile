MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
include ${MAKEFILE_DIR}.env
GO_VERSION := $(shell cat ${MAKEFILE_DIR}.go-version)
NODE_VERSION := $(shell cat ${MAKEFILE_DIR}web/.node-version)

define docker-compose
	GO_VERSION=${GO_VERSION} \
	NODE_VERSION=${NODE_VERSION} \
	docker-compose \
		-f ${MAKEFILE_DIR}development/docker-compose.yml \
		--env-file ${MAKEFILE_DIR}.env \
		-p campus-server \
		$1
endef

## protobuf (usage: `make protobuf arg='protoc --proto_path=./proto --go_out=module=github.com/xhayamix/proto-gen-golang:. ./proto/server/options/api/api.proto'`)
.PHONY: protobuf
protobuf:
	$(call docker-compose, run --rm protobuf $(arg))

## protoc
.PHONY: protoc
protoc:
	$(call docker-compose, run --rm --entrypoint sh protoc ./scripts/protoc.sh)
	goimports -w -local "github.com/QualiArts/campus-server" pkg/domain/proto/client pkg/domain/proto/server pkg/cmd/admin/handler pkg/domain/proto/definition
	gofmt -s -w pkg/domain/proto/client pkg/cmd/admin/handler pkg/domain/proto/definition
	$(call clang-format, --entrypoint sh, -c "find ./campus-proto -type f -name '*.proto' | xargs clang-format -i")
	npm run format --prefix web

## local install
.PHONY: local-install
local-install:
	# 開発する上で必要なものをinstall
	go install golang.org/x/tools/cmd/goimports
	go install go.uber.org/mock/mockgen
	go install github.com/vektra/mockery/v2@v2.38.0
	go install github.com/googleapis/api-linter/cmd/api-linter@v1.32.3
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

