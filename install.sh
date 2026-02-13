#!/bin/sh
# Install script for pray-cli
# Usage: curl -sSL https://raw.githubusercontent.com/AbdElrahmaN31/pray-cli/main/install.sh | sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    msys*|mingw*|cygwin*) OS="windows" ;;
    *)
        echo "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="armv7" ;;
    *)
        echo "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# Get latest version from GitHub API
echo "${YELLOW}Fetching latest version...${NC}"
VERSION=$(curl -s https://api.github.com/repos/AbdElrahmaN31/pray-cli/releases/latest | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$VERSION" ]; then
    echo "${RED}Failed to fetch latest version${NC}"
    exit 1
fi

echo "${GREEN}Latest version: $VERSION${NC}"

# Remove 'v' prefix from version for filename
VERSION_NUM=${VERSION#v}

# Set file extension
if [ "$OS" = "windows" ]; then
    EXT="zip"
    BINARY="pray.exe"
else
    EXT="tar.gz"
    BINARY="pray"
fi

# Construct download URL
FILENAME="pray-cli_${VERSION_NUM}_${OS}_${ARCH}.${EXT}"
URL="https://github.com/AbdElrahmaN31/pray-cli/releases/download/${VERSION}/${FILENAME}"

echo "${YELLOW}Downloading $FILENAME...${NC}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
if [ "$EXT" = "zip" ]; then
    curl -sSL "$URL" -o pray.zip
    unzip -q pray.zip
    rm pray.zip
else
    curl -sSL "$URL" | tar xz
fi

# Install binary
INSTALL_DIR="/usr/local/bin"
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY"
    echo "${GREEN}✓ Installed pray to $INSTALL_DIR${NC}"
else
    echo "${YELLOW}No write permission to $INSTALL_DIR${NC}"
    echo "${YELLOW}Installing to ./pray instead${NC}"
    mv "$BINARY" "$HOME/"
    chmod +x "$HOME/$BINARY"
    echo "${GREEN}✓ Installed pray to $HOME${NC}"
    echo "${YELLOW}Add $HOME to your PATH or move pray to /usr/local/bin manually:${NC}"
    echo "  sudo mv $HOME/pray /usr/local/bin/"
fi

# Cleanup
cd - > /dev/null
rm -rf "$TMP_DIR"

echo ""
echo "${GREEN}Installation complete!${NC}"
echo "Run 'pray --help' to get started"