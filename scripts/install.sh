#!/bin/bash

# Colors for output
# Orange text
ORANGE='\033[38;5;208m'
NC='\033[0m'

# Display banner
echo -e "${ORANGE}
             __                               __                           __               
            /  |                             /  |                         /  |              
   _______  \$\$ |   ______    __    __    ____\$\$ |   ______     ______    _\$\$ |_      ______  
  /       | \$\$ |  /      \\  /  |  /  |  /    \$\$ |  /      \\   /      \\  / \$\$   |    /      \\ 
 /\$\$\$\$\$\$\$/  \$\$ | /\$\$\$\$\$\$  | \$\$ |  \$\$ | /\$\$\$\$\$\$\$ | /\$\$\$\$\$\$  |  \$\$\$\$\$\$  | \$\$\$\$\$/    /\$\$\$\$\$\$  |
 \$\$ |       \$\$ | \$\$ |  \$\$ | \$\$ |  \$\$ | \$\$ |  \$\$ | \$\$ |  \$\$ |  /    \$\$ |   \$\$ | __  \$\$    \$\$ |
 \$\$ \\_____  \$\$ | \$\$ \\__\$\$ | \$\$ \\__\$\$ | \$\$ \\__\$\$ | \$\$ \\__\$\$ | /\$\$\$\$\$\$\$ |   \$\$ |/  | \$\$\$\$\$\$\$\$/ 
 \$\$       | \$\$ | \$\$    \$\$/  \$\$    \$\$/  \$\$    \$\$ | \$\$    \$\$ | \$\$    \$\$ |   \$\$  \$\$/  \$\$       |
  \$\$\$\$\$\$\$/  \$\$/   \$\$\$\$\$\$/    \$\$\$\$\$\$/    \$\$\$\$\$\$\$/   \$\$\$\$\$\$\$ |  \$\$\$\$\$\$/     \$\$\$\$/    \$\$\$\$\$\$\$/ 
                                                  /  \\__\$\$ |                              
                                                  \$\$    \$\$/                               
                                                   \$\$\$\$\$\$/                                
${NC}"

echo "Installing cloudgate..."

# Determine installation directory
if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
else
    INSTALL_DIR="$HOME/.local/bin"
    # Create ~/.local/bin if it doesn't exist
    mkdir -p "$INSTALL_DIR"
fi

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert architecture names
case $ARCH in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
esac

# Set binary name based on detected OS/arch
BINARY_NAME="cloudgate_${OS}_${ARCH}"
LATEST_RELEASE="https://github.com/HenryOwenz/cloudgate/releases/latest/download/${BINARY_NAME}"

# Download the binary
if command -v curl > /dev/null; then
    echo "Downloading with curl..."
    curl -L "$LATEST_RELEASE" -o "$INSTALL_DIR/cg"
elif command -v wget > /dev/null; then
    echo "Downloading with wget..."
    wget -O "$INSTALL_DIR/cg" "$LATEST_RELEASE"
else
    echo "Error: Neither curl nor wget found. Please install either one and try again."
    exit 1
fi

# Make binary executable
chmod +x "$INSTALL_DIR/cg"

# Add to PATH if needed
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.bashrc"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.zshrc" 2>/dev/null || true
fi

echo "Successfully installed cloudgate!"
echo "Please restart your terminal or run 'source ~/.bashrc' to use the 'cg' command." 
