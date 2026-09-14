MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

SHELL := $(shell if command -v bash >/dev/null 2>&1; then echo bash; else echo sh; fi)
ifeq ($(SHELL),bash)
.SHELLFLAGS := -euo pipefail -c
else
.SHELLFLAGS := -c
endif

# Sub-project directories.
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
.DEFAULT_GOAL := help

# ---------------------------------------------------------------------------
# Backend (Go) -- every go/buf/lint tool runs against the api/ module
# ---------------------------------------------------------------------------

.PHONY: backend/generate
backend/generate:  ## Generate enum sources (installs pinned goenums, runs go generate)
	$(API) go install github.com/zarldev/goenums@v0.4.3
	$(API) PATH="$$(go env GOPATH)/bin:$$PATH" go generate ./...

.PHONY: backend/run
backend/run: backend/generate  ## Run the API server
	$(API) $(GO) run cmd/main.go run

.PHONY: backend/build
backend/build: backend/generate  ## Build the API binary
	$(API) $(GO) build cmd/main.go

.PHONY: backend/test
backend/test: backend/generate  ## Run go tests (race + coverage)
	$(API) $(GO) test -race -cover ./...

.PHONY: backend/bench
backend/bench: backend/generate  ## Run go benchmarks
	$(API) $(GO) test -run=XXXXXX -benchtime=10s -bench=./ || exit 1

.PHONY: backend/update
backend/update:  ## Bump go module dependencies
	$(API) $(GO) get -u ./...

.PHONY: backend/deps-update
backend/deps-update:  ## Update + tidy go dependencies
	$(API) $(GO) get -u -t -v ./...
	$(API) $(GO) mod tidy

.PHONY: backend/deps-cleancache
backend/deps-cleancache:  ## Clear the go module cache
	$(API) $(GO) clean -modcache

.PHONY: backend/format-proto
backend/format-proto:  ## Format protobuf sources
	$(API) buf format -w

.PHONY: backend/lint-proto
backend/lint-proto:  ## Lint protobuf sources
	$(API) buf lint

.PHONY: backend/tidy
backend/tidy: backend/format-proto backend/generate  ## Format + tidy go sources
	$(API) $(GO) fmt ./...
	$(API) $(GO) mod tidy
	$(API) $(GO) mod verify
	$(API) goimports -w .
	$(API) golines -m 120 -w --ignore-generated .
	$(API) gci write --skip-generated -s standard -s "prefix($(MODULE))" -s default -s blank -s dot --custom-order .
	$(API) gofumpt -l -w .

.PHONY: backend/lint
backend/lint: backend/lint-proto backend/generate  ## Run go linters (revive + golangci-lint)
	$(API) revive -config revive-config.toml -formatter friendly ./...
	$(API) golangci-lint run ./...

.PHONY: backend/audit
backend/audit: backend/generate  ## Run quality-control checks
	$(API) $(GO) mod verify
	$(API) $(GO) vet ./...
	$(API) staticcheck -checks=all,-ST1000,-U1000 ./...
	$(API) govulncheck ./...
	$(API) $(GO) test -race ./...

.PHONY: backend/vet
backend/vet: backend/generate  ## Run go vet
	$(API) $(GO) vet ./...

.PHONY: backend/install
backend/install:  ## Download go module dependencies
	$(API) $(GO) mod download

# ---------------------------------------------------------------------------
# Frontend (React/TypeScript) -- driven through pnpm scripts
# ---------------------------------------------------------------------------

.PHONY: frontend/install
frontend/install:  ## Install frontend dependencies
	$(FRONT) $(PNPM) install

.PHONY: frontend/build
frontend/build:  ## Build the frontend
	$(FRONT) $(PNPM) run build

.PHONY: frontend/dev
frontend/dev:  ## Start the frontend dev server
	$(FRONT) $(PNPM) run dev

.PHONY: frontend/preview
frontend/preview:  ## Preview the frontend build
	$(FRONT) $(PNPM) run preview

.PHONY: frontend/lint
frontend/lint:  ## Lint the frontend
	$(FRONT) $(PNPM) run lint

.PHONY: frontend/test
frontend/test:  ## Run frontend tests
	$(FRONT) $(PNPM) run test

# ---------------------------------------------------------------------------
# Project-wide (root-level) targets
# ---------------------------------------------------------------------------

.PHONY: all build
all: build  ## Build everything
build: backend/build frontend/build  ## Build backend + frontend

.PHONY: project/install project/test project/lint project/tidy project/vet project/audit
project/install: frontend/install backend/install  ## Install all dependencies

.PHONY: project/test
project/test: backend/test frontend/test  ## Run all tests

.PHONY: project/lint
project/lint: backend/lint frontend/lint  ## Lint backend + frontend

.PHONY: project/tidy
project/tidy: backend/tidy  ## Format + tidy go sources

.PHONY: project/vet
project/vet: backend/vet  ## Run go vet (alias for backend/vet)

.PHONY: project/audit
project/audit: backend/audit  ## Run all quality-control checks

.PHONY: clean
clean:  ## Remove build artifacts
	rm -rf $(API_DIR)/main $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/build

# ---------------------------------------------------------------------------
# Help -- grouped by category
# ---------------------------------------------------------------------------

.PHONY: help
help:  ## Print this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make [category]/target\n\n"} \
		CAT == "" { CAT = "project"; } \
		/^[a-zA-Z0-9_/-]+:.*?##/ { \
			cat = "project"; \
			if ($$1 ~ /^backend\//) cat = "backend"; \
			if ($$1 ~ /^frontend\//) cat = "frontend"; \
			if (cat != CAT) { \
				CAT = cat; \
				printf "\n## %s\n", toupper(cat); \
			} \
			target = $$1; sub(/:.*$$/, "", target); \
			print "  " target; \
		} \
		END { printf "\n" }' $(MAKEFILE_LIST)
