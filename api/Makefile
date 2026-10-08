MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

SHELL := $(shell command -v bash 2>/dev/null || command -v sh)
ifneq ($(findstring bash,$(SHELL)),)
.SHELLFLAGS := -euo pipefail -c
else
.SHELLFLAGS := -c
endif

GO   ?= go
CGO  ?= 0
PNPM ?= pnpm

MODULE   := $(shell $(GO) list -m 2>/dev/null || echo "unknown")
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

GOOGLEAPIS_SHA   ?= d71f2ee15b7ec0e8cc7ce19f3dce68a5a341ace2
GRPC_GATEWAY_SHA ?= 069005e84c857849d0419b10abb9c0e50910631b

NODE_BIN := $(CURDIR)/node_modules/.bin

.DEFAULT_GOAL := help

.PHONY: generate
GOENUMS := $(shell $(GO) env GOPATH)/bin/goenums
$(GOENUMS):
	$(GO) install github.com/zarldev/goenums@v0.4.3
generate: $(GOENUMS)  ## Generate: Generate enum sources (installs pinned goenums, runs go generate)
	PATH="$(dir $(GOENUMS)):$$PATH" $(GO) generate ./pkg/...

.PHONY: run
run: generate  ## Run: Run the application (development)
	CGO_ENABLED=$(CGO) $(GO) run $(BUILDFLAGS) $(CMD_DIR) run --environment=development

.PHONY: build
build: generate  ## Build: Build the application binary
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=$(CGO) $(GO) build $(BUILDFLAGS) -o $(BIN_DIR)/$(NAME) $(CMD_DIR)

.PHONY: build-debug
build-debug: generate  ## Build: Build the application binary with debug enabled
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 $(GO) build -gcflags "all=-N -l" -o $(BIN_DIR)/$(NAME)-debug $(CMD_DIR)

.PHONY: test
test: generate  ## Test: Run tests (race + coverage)
	CGO_ENABLED=1 $(GO) test -race -count=1 -coverprofile=coverage.out -covermode=atomic ./...
	CGO_ENABLED=1 $(GO) tool cover -func=coverage.out | tail -n 1

.PHONY: test-ci
test-ci: generate  ## Test: Run tests (CI-optimized)
  $(GO) test -race -count=1 -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: bench
bench: generate  ## Bench: Run go benchmarks (10s each)
	$(GO) test -run=^$$ -bench=. -benchtime=10s -count=1 ./...

.PHONY: mocks
mocks:  ## Generate: Generate mocks with mockgen
	$(GO) generate ./mocks/gen.go

.PHONY: update
update:  ## Deps: Bump go module dependencies
	$(GO) get -u ./...

.PHONY: deps-update
deps-update:  ## Deps: Update + tidy go dependencies
	$(GO) get -u -t -v ./...
	$(GO) mod tidy

.PHONY: deps-cleancache
deps-cleancache:  ## Deps: Clear the go module cache
	$(GO) clean -modcache

.PHONY: tidy
tidy: proto-format generate  ## Tidy: Format Go code, tidy modules, organize imports
	$(GO) fmt ./...
	$(GO) mod tidy
	$(GO) mod verify
	goimports -w ./...
	golines -m 120 -w --ignore-generated ./...
	gci write --skip-generated -s standard -s "prefix($(MODULE))" -s default -s blank -s dot --custom-order ./...
	gofumpt -l -w .

.PHONY: lint
lint: proto-lint generate  ## Lint: Run go linters (revive + golangci-lint)
	revive -config revive-config.toml -formatter friendly ./...
	golangci-lint run ./...

.PHONY: vet
vet: generate  ## Vet: Run go vet
	$(GO) vet ./...

.PHONY: audit
audit: vet  ## Audit: Security audit and static analysis
	$(GO) mod verify
	staticcheck -checks=all,-ST1000,-U1000 ./...
	govulncheck ./...

.PHONY: check
check: lint test audit  ## Check: lint + test + audit

.PHONY: install
install:  ## Install: Download go module dependencies
	$(GO) mod download

.PHONY: install-tools
install-tools:  ## Install: Install development tools
	$(GO) install github.com/daixiang0/gci@v0.6.10
	$(GO) install go.uber.org/mock/mockgen@v0.5.0
	$(GO) install github.com/mgechev/revive@v1.7.0
	$(GO) install github.com/segmentio/golines@v0.12.2
	$(GO) install golang.org/x/tools/cmd/goimports@v0.31.0
	$(GO) install golang.org/x/vuln/cmd/govulncheck@v1.1.3
	$(GO) install honnef.co/go/tools/cmd/staticcheck@v0.6.1
	$(GO) install mvdan.cc/gofumpt@v0.7.0
	$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.26.3
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.26.3
	$(GO) install github.com/wasilibs/go-protoc-gen-connect-query/cmd/protoc-gen-connect-query@v0.3.0
	$(GO) install github.com/bufbuild/protoc-gen-es@v1.0.0

# ---------------------------------------------------------------------------
# Api proto (gRPC/buf) -- generated code lives under api/, frontend/grpc
# ---------------------------------------------------------------------------

GRPC_PROTO_MODULE := frisboo-bank/openapi-generator-service
GRPC_PROTO_DIR := ./internal/shared/grpc/proto
GRPC_PROTO_SHARED_DIR := ./pkg/rpc/rpc_server/resources/grpc/proto
GRPC_PROTO_VENDOR_DIR := ./internal/shared/grpc/vendors
GRPC_SWAGGER_OUT := ./docs/swagger
GRPC_GO_OUT := .
GRPC_ES_OUT ?= ../frontend/grpc/client

.PHONY: proto-shared
proto-shared:  ## Proto: Generate shared protobuf code
	protoc \
		--proto_path=$(GRPC_PROTO_VENDOR_DIR) \
		--proto_path=$(GRPC_PROTO_SHARED_DIR) \
		--go_out=$(GRPC_GO_OUT) \
		--go_opt=module=$(GRPC_PROTO_MODULE) \
		$(shell find $(GRPC_PROTO_SHARED_DIR) -name '*.proto')

.PHONY: proto-build
proto-build: proto-shared  ## Proto: Generate service protobuf/gRPC code
	mkdir -p $(GRPC_ES_OUT) $(GRPC_SWAGGER_OUT)
	protoc \
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
		--plugin=protoc-gen-es=$(NODE_BIN)/protoc-gen-es \
		--plugin=protoc-gen-connect-query=$(NODE_BIN)/protoc-gen-connect-query \
		--connect-query_out=$(GRPC_ES_OUT) \
		--es_out=$(GRPC_ES_OUT) \
		$(shell find $(GRPC_PROTO_DIR) -name '*.proto')

.PHONY: proto-format
proto-format:  ## Proto: Format protobuf files
	buf format -w

.PHONY: proto-lint
proto-lint:  ## Proto: Lint protobuf schemas
	buf lint

.PHONY: proto-vendor
GOOGLEAPIS_SHA := b1b9c2a5c5e9f8d7a6b4c3d2e1f0a9b8c7d6e5f4
GRPC_GATEWAY_SHA := v2.26.3

proto-vendor:  ## Proto: Download third-party proto dependencies locally
	mkdir -p $(GRPC_PROTO_VENDOR_DIR)/google/api
	mkdir -p $(GRPC_PROTO_VENDOR_DIR)/protoc-gen-openapiv2/options
	echo "Downloading Google API protos..."
	curl -fsSL -o $(GRPC_PROTO_VENDOR_DIR)/google/api/annotations.proto https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_SHA)/google/api/annotations.proto
	curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/google/api/annotations.proto https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_SHA)/google/api/annotations.proto
	curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/google/api/http.proto https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_SHA)/google/api/http.proto
	curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/google/api/field_behavior.proto https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_SHA)/google/api/field_behavior.proto
	echo "Downloading grpc-gateway OpenAPIv2 protos..."
	curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/protoc-gen-openapiv2/options/annotations.proto https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/$(GRPC_GATEWAY_SHA)/protoc-gen-openapiv2/options/annotations.proto
	curl -sSL -o $(GRPC_PROTO_VENDOR_DIR)/protoc-gen-openapiv2/options/openapiv2.proto https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/$(GRPC_GATEWAY_SHA)/protoc-gen-openapiv2/options/openapiv2.proto
	echo "Done vendoring protos."

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
