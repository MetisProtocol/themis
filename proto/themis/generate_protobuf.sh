#!/bin/sh

DIR=$(cd "$(dirname "$0")/" && pwd)

PROTO_FILES=$(find "$DIR" -iname "*.proto")

echo $DIR
echo $PROTO_FILES

# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

PATH="$(go env GOPATH)/bin:$PATH"

# generate golang pb & grpc
protoc -I "$DIR" \
    --go_out="$DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$DIR" --go-grpc_opt=paths=source_relative \
    $PROTO_FILES