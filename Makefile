BIN_NAME := zip-hacker

.PHONY: all test clean pb

all: build

deps:
	go mod tidy

build: deps
	go build -o bin/$(BIN_NAME) .

test:
	go test -v .

fmt:
	gofmt -w ${GOFILES}

clean:
	@rm -rf bin
