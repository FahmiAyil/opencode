#!/bin/bash
set -e

echo "Uninstalling local opencode from your system..."

# Remove global symlink
echo "Removing global symlink..."
npm unlink -g opencode

echo "✅ Local opencode uninstalled successfully!"
echo "Note: This only removes the global symlink. Your local project files remain unchanged."