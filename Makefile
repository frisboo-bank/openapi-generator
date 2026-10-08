MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

SHELL := $(shell if command -v bash >/dev/null 2>&1; then echo bash; else echo sh; fi)
ifeq ($(SHELL),bash)
.SHELLFLAGS := -euo pipefail -c
else
.SHELLFLAGS := -c
endif

API_DIR      := api
FRONTEND_DIR := frontend

.DEFAULT_GOAL := help

.PHONY: build
build:  ## Build api binary and frontend bundle
	$(MAKE) -C $(API_DIR) build
	$(MAKE) -C $(FRONTEND_DIR) build

.PHONY: install
install:  ## Install dependencies for both sub-projects
	$(MAKE) -C $(API_DIR) install
	$(MAKE) -C $(FRONTEND_DIR) install

.PHONY: install-tools
install-tools:  ## Install all development tools
	$(MAKE) -C $(API_DIR) install-tools

.PHONY: generate
generate:  ## Generate code (enums, protobuf, mocks)
	$(MAKE) -C $(API_DIR) generate
	$(MAKE) -C $(API_DIR) mocks

.PHONY: tidy
tidy:  ## Format + tidy sources in both sub-projects
	$(MAKE) -C $(API_DIR) tidy

.PHONY: vet
vet:  ## Run go vet (alias of api/vet)
	$(MAKE) -C $(API_DIR) vet

.PHONY: audit
audit:  ## Run all quality-control checks
	$(MAKE) -C $(API_DIR) audit

.PHONY: test
test:  ## Run tests in both sub-projects
	$(MAKE) -C $(API_DIR) test
	$(MAKE) -C $(FRONTEND_DIR) test

.PHONY: lint
lint:  ## Lint both sub-projects
	$(MAKE) -C $(API_DIR) lint
	$(MAKE) -C $(FRONTEND_DIR) lint

.PHONY: check
check: lint test audit  ## lint + test + audit (both)

.PHONY: clean
clean:  ## Remove build artifacts and coverage reports
	rm -rf $(API_DIR)/bin $(API_DIR)/coverage.out $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/build

.PHONY: run
run:  ## Run the api development server
	$(MAKE) -C $(API_DIR) run

.PHONY: bench
bench:  ## Run go benchmarks
	$(MAKE) -C $(API_DIR) bench

.PHONY: deps-update
deps-update:  ## Update + tidy go dependencies
	$(MAKE) -C $(API_DIR) deps-update

.PHONY: deps-cleancache
deps-cleancache:  ## Clear the go module cache
	$(MAKE) -C $(API_DIR) deps-cleancache

.PHONY: proto-build
proto-build:  ## Generate service protobuf/gRPC code
	$(MAKE) -C $(API_DIR) proto-build

.PHONY: proto-shared
proto-shared:  ## Generate shared protobuf code
	$(MAKE) -C $(API_DIR) proto-shared

.PHONY: format-proto
format-proto:  ## Format protobuf files
	$(MAKE) -C $(API_DIR) format-proto

.PHONY: lint-proto
lint-proto:  ## Lint protobuf schemas
	$(MAKE) -C $(API_DIR) lint-proto

.PHONY: proto-vendor
proto-vendor:  ## Download third-party proto dependencies locally
	$(MAKE) -C $(API_DIR) proto-vendor

.PHONY: help
help:  ## Print this help menu
	@awk 'BEGIN {FS = ":.*## "; printf "\033[1;34mUsage:\033[0m\n  make \033[36m<target>\033[0m\n  (run from project root)\n\n"} \
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
				printf "\033[1;35m## %s\033[0m\n", toupper(cat); \
			} \
			printf "  \033[36m%-25s\033[0m %s\n", $$1, text; \
		} \
		END { printf "\n" }' $(MAKEFILE_LIST)
