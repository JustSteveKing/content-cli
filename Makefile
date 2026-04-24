BINARY  := content
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/juststeveking/content-cli/cmd.Version=$(VERSION)"

.PHONY: build install clean tidy fmt vet test check snapshot help

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY) .

install: ## Install to $$GOPATH/bin
	go install $(LDFLAGS) .

clean: ## Remove built binary
	rm -f $(BINARY)

tidy: ## Tidy go.mod and go.sum
	go mod tidy

fmt: ## Format all Go source files
	go fmt ./...

vet: ## Run go vet across all packages
	go vet ./...

test: ## Run tests
	go test ./...

check: fmt vet ## Run fmt and vet (use in CI)

snapshot: ## Build a local snapshot with GoReleaser (no publish)
	goreleaser release --snapshot --clean
