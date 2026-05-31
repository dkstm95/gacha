#!/usr/bin/env sh
set -eu

root_dir="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

fakebin="$tmpdir/bin"
release_dir="$tmpdir/release"
mkdir -p "$fakebin" "$release_dir"

checksum_file() {
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{ print $1 }'
    return
  fi
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{ print $1 }'
    return
  fi
  echo "scripts/test-install.sh requires shasum or sha256sum" >&2
  exit 1
}

make_release() {
  target="$1"
  version="$2"
  work="$tmpdir/$target"
  mkdir -p "$work"
  cat > "$work/gacha" <<SCRIPT
#!/usr/bin/env sh
if [ "\${1:-}" = "version" ]; then
  echo "$version"
else
  echo "fake gacha $target"
fi
SCRIPT
  chmod +x "$work/gacha"
  tar -C "$work" -czf "$release_dir/gacha-$target.tar.gz" gacha
}

write_checksums() {
  : > "$release_dir/checksums.txt"
  for archive in "$release_dir"/*.tar.gz; do
    printf '%s  %s\n' "$(checksum_file "$archive")" "$(basename "$archive")" >> "$release_dir/checksums.txt"
  done
}

cat > "$fakebin/uname" <<'SCRIPT'
#!/usr/bin/env sh
case "${1:-}" in
  -s) echo "${GACHA_TEST_OS:-Darwin}" ;;
  -m) echo "${GACHA_TEST_ARCH:-arm64}" ;;
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
  */checksums.txt)
    cp "$release_dir/checksums.txt" "\$out"
    ;;
  */gacha-*.tar.gz)
    cp "$release_dir/\$(basename "\$url")" "\$out"
    ;;
  *)
    echo "unexpected url: \$url" >&2
    exit 1
    ;;
esac
SCRIPT
chmod +x "$fakebin/curl"

run_install() {
  os="$1"
  arch="$2"
  install_dir="$3"
  PATH="$fakebin:/usr/bin:/bin" \
    GACHA_TEST_OS="$os" \
    GACHA_TEST_ARCH="$arch" \
    GACHA_INSTALL_DIR="$install_dir" \
    GACHA_VERSION="v9.9.9" \
    sh "$root_dir/install.sh"
}

make_release "darwin-arm64" "9.9.9-darwin"
make_release "linux-amd64" "9.9.9-linux"
write_checksums

darwin_install="$tmpdir/install-darwin"
run_install Darwin arm64 "$darwin_install" >"$tmpdir/darwin.out"
"$darwin_install/gacha" version | grep -qx "9.9.9-darwin"
test -L "$darwin_install/gch"

linux_install="$tmpdir/install-linux"
run_install Linux x86_64 "$linux_install" >"$tmpdir/linux.out"
"$linux_install/gacha" version | grep -qx "9.9.9-linux"
test -L "$linux_install/gch"

alias_bin="$tmpdir/alias-bin"
mkdir -p "$alias_bin"
cat > "$alias_bin/gch" <<'SCRIPT'
#!/usr/bin/env sh
echo existing alias
SCRIPT
chmod +x "$alias_bin/gch"
alias_install="$tmpdir/install-alias"
PATH="$fakebin:$alias_bin:/usr/bin:/bin" \
  GACHA_TEST_OS=Darwin \
  GACHA_TEST_ARCH=arm64 \
  GACHA_INSTALL_DIR="$alias_install" \
  GACHA_VERSION="v9.9.9" \
  sh "$root_dir/install.sh" >"$tmpdir/alias.out"
grep -q "Skipped short alias" "$tmpdir/alias.out"
test ! -e "$alias_install/gch"

printf '0  gacha-darwin-arm64.tar.gz\n' > "$release_dir/checksums.txt"
if run_install Darwin arm64 "$tmpdir/bad-install" >"$tmpdir/bad.out" 2>"$tmpdir/bad.err"; then
  echo "install succeeded with an invalid checksum" >&2
  exit 1
fi
grep -q "checksum mismatch" "$tmpdir/bad.err"

printf '%s  other-asset.tar.gz\n' "$(checksum_file "$release_dir/gacha-darwin-arm64.tar.gz")" > "$release_dir/checksums.txt"
if run_install Darwin arm64 "$tmpdir/missing-install" >"$tmpdir/missing.out" 2>"$tmpdir/missing.err"; then
  echo "install succeeded with a missing checksum entry" >&2
  exit 1
fi
grep -q "checksum for gacha-darwin-arm64.tar.gz not found" "$tmpdir/missing.err"
