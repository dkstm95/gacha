#!/usr/bin/env sh
set -eu

VERSION="${VERSION:-0.1.20}"
OUT_DIR="${OUT_DIR:-dist}"
TARGETS="${TARGETS:-darwin/arm64 darwin/amd64 linux/arm64 linux/amd64}"

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

for target in $TARGETS; do
  os="${target%/*}"
  arch="${target#*/}"
  name="gacha-$os-$arch"
  echo "Building $name"
  GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o "$OUT_DIR/$name/gacha" ./cmd/gacha
  tar -C "$OUT_DIR/$name" -czf "$OUT_DIR/$name.tar.gz" gacha
  rm -rf "$OUT_DIR/$name"
done

if command -v shasum >/dev/null 2>&1; then
  (cd "$OUT_DIR" && shasum -a 256 *.tar.gz > checksums.txt)
elif command -v sha256sum >/dev/null 2>&1; then
  (cd "$OUT_DIR" && sha256sum *.tar.gz > checksums.txt)
fi

echo "Release artifacts written to $OUT_DIR"
