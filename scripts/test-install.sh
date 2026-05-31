#!/usr/bin/env sh
set -eu

root_dir="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

fakebin="$tmpdir/bin"
release_dir="$tmpdir/release"
install_dir="$tmpdir/install"
mkdir -p "$fakebin" "$release_dir" "$install_dir"

cat > "$tmpdir/gacha" <<'SCRIPT'
#!/usr/bin/env sh
if [ "${1:-}" = "version" ]; then
  echo "9.9.9"
else
  echo "fake gacha"
fi
SCRIPT
chmod +x "$tmpdir/gacha"
tar -C "$tmpdir" -czf "$release_dir/gacha-darwin-arm64.tar.gz" gacha

if command -v shasum >/dev/null 2>&1; then
  checksum="$(shasum -a 256 "$release_dir/gacha-darwin-arm64.tar.gz" | awk '{ print $1 }')"
elif command -v sha256sum >/dev/null 2>&1; then
  checksum="$(sha256sum "$release_dir/gacha-darwin-arm64.tar.gz" | awk '{ print $1 }')"
else
  echo "scripts/test-install.sh requires shasum or sha256sum" >&2
  exit 1
fi
printf '%s  gacha-darwin-arm64.tar.gz\n' "$checksum" > "$release_dir/checksums.txt"

cat > "$fakebin/uname" <<'SCRIPT'
#!/usr/bin/env sh
case "${1:-}" in
  -s) echo Darwin ;;
  -m) echo arm64 ;;
  *) exit 1 ;;
esac
SCRIPT
chmod +x "$fakebin/uname"

cat > "$fakebin/curl" <<SCRIPT
#!/usr/bin/env sh
out=""
url=""
while [ "\$#" -gt 0 ]; do
  case "\$1" in
    -o)
      shift
      out="\$1"
      ;;
    http*)
      url="\$1"
      ;;
  esac
  shift
done
case "\$url" in
  */gacha-darwin-arm64.tar.gz)
    cp "$release_dir/gacha-darwin-arm64.tar.gz" "\$out"
    ;;
  */checksums.txt)
    cp "$release_dir/checksums.txt" "\$out"
    ;;
  *)
    echo "unexpected url: \$url" >&2
    exit 1
    ;;
esac
SCRIPT
chmod +x "$fakebin/curl"

PATH="$fakebin:/usr/bin:/bin" GACHA_INSTALL_DIR="$install_dir" GACHA_VERSION="v9.9.9" sh "$root_dir/install.sh" >"$tmpdir/install-ok.out"
"$install_dir/gacha" version | grep -qx "9.9.9"
test -L "$install_dir/gch"

printf '0  gacha-darwin-arm64.tar.gz\n' > "$release_dir/checksums.txt"
if PATH="$fakebin:/usr/bin:/bin" GACHA_INSTALL_DIR="$tmpdir/bad-install" GACHA_VERSION="v9.9.9" sh "$root_dir/install.sh" >"$tmpdir/bad.out" 2>"$tmpdir/bad.err"; then
  echo "install succeeded with an invalid checksum" >&2
  exit 1
fi
grep -q "checksum mismatch" "$tmpdir/bad.err"
