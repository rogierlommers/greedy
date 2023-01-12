#!/bin/bash

# config
VERSION=`git rev-parse HEAD`
LDFLAGS="-X github.com/rogierlommers/greedy/internal/common.CommitHash=`git rev-parse HEAD` -X github.com/rogierlommers/greedy/internal/common.BuildDate=`date +'%d-%B-%Y/%T'`"
BINARY="./target/greedy-${VERSION}"

# start build
mkdir -p ./target
rice embed-go -i ./internal/render/
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s ${LDFLAGS}" -a -installsuffix cgo -o ./bin/greedy main.go
