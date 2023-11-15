.PHONY: all setup dunecli build lint yamllint

LOCALBIN ?= $(shell pwd)/bin
GOLANGCI_LINT_VERSION ?= v1.55.2

all: lint test build

setup: bin/golangci-lint
	go mod download

cli: lint
	go build -o cli cmd/main.go

build: dunecli

bin/golangci-lint:
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

lint: bin/golangci-lint
	go fmt ./...
	go vet ./...
	bin/golangci-lint -c .golangci.yml run ./...
	go mod tidy

test:
	go mod tidy
	go test -timeout=10s -race -benchmem ./...
