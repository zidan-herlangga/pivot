#!/bin/sh
# pivot - Multi-platform runtime version switcher
# Install script for Linux/macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/zidan-herlangga/pivot/main/scripts/install.sh | sh

set -e

REPO="zidan-herlangga/pivot"
BIN_DIR="${HOME}/.pivot/bin"
INSTALL_DIR="${HOME}/.pivot"

# Detect platform
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: ${ARCH}"; exit 1 ;;
esac

case "${OS}" in
    linux|darwin) ;;
    *) echo "Unsupported OS: ${OS}"; exit 1 ;;
esac

# GitHub API to get latest release
echo "  Checking latest version..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name"' | cut -d'"' -f4 || echo "latest")
if [ "${LATEST}" = "latest" ] || [ -z "${LATEST}" ]; then
    echo "  Could not determine latest version, using 'latest'"
    LATEST="latest"
fi

# Download binary
BINARY="pivot-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${BINARY}.tar.gz"
echo "  Downloading ${BINARY}..."
mkdir -p "${BIN_DIR}"
curl -fsSL "${URL}" | tar xz -C "${BIN_DIR}" 2>/dev/null || {
    echo "  Pre-built binary not found. Building from source..."
    if command -v go >/dev/null 2>&1; then
        mkdir -p "${INSTALL_DIR}/src"
        cd "${INSTALL_DIR}/src"
        git clone --depth 1 "https://github.com/${REPO}.git" 2>/dev/null || true
        if [ -d "pivot" ]; then
            cd pivot
            go build -o "${BIN_DIR}/pivot" .
        else
            echo "  Could not clone repository."
            echo "  Install Go manually then run:"
            echo "    go install github.com/${REPO}@latest"
            exit 1
        fi
    else
        echo "  Could not download pre-built binary and Go is not installed."
        echo "  Install Go from https://go.dev or download manually from:"
        echo "    https://github.com/${REPO}/releases"
        exit 1
    fi
}

chmod +x "${BIN_DIR}/pivot"

# Add to PATH if not already there
case ":${PATH}:" in
    *:"${BIN_DIR}":*) ;;
    *)
        SHELL_CONFIG=""
        case "${SHELL}" in
            */zsh) SHELL_CONFIG="${HOME}/.zshrc" ;;
            */bash) SHELL_CONFIG="${HOME}/.bashrc" ;;
            */fish) SHELL_CONFIG="${HOME}/.config/fish/config.fish" ;;
        esac
        
        if [ -n "${SHELL_CONFIG}" ]; then
            echo "" >> "${SHELL_CONFIG}"
            echo "# pivot runtime switcher" >> "${SHELL_CONFIG}"
            echo "export PATH=\"${BIN_DIR}:\$PATH\"" >> "${SHELL_CONFIG}"
            echo "  Added ${BIN_DIR} to PATH in ${SHELL_CONFIG}"
        else
            echo "  Add ${BIN_DIR} to your PATH manually."
        fi
        ;;
esac

echo ""
echo "  pivot installed successfully!"
echo "  Run 'pivot' to start."
echo "  Restart your terminal or run: source ${SHELL_CONFIG}"
