GOFILES=$(wildcard *.go)

default: build

build: clean-bin build-client build-server

install-client: build-client
	@GOPATH=$(GOPATH) GOBIN=$(GOPATH)/bin go install ./cmd/client/godo.go

build-server: clean-server
	go build -v -o ./bin/godo-server ./cmd/server

build-client: clean-client
	go build -v -o ./bin/godo ./cmd/client

install: clean-bin install-client build-server

clean-bin:
	@mkdir -p bin
	@rm -f bin/*

clean-server:
	@mkdir -p bin
	@rm -f bin/godo-server*

clean-client:
	@mkdir -p bin
	@rm -f bin/godo

clean: clean-bin

#run: build
#	./_bin/snake-game

#test:
#	go test -v ./...
