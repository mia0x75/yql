.PHONY: all build test cover vet generate clean help

GO ?= go
GOYACC ?= goyacc

all: test ## Run tests (default)

help: ## Show this help
	@printf 'Usage: make <target>\n\n'
	@printf 'Targets:\n'
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all packages
	@$(GO) build ./...

test: ## Run all tests
	@$(GO) test ./...

cover: ## Run tests and print a coverage summary
	@$(GO) test ./... -coverprofile=/tmp/yql-cover.out
	@$(GO) tool cover -func=/tmp/yql-cover.out | tail -1

vet: ## Run go vet
	@$(GO) vet ./...

generate: ## Regenerate grammar/parser.go with goyacc
	@$(GOYACC) -o grammar/parser.go grammar/parser.y
	@sed 's/return __yyfmt__.Sprintf("%q", c)/return __yyfmt__.Sprintf("%q", rune(c))/' grammar/parser.go > grammar/parser.go.tmp
	@mv grammar/parser.go.tmp grammar/parser.go
	@rm -f y.output

clean: ## Remove generated scratch files
	@rm -f y.output grammar/parser.go.tmp /tmp/yql-cover.out
