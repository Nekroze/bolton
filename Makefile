## Variables

NAME := bolton

GO_FILES := $(shell find . -name '*.go' -type f)
PROJ_DIR := $(realpath .)
GOPATH   := $(shell echo "${GOPATH}" | cut -d':' -f1 )
BINARY   := $(GOPATH)/bin/$(NAME)
GORICH   := $(GOPATH)/bin/richgo
GOCILINT := $(GOPATH)/bin/golangci-lint

GO111MODULE = on

## Phony Targets

.PHONY: binary dev deps depunpin depupdate test lint sanity

binary: $(BINARY)

dev:
	find . -type f | grep -E '(\.(go|mod|sum)|Makefile)$$' | entr -d make

deps: go.sum
	go mod download

depunpin:
	rm -f go.sum

depupdate: depunpin deps

run: $(BINARY)
	$<

test:
	go test -v ./...

lint: $(GOCILINT)
	$< run --enable-all --disable gochecknoglobals,scopelint,dupl

sanity: binary lint test

## Targets

$(GORICH):
	env GO111MODULE=off go get github.com/kyoh86/richgo

$(GOCILINT):
	env GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint

go.sum: go.mod
	go get -d -u -v ./...
	touch $@ $<

$(BINARY): deps $(GO_FILES)
	go build -o $@ .
