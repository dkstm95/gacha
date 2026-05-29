#!/usr/bin/env sh
set -eu

REPO="${GACHA_REPO:-dkstm95/gacha}"
VERSION="${GACHA_VERSION:-latest}"
INSTALL_DIR="${GACHA_INSTALL_DIR:-$HOME/.local/bin}"
BIN_NAME="gacha"
ALIAS_NAME="gch"

need() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "install.sh requires $1" >&2
    exit 1
  }
}

detect_target() {
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"

  case "$os" in
    darwin) os="darwin" ;;
    linux) os="linux" ;;
    *) echo "unsupported OS: $os" >&2; exit 1 ;;
  esac

  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *) echo "unsupported architecture: $arch" >&2; exit 1 ;;
  esac

  echo "$os-$arch"
}

need curl
need tar

target="$(detect_target)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

if [ "$VERSION" = "latest" ]; then
  url="https://github.com/$REPO/releases/latest/download/gacha-$target.tar.gz"
else
  url="https://github.com/$REPO/releases/download/$VERSION/gacha-$target.tar.gz"
fi

mkdir -p "$INSTALL_DIR"
echo "Downloading $url"
curl -fsSL "$url" -o "$tmpdir/gacha.tar.gz"
tar -xzf "$tmpdir/gacha.tar.gz" -C "$tmpdir"
install -m 0755 "$tmpdir/gacha" "$INSTALL_DIR/$BIN_NAME"

echo "Installed $BIN_NAME to $INSTALL_DIR/$BIN_NAME"
if command -v "$ALIAS_NAME" >/dev/null 2>&1; then
  alias_path="$(command -v "$ALIAS_NAME")"
  if [ "$alias_path" = "$INSTALL_DIR/$ALIAS_NAME" ]; then
    ln -sf "$BIN_NAME" "$INSTALL_DIR/$ALIAS_NAME"
    echo "Updated short alias $ALIAS_NAME at $INSTALL_DIR/$ALIAS_NAME"
  else
    echo "Skipped short alias $ALIAS_NAME because it already exists at $alias_path"
  fi
else
  ln -sf "$BIN_NAME" "$INSTALL_DIR/$ALIAS_NAME"
  echo "Installed short alias $ALIAS_NAME to $INSTALL_DIR/$ALIAS_NAME"
fi
case ":$PATH:" in
  *":$INSTALL_DIR:"*)
    echo "$INSTALL_DIR is already on PATH."
    echo "You can now run:"
    echo "  gch version"
    echo "  gch setup"
    ;;
  *)
    echo "$INSTALL_DIR is not on PATH."
    echo "Run this now:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    echo
    echo "Add this to your shell profile:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac

"$INSTALL_DIR/$BIN_NAME" version
