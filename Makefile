# go-apperror Makefile.
#
# Common entry points for development and CI.
# Run `make` (or `make help`) to see available targets.

# .PHONY declares that these targets are NOT real files — they're just labels
# for command sequences. Without this, if a file named (say) "test" happened
# to exist in the project root, Make would treat the target as "already built"
# and skip running its recipe. Listing every target here keeps Make from being
# fooled by any same-named file or directory.
.PHONY: help setup lint test cover bench vuln fmt fmt-check tidy verify ci tools install-hooks clean

# Default target shows help so a bare `make` is never destructive.
.DEFAULT_GOAL := help

# -----------------------------------------------------------------------------
# Help
# -----------------------------------------------------------------------------

help: ## Show this help.
	@printf "Available targets:\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-13s\033[0m %s\n", $$1, $$2}'

# -----------------------------------------------------------------------------
# One-shot bootstrap
# -----------------------------------------------------------------------------

setup: tools install-hooks ## One-shot bootstrap after a fresh clone (tools + hooks + deps + sanity build).
	# Fetch module dependencies and verify their checksums.
	go mod download
	go mod verify
	# Sanity compile to make sure the source builds. For a library this
	# produces no binaries; it just exercises the compiler on every package.
	go build ./...
	@echo ""
	@echo "Setup complete. Next:"
	@echo "  make ci      # run the full CI flow locally"
	@echo "  make help    # see all available targets"

# -----------------------------------------------------------------------------
# Code quality
# -----------------------------------------------------------------------------

lint: ## Run golangci-lint over the whole module.
	golangci-lint run ./...

vuln: ## Scan for known vulnerabilities in dependencies and stdlib.
	govulncheck ./...

fmt: ## Format code with gofumpt (writes changes in place; prints modified files).
	gofumpt -l -w .

fmt-check: ## Check formatting without writing; exit non-zero on diff.
	@diff=$$(gofumpt -d .); \
	if [ -n "$$diff" ]; then \
		echo "$$diff"; \
		echo ""; \
		echo "Run 'make fmt' to fix."; \
		exit 1; \
	fi

tidy: ## Run go mod tidy and verify modules.
	go mod tidy
	go mod verify

# -----------------------------------------------------------------------------
# Testing
# -----------------------------------------------------------------------------

test: ## Run tests with race detector enabled.
	go test -race -count=1 ./...

cover: ## Run tests with coverage; produces coverage.out and coverage.html.
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

bench: ## Run all benchmarks.
	go test -run=^$$ -bench=. -benchmem ./...

# -----------------------------------------------------------------------------
# CI aggregates
# -----------------------------------------------------------------------------

verify: lint vuln test ## Code-quality checks excluding formatting (lint + vuln + test). Used by pre-push.

ci: fmt-check verify ## Everything CI should run (fmt-check + verify).

# -----------------------------------------------------------------------------
# Dev tools
# -----------------------------------------------------------------------------

install-hooks: ## Enable repo-managed git hooks (one-time per clone).
	# Tell this repo's git to look in .githooks/ for hooks. The script is
	# already version-controlled with the executable bit set, so this is
	# the only step needed.
	git config core.hooksPath .githooks
	@echo "git hooks enabled (core.hooksPath = .githooks)"

tools: ## Install required dev tools (golangci-lint v2, govulncheck, gofumpt) into $GOBIN.
	# golangci-lint: meta-linter that aggregates staticcheck, errcheck, govet,
	# gocritic, errorlint, ineffassign, unused, etc. Runs them in parallel and
	# reads .golangci.yml for which checks to enable. v2 module path includes
	# /v2/ per Go's semantic import versioning rule for major version >= 2.
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	# govulncheck: official Go vulnerability scanner. Unlike generic CVE scanners
	# that flag every transitive dependency with a known CVE, this one does call-
	# graph analysis and reports only the vulnerabilities reachable from your code.
	go install golang.org/x/vuln/cmd/govulncheck@latest
	# gofumpt: stricter superset of gofmt. Enforces a more opinionated, more
	# consistent format than gofmt (handles cases gofmt leaves alone). Wired
	# as a formatter in .golangci.yml so `make fmt` / `make fmt-check` use it.
	go install mvdan.cc/gofumpt@latest

# -----------------------------------------------------------------------------
# Cleanup
# -----------------------------------------------------------------------------

clean: ## Remove generated artifacts (coverage files, test binaries).
	rm -f coverage.out coverage.html
	find . -name '*.test' -delete
