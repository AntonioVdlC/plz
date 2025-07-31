# Development Setup Guide

This guide helps you set up the development environment for the `plz` CLI tool.

## Quick Setup

```bash
# 1. Install development tools and pre-commit hooks
make setup

# 2. Add Go bin directory to your PATH (if not already done)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc  # for zsh
# OR
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc # for bash

# 3. Reload your shell or source the config
source ~/.zshrc  # or ~/.bashrc

# 4. Verify everything works
make quality
```

## Manual Setup

If the automatic setup doesn't work, follow these steps:

### 1. Install Go Tools

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Verify installation
golangci-lint version
```

### 2. Install Pre-commit (Choose one method)

```bash
# Method 1: Using pip (Python)
pip install pre-commit
# OR
pip3 install pre-commit

# Method 2: Using Homebrew (macOS)
brew install pre-commit

# Method 3: Using conda
conda install -c conda-forge pre-commit
```

### 3. Setup Pre-commit Hooks (Choose one approach)

#### Option A: Pre-commit Framework (Recommended)
```bash
pre-commit install
pre-commit run --all-files  # Test on all files
```

#### Option B: Simple Git Hooks
```bash
./scripts/install-hooks.sh
./scripts/pre-commit.sh     # Test manually
```

### 4. Verify Setup

```bash
# Test all quality checks
make quality

# Test individual components
make fmt     # Format code
make vet     # Static analysis
make lint    # Comprehensive linting
make test    # Run tests
```

## Troubleshooting

### golangci-lint: command not found

This usually means `$GOPATH/bin` is not in your PATH. Fix it:

```bash
# Check where Go installs binaries
go env GOPATH

# Add to your shell profile (choose your shell)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc   # zsh
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc  # bash

# Reload shell configuration
source ~/.zshrc  # or ~/.bashrc

# Or temporarily for current session
export PATH=$PATH:$(go env GOPATH)/bin
```

### pre-commit: command not found

Install pre-commit using one of the methods above, or use the simple Git hooks:

```bash
./scripts/install-hooks.sh
```

### Make commands fail

If you don't want to install all tools, use basic quality checks:

```bash
make quality-basic  # Only fmt, vet, and test (no golangci-lint)
```

## Development Workflow

```bash
# Before starting work
git pull origin main

# Make your changes
# ...

# Run quality checks
make quality

# Commit (pre-commit hooks will run automatically)
git add .
git commit -m "Your commit message"

# Push
git push origin feature-branch
```

## Available Make Commands

```bash
make help           # Show all available commands
make build          # Build the binary
make test           # Run tests with coverage
make quality        # Run all quality checks
make quality-basic  # Run basic checks (no golangci-lint)
make setup          # Setup development environment
make clean          # Clean build artifacts
```

## IDE Integration

### VS Code

Install these extensions for the best experience:

- Go (by Google)
- golangci-lint (by golangci)
- Error Lens
- GitLens

### GoLand/IntelliJ

GoLand has built-in Go support. Enable golangci-lint integration:

1. Go to Settings → Tools → Go Linter
2. Enable golangci-lint
3. Set configuration file to `.golangci.yml`
