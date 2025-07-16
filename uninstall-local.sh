#!/usr/bin/env bash
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
ORANGE='\033[38;2;255;140;0m'
NC='\033[0m' # No Color

print_message() {
    local level=$1
    local message=$2
    local color=""

    case $level in
        info) color="${GREEN}" ;;
        warning) color="${YELLOW}" ;;
        error) color="${RED}" ;;
    esac

    echo -e "${color}${message}${NC}"
}

INSTALL_DIR=$HOME/.opencode/bin

print_message info "Uninstalling local ${ORANGE}opencode${GREEN} from your system..."

# Remove wrapper script
print_message info "Removing opencode wrapper..."
rm -f "$INSTALL_DIR/opencode"

# Remove install directory if empty
if [[ -d "$INSTALL_DIR" ]] && [[ -z "$(ls -A "$INSTALL_DIR")" ]]; then
    rmdir "$INSTALL_DIR"
    rmdir "$HOME/.opencode" 2>/dev/null || true
fi

print_message info "✅ Local ${ORANGE}opencode${GREEN} uninstalled successfully!"
print_message warning "Note: PATH entries in your shell config remain. Remove manually if needed."
print_message info "Your local project files remain unchanged."