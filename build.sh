#!/bin/bash
echo "start building version: ${BUILD_NUMBER}"
rm -rf ./target
mkdir -p ./target

CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o ./target/go-read go-read.go
echo "build output code: $?"

cp -r ./static ./target/static
echo "static copying output code: $?"
