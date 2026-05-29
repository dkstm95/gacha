# Development

## Build From Source

```bash
git clone https://github.com/dkstm95/investiq.git
cd investiq
go test ./...
go build -o investiq ./cmd/investiq
./investiq doctor
```

## Release

```bash
VERSION=0.1.6 sh scripts/build-release.sh
gh release create v0.1.6 dist/*.tar.gz dist/checksums.txt --title "v0.1.6"
```

GitHub Actions templates are available in:

```text
docs/github-actions/
```

They are kept as templates until the repository token has permission to push workflow files.

## Codex Plugin Assets

```text
.agents/plugins/marketplace.json
plugins/investiq/.codex-plugin/plugin.json
plugins/investiq/skills/investiq/SKILL.md
```
