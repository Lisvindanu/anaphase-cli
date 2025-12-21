#!/bin/bash

set -e

echo "🚀 Installing Anaphase CLI..."

# Install the binary
echo "📦 Installing anaphase binary..."
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest

# Check if go/bin is in PATH
if ! echo "$PATH" | grep -q "$HOME/go/bin"; then
    echo ""
    echo "⚠️  $HOME/go/bin is not in your PATH"
    echo ""

    # Detect shell
    SHELL_NAME=$(basename "$SHELL")

    case "$SHELL_NAME" in
        bash)
            SHELL_CONFIG="$HOME/.bashrc"
            ;;
        zsh)
            SHELL_CONFIG="$HOME/.zshrc"
            ;;
        fish)
            SHELL_CONFIG="$HOME/.config/fish/config.fish"
            ;;
        *)
            SHELL_CONFIG="$HOME/.profile"
            ;;
    esac

    echo "Would you like to add it to your PATH? (y/n)"
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
        # Add to shell config
        if [ "$SHELL_NAME" = "fish" ]; then
            echo 'set -gx PATH $HOME/go/bin $PATH' >> "$SHELL_CONFIG"
        else
            echo 'export PATH="$HOME/go/bin:$PATH"' >> "$SHELL_CONFIG"
        fi

        echo "✅ Added to $SHELL_CONFIG"
        echo ""
        echo "Please run: source $SHELL_CONFIG"
        echo "Or open a new terminal window"
    else
        echo ""
        echo "You can manually add it by running:"
        echo "  echo 'export PATH=\"\$HOME/go/bin:\$PATH\"' >> $SHELL_CONFIG"
    fi
else
    echo "✅ PATH already configured"
fi

echo ""
echo "🎉 Installation complete!"
echo ""
echo "Verify installation:"
echo "  anaphase --version"
echo ""
echo "Get started:"
echo "  anaphase init my-project"
echo "  cd my-project"
echo "  anaphase gen domain --name user --prompt \"User with email and name\""
echo ""
echo "Documentation: https://anaphygon.my.id"
