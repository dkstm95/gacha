# Development

## Build From Source

```bash
git clone https://github.com/dkstm95/gacha.git
cd gacha
go test ./...
go build -o gacha ./cmd/gacha
./gacha doctor
ln -sf gacha gch
./gch doctor
```

## Release

```bash
VERSION=0.1.23 sh scripts/build-release.sh
gh release create v0.1.23 dist/*.tar.gz dist/checksums.txt --title "v0.1.23"
```

GitHub Actions templates are available in:

```text
docs/github-actions/
```

They are kept as templates until the repository token has permission to push workflow files.

## Codex Plugin Assets

```text
.agents/plugins/marketplace.json
plugins/gacha/.codex-plugin/plugin.json
plugins/gacha/skills/gacha/SKILL.md
```
