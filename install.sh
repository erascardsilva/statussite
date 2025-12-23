#!/bin/bash
# Erasmo Cardoso - Dev
# Script de instalação do Status Site Monitor

set -e

APP_NAME="statussite"
BUILD_PATH="./build/bin/$APP_NAME"
INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
ICON_DIR="$HOME/.local/share/icons"

echo "Instalando $APP_NAME..."

# checa se o executavel existe
if [ ! -f "$BUILD_PATH" ]; then
    echo "Erro: Executável não encontrado em $BUILD_PATH"
    echo "Execute 'wails build' primeiro!"
    exit 1
fi

# cria diretorios necessarios
mkdir -p "$INSTALL_DIR"
mkdir -p "$DESKTOP_DIR"
mkdir -p "$ICON_DIR"

# copia executavel
echo "Copiando executável para $INSTALL_DIR..."
cp "$BUILD_PATH" "$INSTALL_DIR/$APP_NAME"
chmod +x "$INSTALL_DIR/$APP_NAME"

# cria atalho no menu
echo "Criando atalho no menu de aplicativos..."
cat > "$DESKTOP_DIR/$APP_NAME.desktop" <<EOF
[Desktop Entry]
Name=Status Site Monitor
Comment=Monitor de disponibilidade de sites
Exec=$INSTALL_DIR/$APP_NAME
Icon=network-workgroup
Terminal=false
Type=Application
Categories=Network;Utility;
Keywords=monitor;status;network;sites;
EOF

chmod +x "$DESKTOP_DIR/$APP_NAME.desktop"

# atualiza cache do menu
if command -v update-desktop-database &> /dev/null; then
    update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
fi

echo "Instalação concluída!"
echo ""
echo "Você pode executar de 3 formas:"
echo "  1. Digite 'statussite' no terminal"
echo "  2. Procure por 'Status Site Monitor' no menu"
echo "  3. Execute: $INSTALL_DIR/$APP_NAME"
echo ""
echo "Para desinstalar: ./uninstall.sh"
