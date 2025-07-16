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

# Remove PATH entries from shell config files
remove_from_path() {
    local config_file=$1
    local temp_file=$(mktemp)
    
    if [[ -f "$config_file" ]]; then
        # Remove the opencode section (comment + export line)
        sed '/^# opencode$/,/^export PATH.*\.opencode\/bin.*$/d' "$config_file" > "$temp_file"
        mv "$temp_file" "$config_file"
        print_message info "Removed opencode PATH from $config_file"
    fi
}

XDG_CONFIG_HOME=${XDG_CONFIG_HOME:-$HOME/.config}

current_shell=$(basename "$SHELL")
case $current_shell in
    fish)
        config_files="$HOME/.config/fish/config.fish"
    ;;
    zsh)
        config_files="$HOME/.zshrc $HOME/.zshenv $XDG_CONFIG_HOME/zsh/.zshrc $XDG_CONFIG_HOME/zsh/.zshenv"
    ;;
    bash)
        config_files="$HOME/.bashrc $HOME/.bash_profile $HOME/.profile $XDG_CONFIG_HOME/bash/.bashrc $XDG_CONFIG_HOME/bash/.bash_profile"
    ;;
    *)
        config_files="$HOME/.bashrc $HOME/.bash_profile"
    ;;
esac

for file in $config_files; do
    if [[ -f $file ]]; then
        remove_from_path "$file"
    fi
done

print_message info "✅ Local ${ORANGE}opencode${GREEN} uninstalled successfully!"
print_message info "PATH entries removed from shell config files."
print_message info "Your local project files remain unchanged."