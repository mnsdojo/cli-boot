#!/bin/bash

# Check if Go is installed
if ! command -v go &>/dev/null; then
  echo "Go is not installed. Please install Go and try again."
  exit 1
fi

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR" || exit

# Clone the repository
git clone https://github.com/mnsdojo/cli-boot.git
cd cli-boot || exit

# Initialize Go module if go.mod doesn't exist
if [ ! -f go.mod ]; then
  go mod init github.com/mnsdojo/cli-boot
  go mod tidy
fi

# Build the program
go build -o cli-boot

# Set the correct permissions
chmod +x cli-boot

# Move the binary to /usr/local/bin (or other specified directory)
INSTALL_DIR=${1:-/usr/local/bin}
sudo mv cli-boot "$INSTALL_DIR/"

# Ensure correct permissions after moving
sudo chmod 755 "$INSTALL_DIR/cli-boot"

# Clean up temporary directory
cd ..
rm -rf "$TEMP_DIR"

# Add cli-boot to both .bashrc and .zshrc for universal support
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  echo 'export PATH=$PATH:'"$INSTALL_DIR" >>~/.bashrc
  echo 'export PATH=$PATH:'"$INSTALL_DIR" >>~/.zshrc
  echo "cli-boot has been added to both your .bashrc and .zshrc files."
  echo "Please run 'source ~/.bashrc' or 'source ~/.zshrc' or start a new terminal session to use cli-boot."
else
  echo "cli-boot is already in your PATH."
fi

echo "cli-boot has been installed successfully. You can now run it by typing 'cli-boot' in the terminal."
