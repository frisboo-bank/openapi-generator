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
backend/generate:  ## Backend: Generate enum sources (installs pinned goenums, runs go generate)
	$(API) go install github.com/zarldev/goenums@v0.4.3
	$(API) PATH="$$(go env GOPATH)/bin:$$PATH" go generate ./...

.PHONY: backend/run
backend/run: backend/generate  ## Backend: Run the API server
	$(API) $(GO) run cmd/main.go run

.PHONY: backend/build
backend/build: backend/generate  ## Backend: Build the API binary
	$(API) $(GO) build cmd/main.go

.PHONY: backend/test
backend/test: backend/generate  ## Backend: Run go tests (race + coverage)
	$(API) $(GO) test -race -cover ./...

.PHONY: backend/bench
backend/bench: backend/generate  ## Backend: Run go benchmarks
	$(API) $(GO) test -run=XXXXXX -benchtime=10s -bench=./ || exit 1

.PHONY: backend/update
backend/update:  ## Backend: Bump go module dependencies
	$(API) $(GO) get -u ./...

.PHONY: backend/deps-update
backend/deps-update:  ## Backend: Update + tidy go dependencies
	$(API) $(GO) get -u -t -v ./...
	$(API) $(GO) mod tidy

.PHONY: backend/deps-cleancache
backend/deps-cleancache:  ## Backend: Clear the go module cache
	$(API) $(GO) clean -modcache

.PHONY: backend/format-proto
backend/format-proto:  ## Backend: Format protobuf sources
	$(API) buf format -w

.PHONY: backend/lint-proto
backend/lint-proto:  ## Backend: Lint protobuf sources
	$(API) buf lint

.PHONY: backend/tidy
backend/tidy: backend/format-proto backend/generate  ## Backend: Format + tidy go sources
	$(API) $(GO) fmt ./...
	$(API) $(GO) mod tidy
	$(API) $(GO) mod verify
	$(API) goimports -w .
	$(API) golines -m 120 -w --ignore-generated .
	$(API) gci write --skip-generated -s standard -s "prefix($(MODULE))" -s default -s blank -s dot --custom-order .
	$(API) gofumpt -l -w .

.PHONY: backend/lint
backend/lint: backend/lint-proto backend/generate  ## Backend: Run go linters (revive + golangci-lint)
	$(API) revive -config revive-config.toml -formatter friendly ./...
	$(API) golangci-lint run ./...

.PHONY: backend/audit
backend/audit: backend/generate  ## Backend: Run quality-control checks
	$(API) $(GO) mod verify
	$(API) $(GO) vet ./...
	$(API) staticcheck -checks=all,-ST1000,-U1000 ./...
	$(API) govulncheck ./...
	$(API) $(GO) test -race ./...

.PHONY: backend/vet
backend/vet: backend/generate  ## Backend: Run go vet
	$(API) $(GO) vet ./...

.PHONY: backend/install
backend/install:  ## Backend: Download go module dependencies
	$(API) $(GO) mod download

# ---------------------------------------------------------------------------
# Frontend (React/TypeScript) -- driven through pnpm scripts
# ---------------------------------------------------------------------------

.PHONY: frontend/install
frontend/install:  ## Frontend: Install frontend dependencies
	$(FRONT) $(PNPM) install

.PHONY: frontend/build
frontend/build:  ## Frontend: Build the frontend
	$(FRONT) $(PNPM) run build

.PHONY: frontend/dev
frontend/dev:  ## Frontend: Start the frontend dev server
	$(FRONT) $(PNPM) run dev

.PHONY: frontend/preview
frontend/preview:  ## Frontend: Preview the frontend build
	$(FRONT) $(PNPM) run preview

.PHONY: frontend/lint
frontend/lint:  ## Frontend: Lint the frontend
	$(FRONT) $(PNPM) run lint

.PHONY: frontend/test
frontend/test:  ## Frontend: Run frontend tests
	$(FRONT) $(PNPM) run test

# ---------------------------------------------------------------------------
# Project-wide (root-level) targets
# ---------------------------------------------------------------------------

.PHONY: build
build: backend/build frontend/build  ## Project: Build backend + frontend

.PHONY: project/install
project/install: frontend/install backend/install  ## Project: Install all dependencies

.PHONY: project/test
project/test: backend/test frontend/test  ## Project: Run all tests

.PHONY: project/lint
project/lint: backend/lint frontend/lint  ## Project: Lint backend + frontend

.PHONY: project/tidy
project/tidy: backend/tidy  ## Project: Format + tidy go sources

.PHONY: project/vet
project/vet: backend/vet  ## Project: Run go vet (alias for backend/vet)

.PHONY: project/audit
project/audit: backend/audit  ## Project: Run all quality-control checks

.PHONY: clean
clean:  ## Project: Remove build artifacts
	rm -rf $(API_DIR)/main $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/build

# ---------------------------------------------------------------------------
# Help -- grouped by category
# ---------------------------------------------------------------------------

.PHONY: help
help:  ## General: Print this help menu
	@awk 'BEGIN {FS = ":.*## "; printf "\033[1;34mUsage:\033[0m\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z0-9_\/-]+:.*?## / { \
			desc = $$2; \
			if (match(desc, /^[A-Za-z0-9_ ]+:/)) { \
				cat = substr(desc, 1, RLENGTH - 1); \
				text = substr(desc, RLENGTH + 2); \
			} else { \
				cat = "General"; \
				text = desc; \
			} \
			if (cat != current_cat) { \
				current_cat = cat; \
				printf "\n\033[1;35m## %s\033[0m\n", toupper(cat); \
			} \
			printf "  \033[36m%-25s\033[0m %s\n", $$1, text; \
		} \
		END { printf "\n" }' $(MAKEFILE_LIST)
