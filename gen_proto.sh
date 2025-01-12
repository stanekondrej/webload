#!/usr/bin/bash

protoc -I="$(pwd)" --go_out="$(pwd)/server/pkg/pb" "$(pwd)/messages.proto"
protoc -I="$(pwd)" --go_out="$(pwd)/client/pkg/pb" "$(pwd)/messages.proto"
