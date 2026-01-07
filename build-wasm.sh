#!/bin/bash

BUILD_TAG=`date -u +%Y-%m-%dT%H:%M:%S`
NAKAMA_URL=${NAKAMA_URL:-http://127.0.0.1:7350}
NAKAMA_KEY=${NAKAMA_KEY:-defaultkey}

GOOS=js GOARCH=wasm go build \
    -ldflags "\
        -X main.BUILD_TAG=${BUILD_TAG} \
        -X main.NAKAMA_URL=${NAKAMA_URL} \
        -X main.NAKAMA_KEY=${NAKAMA_KEY} \
        -X main.IS_WASM_BUILD=1 \
    " \
    -o ebiman.wasm \
    .

cp $(go env GOROOT)/lib/wasm/wasm_exec.js .
