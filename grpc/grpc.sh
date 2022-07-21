#!/bin/bash

set -e

dir=$(cd $(dirname $BASH_SOURCE) && pwd)
output="$dir/protocol"
document="$dir/document"
if [[ -d "$output" ]];then
    rm "$output/*" -rf
else
    mkdir "$output"
fi
if [[ -d "$document" ]];then
    rm "$document/*" -rf
else
    mkdir "$document"
fi

# generate grpc code
protoc -I "$dir/pb" -I "$dir/third_party/googleapis" \
    --go_out=./protocol --go_opt=paths=source_relative \
    --go-grpc_out="$output" --go-grpc_opt=paths=source_relative \
    math/math.proto

# generate gateway code
protoc -I "$dir/pb" -I "$dir/third_party/googleapis" \
    --go_out=./protocol --go_opt=paths=source_relative \
    --grpc-gateway_out="$output" --grpc-gateway_opt=paths=source_relative \
    math/math.proto

protoc -I "$dir/pb" -I "$dir/third_party/googleapis" \
    --go_out=./protocol --go_opt=paths=source_relative \
    --openapiv2_out="$document"  \
    math/math.proto
    