#!/bin/bash
rm -rf ./target
mkdir -p ./target

CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o ./target/go-read go-read.go
cp -r ./static ./target/static
