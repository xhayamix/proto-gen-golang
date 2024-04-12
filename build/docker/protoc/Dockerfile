ARG GO_VERSION
FROM golang:${GO_VERSION}
ARG PROTOBUF_VERSION=26.1

RUN apt-get update && apt-get install unzip

WORKDIR /tmp/protoc
RUN curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip -o protoc.zip
RUN unzip protoc.zip
RUN mv bin/* /usr/local/bin/
RUN mv include/* /usr/local/include/

WORKDIR /usr/src/proto-gen-golang
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN make -C ${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate\@v1.0.4/ build
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

COPY pkg/domain/proto/client/enums pkg/domain/proto/client/enums
COPY pkg/domain/proto/client/options pkg/domain/proto/client/options
COPY pkg/domain/proto/server/options pkg/domain/proto/server/options
COPY pkg/domain/proto/definition/options pkg/domain/proto/definition/options
COPY protoc-gen-all protoc-gen-all
#COPY quapay-server/db/ddl/quapay_spanner.up.sql protoc-gen-all/generator/transaction/output/db/migrations/spanner/quapay_spanner.up.sql
RUN go install ./protoc-gen-all


WORKDIR ${GOPATH}

ENTRYPOINT ["/usr/local/bin/protoc"]
