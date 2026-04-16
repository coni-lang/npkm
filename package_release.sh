#!/bin/bash
set -e

# ======================================================
# NPKM-Coni Build & Package Script
# Cross-compiles npkm-coni for macOS and Windows
# then packages a Windows release zip.
#
# Usage: ./package_release.sh
# ======================================================

# Define which Coni source tree to use
CONI_SRC="/Users/nico/cool/s5/coni-lang-gitea"
export CONI_HOME="$CONI_SRC"

# Ensure typical paths for Go are available
export PATH="$PATH:/usr/local/go/bin:/opt/homebrew/bin"

BUILD_DATE=$(date '+%Y-%m-%d-%H%M')
DIST_DIR="dist"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "============================================"
echo "  NPKM-Coni Build & Package"
echo "  Date: $BUILD_DATE"
echo "  Using Coni Source: $CONI_SRC"
echo "============================================"

# Build the fresh compiler binary
TEMP_CONI_BIN="/tmp/coni-compiler"
echo ""
echo "▸ Building latest Coni compiler from source..."
cd "$CONI_SRC"
go build -o "$TEMP_CONI_BIN" .
echo "  ✓ Compiler built at $TEMP_CONI_BIN"

# 0. Run tests
echo ""
echo "▸ Running tests..."
cd "$SCRIPT_DIR/npkm-coni"
"$TEMP_CONI_BIN" test ...

# 1. Clean dist
cd "$SCRIPT_DIR"
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

# 2. Build macOS (native arm64)
echo ""
echo "▸ Building macOS binary (darwin/arm64)..."
cd "$SCRIPT_DIR/npkm-coni"
"$TEMP_CONI_BIN" build . -o "$SCRIPT_DIR/$DIST_DIR/npkm-coni"

# 3. Build Windows (cross-compile amd64)
echo ""
echo "▸ Building Windows binary (windows/amd64)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 "$TEMP_CONI_BIN" build . -o "$SCRIPT_DIR/$DIST_DIR/npkm-coni.exe"

cd "$SCRIPT_DIR"

# 4. Copy binaries back into npkm-coni/
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

# 6. Deploy to samba share
SAMBA_DIR="/Volumes/share/npkm"
if [ -d "$SAMBA_DIR" ]; then
  echo ""
  echo "▸ Deploying to samba share..."
  pv "$DIST_DIR/$ARCHIVE_NAME" > "$SAMBA_DIR/$ARCHIVE_NAME"
  echo "  ✓ Copied to $SAMBA_DIR/$ARCHIVE_NAME"
else
  echo ""
  echo "⚠ Samba share not mounted at $SAMBA_DIR — skipping deploy"
  echo "  Mount it and run:"
  echo "  pv $DIST_DIR/$ARCHIVE_NAME > $SAMBA_DIR/$ARCHIVE_NAME"
fi
