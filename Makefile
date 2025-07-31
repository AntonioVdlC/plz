.PHONY: build test lint fmt vet clean install-tools pre-commit help

# Build configuration
BINARY_NAME=plz
BUILD_DIR=.
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	go test -race -coverprofile=coverage.out -covermode=atomic ./cmd/...
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -race ./cmd/...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		echo "Please ensure $$GOPATH/bin is in your PATH"; \
		echo "You can add this to your shell profile:"; \
		echo "  export PATH=\$$PATH:\$$(go env GOPATH)/bin"; \
		echo "Or run: export PATH=\$$PATH:\$$(go env GOPATH)/bin"; \
		echo "Attempting to run golangci-lint with full path..."; \
		$$(go env GOPATH)/bin/golangci-lint run --timeout=5m; \
	fi

# Run all quality checks
quality: fmt vet lint test
	@echo "All quality checks completed!"

# Run basic quality checks (without golangci-lint)
quality-basic: fmt vet test
	@echo "Basic quality checks completed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	go clean

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		echo "golangci-lint installed successfully"; \
	else \
		echo "golangci-lint already installed"; \
	fi
	@echo "Installing pre-commit (requires Python)..."
	@if ! command -v pre-commit >/dev/null 2>&1; then \
		if command -v pip3 >/dev/null 2>&1; then \
			pip3 install pre-commit; \
		elif command -v pip >/dev/null 2>&1; then \
			pip install pre-commit; \
		elif command -v brew >/dev/null 2>&1; then \
			echo "Using Homebrew to install pre-commit..."; \
			brew install pre-commit; \
		else \
			echo "Unable to install pre-commit automatically."; \
			echo "Please install pre-commit manually:"; \
			echo "  https://pre-commit.com/#installation"; \
			echo "Or use the simple Git hooks: ./scripts/install-hooks.sh"; \
		fi; \
	else \
		echo "pre-commit already installed"; \
	fi
	@echo "Checking Go installation..."
	@go version

# Install pre-commit hooks
pre-commit-install:
	@echo "Installing pre-commit hooks..."
	pre-commit install

# Run pre-commit on all files
pre-commit-all:
	@echo "Running pre-commit on all files..."
	pre-commit run --all-files

# Update pre-commit hooks
pre-commit-update:
	@echo "Updating pre-commit hooks..."
	pre-commit autoupdate

# Run pre-commit checks manually (without installing hooks)
pre-commit: fmt vet lint test
	@echo "Manual pre-commit checks completed!"

# Install and setup everything for development
setup: install-tools pre-commit-install
	@echo "Development environment setup completed!"
	@echo "You can now run 'make quality' to run all checks"

# Cross-platform build
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-linux-arm64 $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe $(BUILD_DIR)
	GOOS=windows GOARCH=arm64 go build -o $(BINARY_NAME)-windows-arm64.exe $(BUILD_DIR)
	@echo "Cross-platform builds completed!"

# Help target
help:
	@echo "Available targets:"
	@echo "  build           Build the binary"
	@echo "  test            Run tests with coverage"
	@echo "  test-verbose    Run tests with verbose output"
	@echo "  fmt             Format code"
	@echo "  vet             Run go vet"
	@echo "  lint            Run golangci-lint"
	@echo "  quality         Run all quality checks (fmt, vet, lint, test)"
	@echo "  quality-basic   Run basic quality checks (fmt, vet, test)"
	@echo "  pre-commit      Run manual pre-commit checks"
	@echo "  clean           Clean build artifacts"
	@echo "  install-tools   Install development tools"
	@echo "  setup           Setup development environment"
	@echo "  pre-commit-install   Install pre-commit hooks"
	@echo "  pre-commit-all       Run pre-commit on all files"
	@echo "  pre-commit-update    Update pre-commit hooks"
	@echo "  build-all       Build for multiple platforms"
	@echo "  help            Show this help message"
