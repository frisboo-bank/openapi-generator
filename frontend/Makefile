MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

SHELL := $(shell if command -v bash >/dev/null 2>&1; then echo bash; else echo sh; fi)
ifeq ($(SHELL),bash)
.SHELLFLAGS := -euo pipefail -c
else
.SHELLFLAGS := -c
endif

PNPM ?= pnpm

.DEFAULT_GOAL := help

.PHONY: install
install: ## Install: Install dependencies
	$(PNPM) install

.PHONY: dev
dev: ## Run: Start the development server
	$(PNPM) dev

.PHONY: build
build: ## Build: Build the production bundle
	$(PNPM) build

.PHONY: preview
preview: ## Build: Preview the production build
	$(PNPM) preview

.PHONY: lint
lint: ## Lint: Lint sources
	$(PNPM) lint

.PHONY: test
test: ## Test: Run tests
	$(PNPM) test

.PHONY: test-watch
test-watch: ## Test: Run tests in watch mode
	$(PNPM) test:watch

.PHONY: check
check: lint test  ## Check: lint + test

.PHONY: help
help: ## General: Print this help menu
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
