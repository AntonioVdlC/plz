#!/bin/bash

# Script to install Git hooks for the project

set -e

echo "Installing Git hooks..."

# Check if we're in a Git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "Error: Not in a Git repository"
    exit 1
fi

# Get the Git hooks directory
HOOKS_DIR=$(git rev-parse --git-dir)/hooks
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Git hooks directory: $HOOKS_DIR"
echo "Scripts directory: $SCRIPT_DIR"

# Install pre-commit hook
echo
echo "Installing pre-commit hook..."
if [ -f "$HOOKS_DIR/pre-commit" ]; then
    echo "Backing up existing pre-commit hook..."
    mv "$HOOKS_DIR/pre-commit" "$HOOKS_DIR/pre-commit.backup.$(date +%Y%m%d_%H%M%S)"
fi

# Create symlink to our pre-commit script
ln -sf "$SCRIPT_DIR/pre-commit.sh" "$HOOKS_DIR/pre-commit"
chmod +x "$HOOKS_DIR/pre-commit"

echo "✅ Pre-commit hook installed successfully!"
echo
echo "The hook will run automatically before each commit."
echo "To bypass the hook for a specific commit, use: git commit --no-verify"
echo
echo "To test the hook manually, run: ./scripts/pre-commit.sh"
echo
