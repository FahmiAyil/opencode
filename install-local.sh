#!/bin/bash
set -e

echo "Installing opencode locally to your system..."

# Install dependencies
echo "Installing dependencies..."
bun install

# Build the project (if needed)
echo "Building project..."
cd packages/opencode

# Create global symlink
echo "Creating global symlink..."
npm link

echo "✅ opencode installed! You can now run 'opencode' from anywhere."
echo "To uninstall, run: npm unlink -g opencode"