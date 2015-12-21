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

release:
	rice embed-go -i ./internal/render/
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -ldflags "-s $(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-darwin-386 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-darwin-amd64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-linux-386 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-linux-amd64 main.go
	zip -m -9 $(BINARY)-darwin-386.zip $(BINARY)-darwin-386
	zip -m -9 $(BINARY)-darwin-amd64.zip $(BINARY)-darwin-amd64
	zip -m -9 $(BINARY)-linux-386.zip $(BINARY)-linux-386
	zip -m -9 $(BINARY)-linux-amd64.zip $(BINARY)-linux-amd64
