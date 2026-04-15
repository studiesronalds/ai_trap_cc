#!/usr/bin/env bash
set -euo pipefail

BINARY="aisim"
PKG="./cmd/aisim"
DIST="dist"

TARGETS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

rm -rf "$DIST"
mkdir -p "$DIST"

for target in "${TARGETS[@]}"; do
    OS="${target%/*}"
    ARCH="${target#*/}"
    OUT="$DIST/${BINARY}-${OS}-${ARCH}"
    [ "$OS" = "windows" ] && OUT="${OUT}.exe"
    echo "Building $OS/$ARCH → $OUT"
    GOOS="$OS" GOARCH="$ARCH" go build -o "$OUT" "$PKG"
done

echo ""
echo "Done. Binaries in $DIST/:"
ls -lh "$DIST"
