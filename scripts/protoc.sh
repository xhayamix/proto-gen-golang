#!/bin/sh

proto_file_dir=./proto
validate_proto_file_dir=${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate\@v1.0.4
bq_schema_proto_file_dir=${GOPATH}/pkg/mod/github.com/!google!cloud!platform/protoc-gen-bq-schema\@v0.0.0-20230915083002-8edab4bb6c81
def_proto_dir=${proto_file_dir}/definition
client_proto_dir=${proto_file_dir}/client
server_proto_dir=${proto_file_dir}/server
log_proto_dir=${proto_file_dir}/log

out_proto_dir=./pkg/domain/proto
mkdir -p ${out_proto_dir}

call_protoc() {
  protoc "$@" || exit $?
}

# 1. optionsから生成 # options配下のprotoから自動生成したい時だけコメントアウト
 def_option_proto_files=$(find ${server_proto_dir}/options -type f -name '*.proto')
 call_protoc \
   --proto_path=${proto_file_dir} \
  	--go_out=paths=source_relative:${out_proto_dir} \
   ${def_option_proto_files}

# 2. definitionからEnumのProtoを生成
def_enum_proto_files=$(find ${def_proto_dir}/enums -type f -name '*.proto')
call_protoc \
  --proto_path=${proto_file_dir} \
	--all_out=gen_enum,paths=source_relative:. \
  ${def_enum_proto_files}

# 3. 2の生成物をもとにEnum関連のファイルを生成
server_enums_proto_files=$(find ${server_proto_dir}/enums -type f -name '*.proto')
client_enums_proto_files=$(find ${client_proto_dir}/enums -type f -name '*.proto')
call_protoc \
  --proto_path=${proto_file_dir} \
	--go_out=paths=source_relative:${out_proto_dir} \
  ${server_enums_proto_files} ${client_enums_proto_files}


# 4. 3の生成物をもとにサーバ関連のファイルを生成
#server_common_proto_files=$(find ${server_proto_dir}/common -type f -name '*.proto')
server_api_proto_files=$(find ${server_proto_dir}/api -type f -name '*.proto')
server_transaction_proto_files=$(find ${server_proto_dir}/transaction -type f -name '*.proto')
call_protoc \
  --proto_path=${validate_proto_file_dir} \
  --proto_path=${proto_file_dir} \
	--all_out=gen_api,gen_transaction,paths=source_relative:. \
  ${server_api_proto_files} ${server_transaction_proto_files} ${server_enums_proto_files}


# 5. 4の生成物をもとにprotobuf、grpc、validatorの実装を生成
client_api_proto_files=$(find ${client_proto_dir}/api -type f -name '*.proto' | sort)
call_protoc \
  --proto_path=${validate_proto_file_dir} \
  --proto_path=${proto_file_dir} \
	--go_out=paths=source_relative:${out_proto_dir} \
	--go-grpc_out=require_unimplemented_servers=false,paths=source_relative:${out_proto_dir} \
	--validate_out=lang=go,paths=source_relative:${out_proto_dir} \
  ${client_api_proto_files}
