MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

SHELL := $(shell if command -v bash >/dev/null 2>&1; then echo bash; else echo sh; fi)
ifeq ($(SHELL),bash)
.SHELLFLAGS := -euo pipefail -c
else
.SHELLFLAGS := -c
endif

SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

GO ?= go

MODULE ?= $(shell go list -m)
NAME ?= $(notdir $(MODULE))

.PHONY: run
run:
	$(GO) run cmd/main.go run

.PHONY: build
build:
	$(GO) build cmd/main.go

.PHONY: test
test:
	$(GO) test -race -cover ./...

.PHONY: bench
bench:
	$(GO) test -run=XXXXXX -benchtime=10s -bench=./ || exit 1

.PHONY: update
update:
	$(GO) get -u ./...

.PHONY: deps-update
deps-update:
	$(GO) get -u -t -v ./...
	$(GO) mod tidy

.PHONY: deps-cleancache
deps-cleancache:
	$(GO) clean -modcache

.PHONY: format-proto
format-proto:
	buf format -w

.PHONY: lint-proto
lint-proto:
	buf lint

.PHONY: tidy
tidy: format-proto
	$(GO) fmt ./...
	$(GO) mod tidy
	$(GO) mod verify
	goimports -w .
	golines -m 120 -w --ignore-generated .
	gci write --skip-generated -s standard -s "prefix($(MODULE))" -s default -s blank -s dot --custom-order  .
	gofumpt -l -w .

## quality/lint: Run linters
.PHONY: lint
lint: lint-proto
	revive -config revive-config.toml -formatter friendly ./...
	golangci-lint run ./...

## quality/audit: Run quality control checks
.PHONY: audit
audit:
	$(GO) mod verify
	$(GO) vet ./...
	staticcheck -checks=all,-ST1000,-U1000 ./...
	govulncheck ./...
	$(GO) test -race ./...

## quality/vet: Run $(GO) vet
.PHONY: vet
vet:
	$(GO) vet ./...

