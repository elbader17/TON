#!/bin/bash
set -e

echo "🚀 TON Framework - Instalando..."

REPO_URL="${1:-https://github.com/elbader17/TON}"
INSTALL_DIR="${2:-$(pwd)}"

cd "$INSTALL_DIR"

if [ ! -d ".git" ]; then
    echo "📥 Clonando repositorio..."
    git clone --depth 1 "$REPO_URL" . 2>/dev/null || {
        echo "❌ Error: No se pudo clonar el repositorio"
        echo "   Asegúrate de que el repositorio existe y tienes acceso"
        exit 1
    }
fi

echo "📦 Instalando dependencias..."
go mod tidy 2>/dev/null || go mod init github.com/ton/framework && go mod tidy

echo ""
echo "✅ ¡Listo!"
echo ""
echo "📌 Para levantar el servidor:"
echo "   go run ./cmd/ton"
echo ""
echo "🌐 Servidor disponible en: http://localhost:8080"
echo "   Endpoint: POST /tools"