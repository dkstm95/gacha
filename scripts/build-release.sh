#!/usr/bin/env sh
set -eu

VERSION="${VERSION:-0.1.32}"
OUT_DIR="${OUT_DIR:-dist}"
TARGETS="${TARGETS:-darwin/arm64 darwin/amd64 linux/arm64 linux/amd64 windows/amd64 windows/arm64}"

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

for target in $TARGETS; do
  os="${target%/*}"
  arch="${target#*/}"
  name="gacha-$os-$arch"
  binary="gacha"
  archive="$name.tar.gz"
  if [ "$os" = "windows" ]; then
    binary="gacha.exe"
    archive="$name.zip"
    command -v zip >/dev/null 2>&1 || {
      echo "scripts/build-release.sh requires zip for Windows artifacts" >&2
      exit 1
    }
  fi
  echo "Building $name"
  GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o "$OUT_DIR/$name/$binary" ./cmd/gacha
  if [ "$os" = "windows" ]; then
    (cd "$OUT_DIR/$name" && zip -q "../$archive" "$binary")
  else
    tar -C "$OUT_DIR/$name" -czf "$OUT_DIR/$archive" "$binary"
  fi
  rm -rf "$OUT_DIR/$name"
done

if command -v shasum >/dev/null 2>&1; then
  (cd "$OUT_DIR" && find . -maxdepth 1 \( -name '*.tar.gz' -o -name '*.zip' \) -type f -print | sed 's#^\./##' | sort | xargs shasum -a 256 > checksums.txt)
elif command -v sha256sum >/dev/null 2>&1; then
  (cd "$OUT_DIR" && find . -maxdepth 1 \( -name '*.tar.gz' -o -name '*.zip' \) -type f -print | sed 's#^\./##' | sort | xargs sha256sum > checksums.txt)
fi

echo "Release artifacts written to $OUT_DIR"
