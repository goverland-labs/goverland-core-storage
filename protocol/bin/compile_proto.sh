#!/bin/sh

# Format of Using:
#   sh bin/compile_proto.sh PROTO_SRC_PATH OUT_ROOT_PATH PROTO_FILES
# Example:
#   sh bin/compile_proto.sh proto proto/gen proto/*/*.proto

# exit when any command fails
set -e

# create directory if not exists
mkdir -p $2

# remove previously generated .pb.go files
find $2 -type f -name "*.pb.go" | xargs -r -L1 rm

protoc --proto_path=$1 --go_out=$2 --go-grpc_out=$2 $3 --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

echo "Files '$3' were compiled"
