.PHONY: build

default: build

run:
	go run ./cmd/web

format:
	gofmt -w -s .

ci_lint:
	golangci-lint run ./... --fix

linter: format ci_lint

tests:
	go test -v ./... | { grep -v 'no test files'; true; }

build: run