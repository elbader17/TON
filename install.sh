#!/bin/bash
set -e

echo "🚀 TON Framework - Descargando..."

REPO_URL="${1:-https://github.com/elbader17/TON}"
PROJECT_NAME="ton-framework"
TEMP_DIR=$(mktemp -d)

cd "$TEMP_DIR"

echo "📥 Clonando repositorio..."
git clone --depth 1 "$REPO_URL" "$PROJECT_NAME" 2>/dev/null || {
    echo "❌ Error: No se pudo clonar el repositorio"
    echo "   Asegúrate de que el repositorio existe y tienes acceso"
    exit 1
}

cd "$PROJECT_NAME"

echo "📦 Instalando dependencias..."
go mod tidy 2>/dev/null || go mod init github.com/ton/framework && go mod tidy

echo ""
echo "✅ ¡Listo! El proyecto está en: $(pwd)"
echo ""
echo "📌 Para levantar el servidor:"
echo "   cd $(pwd)"
echo "   go run ./cmd/ton"
echo ""
echo "🌐 El servidor estará disponible en: http://localhost:8080"
echo "   Endpoint: POST /tools"