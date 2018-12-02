GOFILES=$(wildcard *.go)

default: build

build:
	go build -v -o ./godo  ./client

install-client:
	@GOPATH=$(GOPATH) GOBIN=$(GOPATH)/bin go install client/godo.go

install-server:
	@GOPATH=$(GOPATH) GOBIN=$(GOPATH)/bin go install servers/server.go

install: install-client install-server


#run: build
#	./_bin/snake-game

#test:
#	go test -v ./...
