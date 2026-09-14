MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

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

GO    ?= go
PNPM  ?= pnpm
CGO   ?= 0

MODULE   := $(shell go -C $(API_DIR) list -m 2>/dev/null || echo "unknown")
NAME     := $(notdir $(MODULE))
CMD_DIR  := ./cmd
BIN_DIR  := ./bin

VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT     ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE       ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_LDFLAGS ?= -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(DATE) \
	-buildid= \
	-extldflags "-static"

GOFLAGS    ?= -trimpath
BUILDFLAGS ?= $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -tags "osusergo,netgo"

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
backend/run: backend/generate  ## Backend: Run the application (development)
	$(API) CGO_ENABLED=$(CGO) $(GO) run $(BUILDFLAGS) $(CMD_DIR) run --environment=development

.PHONY: backend/build
backend/build: backend/generate  ## Backend: Build the application binary
	$(API) mkdir -p $(BIN_DIR)
	$(API) CGO_ENABLED=$(CGO) $(GO) build $(BUILDFLAGS) -o $(BIN_DIR)/$(NAME) $(CMD_DIR)

.PHONY: backend/build-debug
backend/build-debug: backend/generate  ## Backend: Build the application binary with debug enabled
	$(API) mkdir -p $(BIN_DIR)
	$(API) CGO_ENABLED=1 $(GO) build -gcflags "all=-N -l" -o $(BIN_DIR)/$(NAME)-debug $(CMD_DIR)

.PHONY: backend/test
backend/test: backend/generate  ## Backend: Run tests (race + coverage)
	$(API) $(GO) test -race -coverprofile=coverage.out -covermode=atomic ./...
	$(API) $(GO) tool cover -func=coverage.out | tail -n 1

.PHONY: backend/test-ci
backend/test-ci: backend/generate  ## Backend: Run tests (CI-optimized, no race detector)
	$(API) $(GO) test -race -count=1 -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: backend/bench
backend/bench: backend/generate  ## Backend: Run go benchmarks (10s each)
	$(API) $(GO) test -run=^$$ -bench=. -benchtime=10s -count=1 ./...

.PHONY: backend/mocks
backend/mocks:  ## Backend: Generate mocks with mockgen
	$(API) $(GO) generate ./mocks/gen.go

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

.PHONY: backend/tidy
backend/tidy: backend/format-proto backend/generate  ## Backend: Format Go code, tidy modules, organize imports
	$(API) $(GO) fmt ./...
	$(API) $(GO) mod tidy
	$(API) $(GO) mod verify
	$(API) goimports -w ./...
	$(API) golines -m 120 -w --ignore-generated ./...
	$(API) gci write --skip-generated -s standard -s "prefix($(MODULE))" -s default -s blank -s dot --custom-order ./...
	$(API) gofumpt -l -w .

.PHONY: backend/lint
backend/lint: backend/lint-proto backend/generate  ## Backend: Run go linters (revive + golangci-lint)
	$(API) revive -config revive-config.toml -formatter friendly ./...
	$(API) golangci-lint run ./...

.PHONY: backend/vet
backend/vet: backend/generate  ## Backend: Run go vet
	$(API) $(GO) vet ./...

.PHONY: backend/audit
backend/audit: backend/vet  ## Backend: Security audit and static analysis
	$(API) $(GO) mod verify
	$(API) staticcheck -checks=all,-ST1000,-U1000 ./...
	$(API) govulncheck ./...

.PHONY: backend/check
backend/check: backend/lint backend/test backend/audit  ## Backend: lint + test + audit

.PHONY: backend/install
backend/install:  ## Backend: Download go module dependencies
	$(API) $(GO) mod download

.PHONY: backend/install-tools
backend/install-tools:  ## Backend: Install development tools
	$(API) $(GO) install github.com/daixiang0/gci@latest
	$(API) $(GO) install github.com/golang/mock/mockgen@latest
	$(API) $(GO) install github.com/mgechev/revive@latest
	$(API) $(GO) install github.com/segmentio/golines@latest
	$(API) $(GO) install golang.org/x/tools/cmd/goimports@latest
	$(API) $(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(API) $(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	$(API) $(GO) install mvdan.cc/gofumpt@latest
	$(API) $(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	$(API) $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(API) $(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(API) $(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	$(API) $(GO) install github.com/wasilibs/go-protoc-gen-connect-query/cmd/protoc-gen-connect-query@latest

# ---------------------------------------------------------------------------
# Backend proto (gRPC/buf) -- generated code lives under api/, frontend/grpc
# ---------------------------------------------------------------------------

GRPC_PROTO_MODULE := frisboo-bank/openapi-generator-service
GRPC_PROTO_DIR := ./internal/shared/grpc/proto
GRPC_PROTO_SHARED_DIR := ./pkg/rpc/rpc_server/resources/grpc/proto
GRPC_PROTO_VENDOR_DIR := ./internal/shared/grpc/vendors
GRPC_SWAGGER_OUT := ./docs/swagger
GRPC_GO_OUT := .
GRPC_ES_OUT := ../frontend/grpc/client

.PHONY: backend/proto-shared
backend/proto-shared:  ## Backend: Generate shared protobuf code
	$(API) protoc \
		--proto_path=$(GRPC_PROTO_VENDOR_DIR) \
		--proto_path=$(GRPC_PROTO_SHARED_DIR) \
		--go_out=$(GRPC_GO_OUT) \
		--go_opt=module=$(GRPC_PROTO_MODULE) \
		$(shell find $(GRPC_PROTO_SHARED_DIR) -name '*.proto')

.PHONY: backend/proto-build
backend/proto-build: backend/proto-shared  ## Backend: Generate service protobuf/gRPC code
	$(API) mkdir -p $(GRPC_ES_OUT) $(GRPC_SWAGGER_OUT)
	$(API) protoc \
		--proto_path=$(GRPC_PROTO_VENDOR_DIR) \
		--proto_path=$(GRPC_PROTO_SHARED_DIR) \
		--proto_path=$(GRPC_PROTO_DIR) \
		--go_out=$(GRPC_GO_OUT) \
		--go_opt=module=$(GRPC_PROTO_MODULE) \
		--go-grpc_out=$(GRPC_GO_OUT) \
		--go-grpc_opt=module=$(GRPC_PROTO_MODULE) \
		--grpc-gateway_out=$(GRPC_GO_OUT) \
		--grpc-gateway_opt=module=$(GRPC_PROTO_MODULE),generate_unbound_methods=true \
		--openapiv2_out=$(GRPC_SWAGGER_OUT) \
		--openapiv2_opt=allow_merge=true,merge_file_name=entity_service \
		--connect-query_out=$(GRPC_ES_OUT) \
		--es_out=$(GRPC_ES_OUT) \
		$(shell find $(GRPC_PROTO_DIR) -name '*.proto')

.PHONY: backend/format-proto
backend/format-proto:  ## Backend: Format protobuf files
	$(API) buf format -w

.PHONY: backend/lint-proto
backend/lint-proto:  ## Backend: Lint protobuf schemas
	$(API) buf lint

.PHONY: backend/proto-vendor
backend/proto-vendor:  ## Backend: Download third-party proto dependencies locally
	$(API) mkdir -p $(GRPC_PROTO_VENDOR_DIR)/google/api
	$(API) mkdir -p $(GRPC_PROTO_VENDOR_DIR)/protoc-gen-openapiv2/options
	$(API) echo "Downloading Google API protos..."
	$(API) curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/google/api/annotations.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
	$(API) curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/google/api/http.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto
	$(API) curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/google/api/field_behavior.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto
	$(API) echo "Downloading grpc-gateway OpenAPIv2 protos..."
	$(API) curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/protoc-gen-openapiv2/options/annotations.proto https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/annotations.proto
	$(API) curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/protoc-gen-openapiv2/options/openapiv2.proto https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/openapiv2.proto
	$(API) echo "Done vendoring protos."

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

.PHONY: project/check
project/check: backend/check frontend/lint frontend/test  ## Project: lint + test + audit (both)

.PHONY: clean
clean:  ## Project: Remove build artifacts and coverage reports
	rm -rf $(API_DIR)/$(BIN_DIR) $(API_DIR)/coverage.out $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/build

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
