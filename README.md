# plz

A collection of useful CLI utilities for everyday development tasks.

## Installation

### Build from source

```bash
git clone <repository-url>
cd plz
go build -o plz .
```

### Usage

```bash
./plz [command] [flags] [arguments]
```

For help with any command:
```bash
./plz --help
./plz [command] --help
```

## Commands

### `hash` - Generate hash values

Generate MD5, SHA1, or SHA256 hash values for strings or files.

```bash
# Hash a string (default: SHA256)
./plz hash "hello world"

# Use different hash types
./plz hash --type md5 "hello world"
./plz hash --type sha1 "hello world"

# Hash file contents
./plz hash --file myfile.txt
./plz hash --file --type md5 myfile.txt
```

### `encode` - Encode/decode strings

Encode or decode strings using base64 or URL encoding.

```bash
# Base64 encode (default)
./plz encode "hello world"

# Base64 decode
./plz encode --decode "aGVsbG8gd29ybGQ="

# URL encode
./plz encode --type url "hello world!"

# URL decode
./plz encode --type url --decode "hello%20world%21"
```

### `random` - Generate random values

Generate random strings, numbers, or UUIDs.

```bash
# Random string (default: 16 characters)
./plz random

# Custom length string
./plz random --type string --length 32

# Random number in range
./plz random --type number --min 1 --max 100

# Generate UUID
./plz random --type uuid
```

### `time` - Convert and format timestamps

Convert between Unix timestamps and human-readable dates.

```bash
# Current time
./plz time

# Convert Unix timestamp to date
./plz time 1640995200

# Convert date to timestamp
./plz time --timestamp "2022-01-01 00:00:00"

# Use different timezone
./plz time --timezone "America/New_York" 1640995200

# Custom format
./plz time --format "2006-01-02 15:04:05" 1640995200
```

### `json` - Process JSON data

Pretty print, minify, or validate JSON data.

```bash
# Pretty print JSON
./plz json '{"name":"test","items":[1,2,3]}'

# Minify JSON
./plz json --minify '{"name": "test", "items": [1, 2, 3]}'

# Validate JSON
./plz json --validate '{"valid": true}'

# Process JSON file
./plz json --file data.json
./plz json --file --minify data.json
```

## Development

### Prerequisites

- Go 1.24.4 or later

### Building

```bash
# Build binary
go build -o plz .

# Run without building
go run . [command] [args]

# Format code
go fmt ./...

# Run static analysis
go vet ./...

# Run tests
go test ./cmd/...

# Run tests with coverage
go test -race -coverprofile=coverage.out -covermode=atomic ./cmd/...

# Development automation (Makefile)
make build          # Build binary
make quality        # Run all quality checks (fmt, vet, lint, test)
make setup          # Setup development environment with pre-commit hooks
```

### Development Setup

Quick setup for development:

```bash
# Install tools and setup pre-commit hooks
make setup

# Add Go bin to PATH (if needed)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# Verify setup
make quality
```

**Having issues?** See [SETUP.md](SETUP.md) for detailed setup instructions and troubleshooting.

### Adding New Commands

1. Create a new file `cmd/cmd_newcommand.go`
2. Follow the existing command pattern using Cobra
3. Register the command in the `init()` function
4. Add tests in `cmd/cmd_newcommand_test.go`
5. See `CLAUDE.md` for detailed patterns and conventions

## CI/CD

[![CI](https://github.com/AntonioVdlc/plz/workflows/CI/badge.svg)](https://github.com/AntonioVdlc/plz/actions)

This project uses GitHub Actions for continuous integration:

- **Tests**: Runs across Go versions 1.21, 1.22, 1.23
- **Code Quality**: golangci-lint, go vet, go fmt checks
- **Coverage**: 93.5% test coverage
- **Cross-platform builds**: Linux, macOS, Windows (AMD64, ARM64)

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [pflag](https://github.com/spf13/pflag) - Flag parsing
