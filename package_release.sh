#!/bin/bash
set -e

# ======================================================
# NPKM-Coni Build & Package Script
# Uses `coni build` with GOOS/GOARCH for cross-compilation
# ======================================================

BUILD_DATE=$(date '+%Y-%m-%d-%H%M')
DIST_DIR="dist"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "============================================"
echo "  NPKM-Coni Build & Package"
echo "  Date: $BUILD_DATE"
echo "============================================"

# 0. Run tests first
echo ""
echo "▸ Running tests..."
cd "$SCRIPT_DIR/npkm-coni"
coni test ...
cd "$SCRIPT_DIR"

# 1. Clean dist
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

# 2. Build macOS (native arm64)
echo ""
echo "▸ Building macOS binary (darwin/arm64)..."
cd "$SCRIPT_DIR/npkm-coni"
coni build . -o "$SCRIPT_DIR/$DIST_DIR/npkm-coni"
echo "  ✓ npkm-coni (macOS arm64)"

# 3. Build Windows (cross-compile amd64)
echo ""
echo "▸ Building Windows binary (windows/amd64)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 coni build . -o "$SCRIPT_DIR/$DIST_DIR/npkm-coni.exe"
echo "  ✓ npkm-coni.exe (Windows amd64)"

cd "$SCRIPT_DIR"

# 4. Copy binaries back into npkm-coni/ for convenience
echo ""
echo "▸ Updating local binaries..."
cp "$DIST_DIR/npkm-coni" "npkm-coni/npkm-coni"
cp "$DIST_DIR/npkm-coni.exe" "npkm-coni/npkm-coni.exe"

# 5. Package Windows release zip
ARCHIVE_NAME="npkm-coni-windows-amd64-${BUILD_DATE}.zip"
echo ""
echo "▸ Packaging Windows release: $ARCHIVE_NAME"

cd "$DIST_DIR"
cp "$SCRIPT_DIR/README.md" .
cp "$SCRIPT_DIR/npkm-coni/test-playbook.edn" .
cp "$SCRIPT_DIR/test-playbook.yml" .

zip -r "$ARCHIVE_NAME" npkm-coni.exe README.md test-playbook.edn test-playbook.yml
cd "$SCRIPT_DIR"

echo ""
echo "============================================"
echo "  ✅ Build & Package Complete"
echo "============================================"
echo ""
echo "Artifacts:"
ls -lh "$DIST_DIR/npkm-coni"
ls -lh "$DIST_DIR/npkm-coni.exe"
ls -lh "$DIST_DIR/$ARCHIVE_NAME"
