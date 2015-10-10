#!/usr/bin/env bash
VERSION="1.0"
BINARY="./bin/greedy-${VERSION}"

for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    godep go build -ldflags "-X main.BuildDate=`date +"%d-%B-%Y/%T"`" -a -tags netgo -installsuffix netgo -o ${BINARY}-$GOOS-$GOARCH main.go
  done
done
