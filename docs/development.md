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
VERSION=0.1.30 sh scripts/build-release.sh
gh release create v0.1.30 dist/*.tar.gz dist/*.zip dist/checksums.txt --title "v0.1.30"
```

GitHub Actions templates are available in:

```text
docs/github-actions/
```

They are kept as templates until the repository token has permission to push workflow files.

## Agent Assets

```text
internal/agent/system-prompt.md
internal/agent/workflows/
internal/agent/templates/
```

These files are embedded into the `gacha` binary with `go:embed`. Treat
`internal/agent/` as the source of truth for the runtime prompt, workflows, and
report template.
