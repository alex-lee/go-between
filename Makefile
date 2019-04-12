# Makefile for go-between udp proxy

.PHONY: all deps cmd racecmd test coverage

all: cmd

deps:
	go mod tidy

cmd:
	CGO_ENABLED=0 go build github.com/alex-lee/go-between/cmd/go-between

racecmd:
	go build -race github.com/alex-lee/go-between/cmd/go-between

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
