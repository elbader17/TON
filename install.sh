#!/bin/bash
set -e

echo "🚀 TON Framework - Instalando..."

REPO_URL="${1:-https://github.com/elbader17/TON}"
INSTALL_DIR="${2:-$(pwd)}"

cd "$INSTALL_DIR"

if ! command -v git &> /dev/null; then
    echo "❌ Error: Git no está instalado"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "❌ Error: Go no está instalado"
    echo "   Instálalo desde: https://go.dev/doc/install"
    exit 1
fi

if [ ! -d ".git" ]; then
    echo "📥 Clonando repositorio..."
    git clone --depth 1 "$REPO_URL" . || {
        echo "❌ Error: No se pudo clonar el repositorio"
        echo "   Asegúrate de que el repositorio existe y tienes acceso"
        exit 1
    }
elif [ -f "go.mod" ]; then
    echo "📦 Ejecutando go mod tidy..."
    go mod tidy || {
        echo "❌ Error: Falló 'go mod tidy'"
        exit 1
    }
else
    echo "📦 Inicializando módulo Go..."
    go mod init github.com/ton/framework && go mod tidy || {
        echo "❌ Error: Falló la inicialización del módulo Go"
        exit 1
    }
fi

echo ""
echo "✅ ¡Listo!"
echo ""
echo "📌 Para levantar el servidor:"
echo "   go run ./cmd/ton"
echo ""
echo "🌐 Servidor disponible en: http://localhost:8080"
echo "   Endpoint: POST /tools"