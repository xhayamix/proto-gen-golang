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

WORKDIR ${GOPATH}
