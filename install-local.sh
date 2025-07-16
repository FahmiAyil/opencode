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
PROJECT_DIR="/home/fahmi/Pandora/Projects/OPENCODE"

print_message info "Installing ${ORANGE}opencode${GREEN} locally to your system..."

# Install dependencies
print_message info "Installing dependencies..."
bun install

# Create install directory
mkdir -p "$INSTALL_DIR"

# Create wrapper script
print_message info "Creating opencode wrapper..."
cat > "$INSTALL_DIR/opencode" << EOF
#!/usr/bin/env bash
cd "$PROJECT_DIR"
exec bun run packages/opencode/src/index.ts "\$@"
EOF

chmod +x "$INSTALL_DIR/opencode"

add_to_path() {
    local config_file=$1
    local command=$2

    if grep -Fxq "$command" "$config_file"; then
        print_message info "Command already exists in $config_file, skipping write."
    elif [[ -w $config_file ]]; then
        echo -e "\n# opencode" >> "$config_file"
        echo "$command" >> "$config_file"
        print_message info "Successfully added ${ORANGE}opencode ${GREEN}to \$PATH in $config_file"
    else
        print_message warning "Manually add the directory to $config_file (or similar):"
        print_message info "  $command"
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

config_file=""
for file in $config_files; do
    if [[ -f $file ]]; then
        config_file=$file
        break
    fi
done

if [[ -z $config_file ]]; then
    print_message error "No config file found for $current_shell."
    exit 1
fi

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    case $current_shell in
        fish)
            add_to_path "$config_file" "fish_add_path $INSTALL_DIR"
        ;;
        *)
            add_to_path "$config_file" "export PATH=$INSTALL_DIR:\$PATH"
        ;;
    esac
fi

print_message info "✅ ${ORANGE}opencode${GREEN} installed! You can now run ${YELLOW}opencode${GREEN} from anywhere."
print_message info "To uninstall, run: ${YELLOW}./uninstall-local.sh${NC}"