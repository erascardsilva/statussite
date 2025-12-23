#!/bin/bash
# Erasmo Cardoso - Dev
# Script de desinstalação do Status Site Monitor

set -e

APP_NAME="statussite"
INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"

echo "Desinstalando $APP_NAME..."

# remove executavel
if [ -f "$INSTALL_DIR/$APP_NAME" ]; then
    rm "$INSTALL_DIR/$APP_NAME"
    echo "Executável removido"
else
    echo "Executável não encontrado"
fi

# remove atalho do menu
if [ -f "$DESKTOP_DIR/$APP_NAME.desktop" ]; then
    rm "$DESKTOP_DIR/$APP_NAME.desktop"
    echo "Atalho do menu removido"
else
    echo "Atalho não encontrado"
fi

# atualiza cache
if command -v update-desktop-database &> /dev/null; then
    update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
fi

echo "Desinstalação concluída!"
