#!/bin/bash

# Pre-commit hook for Go projects
# This script runs formatting, linting, and tests before allowing a commit

set -e

echo "Running pre-commit checks..."

# Check if we're in a Git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "Error: Not in a Git repository"
    exit 1
fi

# Get list of staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' | grep -v vendor/ || true)

if [ -z "$STAGED_GO_FILES" ]; then
    echo "No Go files staged for commit"
    exit 0
fi

echo "Staged Go files:"
echo "$STAGED_GO_FILES"

# Format code
echo
echo "1. Formatting code..."
go fmt ./...
gofmt -s -w $STAGED_GO_FILES

# Add formatted files back to staging
for file in $STAGED_GO_FILES; do
    git add "$file"
done

# Run go vet
echo
echo "2. Running go vet..."
if ! go vet ./...; then
    echo "Error: go vet failed"
    exit 1
fi

# Run golangci-lint if available
echo
echo "3. Running linter..."
if command -v golangci-lint >/dev/null 2>&1; then
    if ! golangci-lint run --timeout=5m; then
        echo "Error: golangci-lint failed"
        exit 1
    fi
else
    echo "golangci-lint not found, skipping linting"
    echo "Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi

# Run tests
echo
echo "4. Running tests..."
if ! go test ./cmd/...; then
    echo "Error: tests failed"
    exit 1
fi

# Run go mod tidy
echo
echo "5. Running go mod tidy..."
go mod tidy

# Check if go.mod or go.sum changed
if ! git diff --exit-code go.mod go.sum >/dev/null 2>&1; then
    echo "go.mod or go.sum changed, adding to commit"
    git add go.mod go.sum
fi

echo
echo "✅ All pre-commit checks passed!"
echo
