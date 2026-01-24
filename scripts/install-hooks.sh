#!/bin/bash
# Install git hooks for Pulsar

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

echo "Installing git hooks..."

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

# Install pre-commit hook
cp "$SCRIPT_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
chmod +x "$HOOKS_DIR/pre-commit"

echo "Installed: pre-commit hook (gitleaks secret scanning)"

# Check if gitleaks is installed
if ! command -v gitleaks &> /dev/null && [ ! -x "$HOME/.local/bin/gitleaks" ]; then
    echo ""
    echo "Warning: gitleaks is not installed."
    echo "Install with:"
    echo "  mkdir -p ~/.local/bin"
    echo "  curl -sSfL https://github.com/gitleaks/gitleaks/releases/download/v8.18.4/gitleaks_8.18.4_linux_x64.tar.gz | tar -xz -C ~/.local/bin"
    echo ""
fi

echo "Done! Git hooks installed successfully."
