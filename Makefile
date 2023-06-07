.DEFAULT_GOAL := all

# Pretty colours

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

# Import a local .env file and export variables for use
ifneq (,$(wildcard ./.env))
	include .env
    export
endif

# Variables

BUF			?= buf
PROTOS		:= $(shell find api -name "*.proto")
CMDS		:= $(wildcard ./cmd/*)
BINARIES	:= $(subst cmd,bin,$(CMDS))
GO_FILES	:= $(shell find . -name "*.go")
LOCAL		?= TRUE

.PHONY: all
all: lint api test build

# Build
build: $(BINARIES) ## Build project

bin/%: $(GO_FILES)
	go build -o $@ ./cmd/$*

# Run
run-%:
	go run ./cmd/$*

# Test
.PHONY: test
test:
	go test -race ./... -coverprofile cover.out

# Lint
.PHONY: lint
lint: ## Run linters over assets
	$(BUF) lint -v
	goimports -d -e -w ./cmd ./pkg
	golangci-lint run -v ./...

# Formatting
.PHONY: fmt
fmt:
	go fmt ./...

# Vendor
.PHONY: vendor
vendor: ## Tidy go.mod, update dependencies, and vendor them
	go mod tidy
	go mod vendor

# API generation
api: ./pkg/api/.gen ## Generate API client/server code

pkg/api/.gen: .env $(PROTOS)
	$(BUF) generate -v
	@touch pkg/api/.gen

# Manage dependencies
.PHONY: deps
deps: ## Install dependencies for this project
	go install github.com/bufbuild/buf/cmd/buf@$(BUF_VER)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VER)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VER)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VER)
	go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VER)

# Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
