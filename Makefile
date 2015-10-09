default: run

setup:
	go get github.com/tools/godep

build: setup
	rm -rf ./target
	mkdir -p ./target
	godep go build -ldflags "-X main.BuildDate=`date +"%d-%B-%Y/%T"`" -a -tags netgo -installsuffix netgo -o ./target/greedy main.go

run:
	godep go run *.go
