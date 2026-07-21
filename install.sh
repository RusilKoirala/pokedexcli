#!/bin/bash

set -e

# colors for output :D
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' 

# detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)


case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# gitHub repository info
REPO="rusilkoirala/pokedexcli"
BINARY_NAME="pokedexcli"

echo -e "${GREEN}Installing PokedexCLI...${NC}"
echo "OS: $OS"
echo "Architecture: $ARCH"


echo -e "${YELLOW}Fetching latest release...${NC}"
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Failed to fetch latest version${NC}"
    exit 1
fi

echo "Latest version: $LATEST_VERSION"


DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/${BINARY_NAME}-${OS}-${ARCH}"


if [ "$OS" = "windows" ]; then
    DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
fi

# download binary
echo -e "${YELLOW}Downloading $BINARY_NAME...${NC}"
TMP_FILE="/tmp/$BINARY_NAME"

if ! curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"; then
    echo -e "${RED}Failed to download binary${NC}"
    echo "URL attempted: $DOWNLOAD_URL"
    exit 1
fi


chmod +x "$TMP_FILE"


INSTALL_DIR="/usr/local/bin"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    echo -e "${GREEN}Installed to $INSTALL_DIR/$BINARY_NAME${NC}"
else

    if command -v sudo &> /dev/null; then
        echo -e "${YELLOW}Installing to $INSTALL_DIR (requires sudo)...${NC}"
        sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
        echo -e "${GREEN}Installed to $INSTALL_DIR/$BINARY_NAME${NC}"
    else

        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
        mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
        echo -e "${GREEN}Installed to $INSTALL_DIR/$BINARY_NAME${NC}"
        echo -e "${YELLOW}Note: Make sure $INSTALL_DIR is in your PATH${NC}"
        echo "Add this to your ~/.bashrc or ~/.zshrc:"
        echo "  export PATH=\"\$HOME/bin:\$PATH\""
    fi
fi

# verify installation
if command -v $BINARY_NAME &> /dev/null; then
    echo -e "${GREEN}✓ Installation successful!${NC}"
    echo ""
    echo "Run it with: $BINARY_NAME"
else
    echo -e "${YELLOW}Installation complete, but $BINARY_NAME is not in PATH${NC}"
    echo "Try running: $INSTALL_DIR/$BINARY_NAME"
fi
