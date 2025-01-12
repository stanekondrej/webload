#!/usr/bin/bash

protoc -I="$(pwd)" --go_out="$(pwd)/server/pkg/pb" "$(pwd)/messages.proto"
protoc -I="$(pwd)" --go_out="$(pwd)/client/pkg/pb" "$(pwd)/messages.proto"
protoc -I="$(pwd)" --plugin="./frontend/node_modules/.bin/protoc-gen-ts_proto" --ts_proto_out="$(pwd)/frontend/proto" "$(pwd)/messages.proto"
