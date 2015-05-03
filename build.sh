#!/bin/bash
CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' go-read.go
