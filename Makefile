default: run

VERSION := 1.0
LDFLAGS := -X github.com/rogierlommers/greedy/internal/common.CommitHash=`git rev-parse HEAD` -X github.com/rogierlommers/greedy/internal/common.BuildDate=`date +"%d-%B-%Y/%T"`
BINARY := ./bin/greedy-${VERSION}

build:
	mkdir -p ./target
	rice embed-go -i ./internal/render/
	CGO_ENABLED=0 go build -ldflags "-s $(LDFLAGS)" -a -installsuffix cgo -o ./target/greedy main.go

run:
	go run *.go
