#!/usr/bin/env sh
set -eu

REPO="${INVESTIQ_REPO:-dkstm95/investiq}"
VERSION="${INVESTIQ_VERSION:-latest}"
INSTALL_DIR="${INVESTIQ_INSTALL_DIR:-$HOME/.local/bin}"
BIN_NAME="investiq"

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
  url="https://github.com/$REPO/releases/latest/download/investiq-$target.tar.gz"
else
  url="https://github.com/$REPO/releases/download/$VERSION/investiq-$target.tar.gz"
fi

mkdir -p "$INSTALL_DIR"
echo "Downloading $url"
curl -fsSL "$url" -o "$tmpdir/investiq.tar.gz"
tar -xzf "$tmpdir/investiq.tar.gz" -C "$tmpdir"
install -m 0755 "$tmpdir/investiq" "$INSTALL_DIR/$BIN_NAME"

echo "Installed $BIN_NAME to $INSTALL_DIR/$BIN_NAME"
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo "Add this to your shell profile if needed:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac

"$INSTALL_DIR/$BIN_NAME" version
