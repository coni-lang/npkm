#!/bin/bash
set -e

# ======================================================
# NPKM-Coni Build & Package Wrapper
# Delegates to the native EDN playbook.
# ======================================================

echo "▸ Bootstrapping release process via NPKM-Coni playbook..."
cd "$(dirname "$0")"

if [ ! -f "npkm-coni/npkm-coni" ]; then
    echo "⚠ Local npkm-coni binary not found! Please build it first."
    exit 1
fi

./npkm-coni/npkm-coni package_release.edn

