SHELL := /bin/bash

TAG = dev
IMG ?= opa-bundle-api:$(TAG)
TEST_ENV_FILE = tmp/test_env
VERSION ?= "v0.0.0-dev"
REVISION ?= ""
CREATED ?= ""

ifneq (,$(wildcard $(TEST_ENV_FILE)))
    include $(TEST_ENV_FILE)
    export
endif

.PHONY: all
.SILENT: all
all: tidy lint fmt vet gosec test build

.PHONY: lint
.SILENT: lint
lint:
	golangci-lint run

.PHONY: fmt
.SILENT: fmt
fmt:
	go fmt ./...

.PHONY: tidy
.SILENT: tidy
tidy:
	go mod tidy

.PHONY: vet
.SILENT: vet
vet:
	go vet ./...

.PHONY: test
.SILENT: test
test:
	mkdir -p tmp
	go test -timeout 1m ./... -cover

.PHONY: gosec
.SILENT: gosec
gosec:
	gosec ./...

.PHONY: cover
.SILENT: cover
cover:
	go test -timeout 1m ./... -coverprofile=tmp/coverage.out                                                                                                                                                                                         16:10:38
	go tool cover -html=tmp/coverage.out	

.PHONY: run
.SILENT: run
run:
	go run cmd/opa-bundle-api/main.go

.PHONY: gen-docs
.SILENT: gen-docs
gen-docs:
	go run cmd/gen-docs/main.go

.PHONY: build
.SILENT: build
build:
	go build -ldflags "-w -s -X main.Version=$(VERSION) -X main.Revision=$(REVISION) -X main.Created=$(CREATED)" -o bin/opa-bundle-api cmd/opa-bundle-api/main.go

.PHONY: opa-eval
.SILENT: opa-eval
opa-eval:
	opa eval --data pkg/bundle/static/rule.rego --data test/opa/data.json --input test/opa/input.json --format pretty 'data.rule'

.PHONY: opa-run
.SILENT: opa-run
opa-run:
	opa run --server --addr :8181 --config-file ./test/opa/config.yaml