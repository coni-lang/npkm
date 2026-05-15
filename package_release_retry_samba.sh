#!/bin/bash
set -e

echo "▸ Retrying deploy to samba share..."
cd "$(dirname "$0")/dist"

LATEST_ZIP=$(ls -t npkm-coni-release-*.zip 2>/dev/null | head -n 1)

if [ -z "$LATEST_ZIP" ]; then
    echo "⚠ No release zip found in dist/! Run package_release.sh first."
    exit 1
fi

echo "Found release artifact: $LATEST_ZIP"

if [ -d "/Volumes/share/npkm" ]; then
    echo "Copying to samba share..."
    pv "$LATEST_ZIP" > "/Volumes/share/npkm/$LATEST_ZIP"
    echo "Done."
else
    echo "Samba share not mounted at /Volumes/share/npkm — skipping deploy"
fi
