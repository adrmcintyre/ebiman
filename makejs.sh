#!/bin/bash

# build
env GOOS=js GOARCH=wasm go build \
    -o pacman.wasm \
	"github.com/adrmcintyre/poweraid/game"

# grab js
cp $(go env GOROOT)/lib/wasm/wasm_exec.js .


