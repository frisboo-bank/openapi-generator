MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

SHELL := $(shell if command -v bash >/dev/null 2>&1; then echo bash; else echo sh; fi)
ifeq ($(SHELL),bash)
.SHELLFLAGS := -euo pipefail -c
else
.SHELLFLAGS := -c
endif

# Sub-project directories. The Go module lives in api/ and the React app
# in frontend/; the repo root itself has no go.mod.
API_DIR := api
FRONTEND_DIR := frontend

# Prefix that runs a command inside a sub-project directory. Each recipe
# line is its own shell, so the cd only applies to that line.
API := cd $(API_DIR) &&
FRONT := cd $(FRONTEND_DIR) &&

GO ?= go
PNPM ?= pnpm

MODULE ?= $(shell go -C $(API_DIR) list -m)
NAME := $(notdir $(MODULE))

# Default goal is a build, not the foreground server (run).
.DEFAULT_GOAL := all

# ---------------------------------------------------------------------------
# Backend (Go) -- every go/buf/lint tool runs against the api/ module
# ---------------------------------------------------------------------------

.PHONY: generate-backend
generate-backend:  ## Generate enum sources (installs pinned goenums, runs go generate)
	$(API) go install github.com/zarldev/goenums@v0.4.3
	$(API) PATH="$$(go env GOPATH)/bin:$$PATH" go generate ./...

.PHONY: run
run: generate-backend  ## Run the API server
	$(API) $(GO) run cmd/main.go run

.PHONY: build-backend
build-backend: generate-backend  ## Build the API binary
	$(API) $(GO) build cmd/main.go

.PHONY: test-backend
test-backend: generate-backend  ## Run go tests (race + coverage)
	$(API) $(GO) test -race -cover ./...

.PHONY: bench-backend
bench-backend: generate-backend  ## Run go benchmarks
	$(API) $(GO) test -run=XXXXXX -benchtime=10s -bench=./ || exit 1

.PHONY: update-backend
update-backend:  ## Bump go module dependencies
	$(API) $(GO) get -u ./...

.PHONY: deps-update-backend
deps-update-backend:  ## Update + tidy go dependencies
	$(API) $(GO) get -u -t -v ./...
	$(API) $(GO) mod tidy

.PHONY: deps-cleancache-backend
deps-cleancache-backend:  ## Clear the go module cache
	$(API) $(GO) clean -modcache

.PHONY: format-proto
format-proto:  ## Format protobuf sources
	$(API) buf format -w

.PHONY: lint-proto
lint-proto:  ## Lint protobuf sources
	$(API) buf lint

.PHONY: tidy-backend
tidy-backend: format-proto generate-backend  ## Format + tidy go sources
	$(API) $(GO) fmt ./...
	$(API) $(GO) mod tidy
	$(API) $(GO) mod verify
	$(API) goimports -w .
	$(API) golines -m 120 -w --ignore-generated .
	$(API) gci write --skip-generated -s standard -s "prefix($(MODULE))" -s default -s blank -s dot --custom-order .
	$(API) gofumpt -l -w .

.PHONY: lint-backend
lint-backend: lint-proto generate-backend  ## Run go linters (revive + golangci-lint)
	$(API) revive -config revive-config.toml -formatter friendly ./...
	$(API) golangci-lint run ./...

.PHONY: audit-backend
audit-backend: generate-backend  ## Run quality-control checks
	$(API) $(GO) mod verify
	$(API) $(GO) vet ./...
	$(API) staticcheck -checks=all,-ST1000,-U1000 ./...
	$(API) govulncheck ./...
	$(API) $(GO) test -race ./...

.PHONY: vet-backend
vet-backend: generate-backend  ## Run go vet
	$(API) $(GO) vet ./...

# ---------------------------------------------------------------------------
# Frontend (React/TypeScript) -- driven through pnpm scripts
# ---------------------------------------------------------------------------

.PHONY: install-frontend
install-frontend:  ## Install frontend dependencies
	$(FRONT) $(PNPM) install

.PHONY: build-frontend
build-frontend:  ## Build the frontend
	$(FRONT) $(PNPM) run build

.PHONY: dev-frontend
dev-frontend:  ## Start the frontend dev server
	$(FRONT) $(PNPM) run dev

.PHONY: preview-frontend
preview-frontend:  ## Preview the frontend build
	$(FRONT) $(PNPM) run preview

.PHONY: lint-frontend
lint-frontend:  ## Lint the frontend
	$(FRONT) $(PNPM) run lint

.PHONY: test-frontend
test-frontend:  ## Run frontend tests
	$(FRONT) $(PNPM) run test

# ---------------------------------------------------------------------------
# Project-wide (root-level) targets
# ---------------------------------------------------------------------------

.PHONY: all build
all: build  ## Build everything
build: build-backend build-frontend  ## Build backend + frontend

.PHONY: install install-backend
install: install-frontend install-backend  ## Install all dependencies
install-backend:  ## Download go module dependencies
	$(API) $(GO) mod download

.PHONY: test
test: test-backend test-frontend  ## Run all tests

.PHONY: lint
lint: lint-backend lint-frontend  ## Lint backend + frontend

.PHONY: tidy vet
tidy: tidy-backend  ## Format + tidy go sources
vet: vet-backend  ## Run go vet (alias for vet-backend)

.PHONY: audit
audit: audit-backend  ## Run all quality-control checks

.PHONY: clean
clean:  ## Remove build artifacts
	rm -rf $(API_DIR)/main $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/build

.PHONY: help
help:  ## Print this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make [target]\n\nTargets:\n"} \
	  /^[a-zA-Z0-9_-]+:.*?##/ { printf "  %-22s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
