# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**plz** is a Go CLI tool providing useful utilities for development tasks. Built with Cobra CLI framework using a single-package, modular architecture.

## Build and Development Commands

```bash
# Build binary
go build -o plz .

# Run without building
go run . [command] [args]

# Development
go mod tidy
go fmt ./...
go vet ./...

# Testing
go test ./cmd/...       # Run all tests
go test ./cmd/ -v       # Run tests with verbose output

# Pre-commit hooks (choose one approach)
make setup              # Install tools and setup pre-commit framework
./scripts/install-hooks.sh  # Install simple Git hooks (alternative)

# Development workflow
make quality            # Run all quality checks (fmt, vet, lint, test)
```

## Architecture

### File Organization Pattern
- `main.go` - Application entry point
- `cmd/` - Command package directory
  - `cmd/root.go` - Root command definition and Execute function
  - `cmd/cmd_*.go` - Individual command implementations (one file per command)
  - `cmd/*_test.go` - Test files for each command
  - `cmd/testdata/` - Test data files
- `.github/workflows/` - GitHub Actions CI/CD workflows
- `scripts/` - Development scripts and Git hooks
- `.golangci.yml` - Linter configuration
- `.pre-commit-config.yaml` - Pre-commit framework configuration
- `Makefile` - Development automation
- Package structure: `main` package imports `plz/cmd` package

### Command Structure Template
Each command follows this pattern:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var commandCmd = &cobra.Command{
    Use:   "command [args]",
    Short: "Brief description",
    Long:  "Detailed description",
    Args:  cobra.ExactArgs(1), // or other validators
    RunE:  runCommand,
}

var (
    commandFlag string
    commandBool bool
)

func init() {
    commandCmd.Flags().StringVarP(&commandFlag, "flag", "f", "default", "Description")
    commandCmd.Flags().BoolVarP(&commandBool, "bool", "b", false, "Description")
    rootCmd.AddCommand(commandCmd)
}

func runCommand(cmd *cobra.Command, args []string) error {
    // Implementation
    return nil
}
```

### Key Conventions
- **Error Handling**: Return `error` from run functions, use `fmt.Errorf()` with wrapping
- **Flag Naming**: Use both short (`-f`) and long (`--flag`) versions
- **File Input Pattern**: Many commands support `--file` flag for file vs string input
- **Output Format**: Include operation context in output (e.g., "Encoded (BASE64): ...")

## Adding New Commands

1. Create `cmd/cmd_newcommand.go` in the cmd package
2. Follow the command structure template above (use `package cmd`)
3. Register command in `init()` function with `rootCmd.AddCommand(newcommandCmd)`
4. Implement the `runNewcommand` function with proper error handling

## Current Commands

- `encode` - Base64/URL encoding and decoding
- `hash` - MD5/SHA1/SHA256 hash generation for strings or files
- `json` - JSON pretty printing, minification, and validation
- `random` - Random string/number/UUID generation
- `time` - Timestamp conversion and formatting

## Dependencies

- **Cobra CLI**: `github.com/spf13/cobra` - Primary CLI framework
- **Go Version**: 1.24.4

## CI/CD

The project uses GitHub Actions for continuous integration:

### Workflows
- **ci.yml** - Runs on push/PR to main branch
  - Tests across Go versions 1.21, 1.22, 1.23
  - Runs `go vet`, `go fmt` check, and all tests
  - Generates test coverage reports
  - Runs golangci-lint for code quality
  - Tests binary functionality

- **release.yml** - Runs on git tags/releases
  - Builds cross-platform binaries (Linux, macOS, Windows)
  - Supports AMD64 and ARM64 architectures
  - Uploads release artifacts

### Coverage
Current test coverage: **93.5%** of statements

### Code Quality
- Uses golangci-lint with custom configuration
- Enforces formatting with `gofmt`
- Static analysis with `go vet`

## Development Tools

### Pre-commit Hooks
Two approaches available for pre-commit hooks:

#### 1. Pre-commit Framework (Recommended)
```bash
# Setup (installs tools and hooks)
make setup

# Or manually:
pip install pre-commit
pre-commit install

# Run on all files
pre-commit run --all-files
```

#### 2. Simple Git Hooks
```bash
# Install custom Git hooks
./scripts/install-hooks.sh

# Test manually
./scripts/pre-commit.sh
```

### Makefile Commands
```bash
make build          # Build binary
make test           # Run tests with coverage
make quality        # Run all quality checks
make clean          # Clean build artifacts
make help           # Show all commands
```
