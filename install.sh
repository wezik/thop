#!/bin/sh

set -eu

# Check for curl early since it's required
if ! command -v curl >/dev/null 2>&1; then
    echo "!!! curl is not installed. Please install it and try again."
    exit 1
fi

# Check if thop is installed
if ! command -v thop >/dev/null 2>&1; then
    echo "Thop is not installed. Installing..."
else
    echo "Thop is already installed. Updating..."
fi

download_url="https://github.com/wezik/thop/releases/latest/download/thop"
binary="thop"
tmp_file=$(mktemp /tmp/thop.XXXXXX)
install_dir="/usr/local/bin"
version=$(curl -s https://api.github.com/repos/wezik/thop/releases/latest | \
  grep '"tag_name":' | \
  sed -E 's/.*"([^"]+)".*/\1/'
)

echo "Downloading latest version($version) of Thop..."
curl -Ls "$download_url" -o "$tmp_file"

echo "Installing Thop to $install_dir (requires sudo)..."
sudo mv "$tmp_file" "$install_dir/$binary"

echo "Updating permissions..."
sudo chmod +x "$install_dir/$binary"

echo "Thop successfully installed. Run 'thop' to start. Or 'thop --help' for more info."

# Check runtime dependencies
for dep in tmux fzf; do
    if ! command -v "$dep" >/dev/null 2>&1; then
        echo "!!! $dep is not installed. It's required for thop to work. Please install it."
        exit 1
    fi
done
