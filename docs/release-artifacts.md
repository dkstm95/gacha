# Release Artifacts

Gacha releases publish one archive per supported target plus a checksum manifest.

## Archives

Archive names must use this format:

```text
gacha-<goos>-<goarch>.<extension>
```

Current archive extensions:

```text
darwin-amd64  .tar.gz
darwin-arm64  .tar.gz
linux-amd64   .tar.gz
linux-arm64   .tar.gz
windows-amd64 .zip
windows-arm64 .zip
```

The Unix archives must contain an executable named `gacha` at the archive root.
The Windows archives must contain `gacha.exe`.

## Checksums

Every release must include `checksums.txt` beside the archives. Each line uses
the standard `sha256sum` shape:

```text
<sha256>  <artifact-name>
```

The installer and `gch update` treat this manifest as part of the release
contract. They must fail closed when the manifest is missing, the target asset is
not listed, or the downloaded checksum does not match.

## Automation

The canonical GitHub Actions workflow is:

```text
.github/workflows/release.yml
```

Do not keep duplicate workflow templates under `docs/`. If the release workflow
changes, update the canonical workflow directly and keep this artifact contract
in sync.
