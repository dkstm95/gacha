# investiq

`investiq` is a standalone investment research agent harness.

It routes your investment question to the best available AI CLI on your machine, while enforcing fresh-data research, source links, risk review, counterarguments, and provenance. If no supported AI CLI is ready, it prints a complete prompt you can paste into any web AI with browsing.

Korean documentation: [docs/ko/README.md](docs/ko/README.md)

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/dkstm95/investiq/main/install.sh | sh
```

The installer downloads a standalone binary from GitHub Releases and installs:

- `investiq`
- `iq`, the short command most users should use

No Node, npm, Python, or Go runtime is required.

If the installer says your install directory is not on `PATH`, run the printed `export PATH=...` command.

## Quick Start

```bash
iq init
iq doctor
investiq
```

Running `investiq` opens the interactive UI:

```text
investiq
Fresh-data investment research agent

Ask an investment question. investiq will classify it and route it automatically.
Type /help for commands, /doctor to check AI platforms, /quit to exit.

iq>
```

Type your question at the `iq>` prompt:

```text
iq> Should I buy NVDA now?
```

You do not need to choose `entry`, `exit`, or a platform. `investiq` classifies the request internally and routes it automatically.

More examples:

```text
iq> What should I invest in for the next 6 to 12 months?
iq> I want exposure to AI infrastructure. Which stocks or ETFs should I compare?
iq> I own TSLA. When should I trim, sell, or stop out?
iq> Review my portfolio: AAPL 35%, NVDA 30%, SGOV 35%.
```

For one-shot usage, you can still run:

```bash
iq "Should I buy NVDA now?"
```

## How It Works

`investiq` checks your locally available AI CLIs and uses the first working platform in this routing order:

```text
Claude Code -> Codex -> OpenCode -> Gemini CLI -> manual prompt
```

If a detected platform fails at runtime, `investiq` falls back to printing a prompt instead of blocking you.

`iq init` writes:

```text
~/.investiq/config.json
```

You can edit that file if you want to change routing priority or command names.

## Update

```bash
iq update
```

`iq update` checks the latest GitHub Release, downloads the matching standalone binary for your OS/CPU, and replaces the installed `investiq` binary.

## Fresh Data Rule

Every investment workflow requires current web or market data, even if the user does not explicitly ask for "latest", "current", or "recent" data.

If fresh data cannot be verified, the AI should not make a recommendation.

Required report elements:

- data freshness
- source links
- current price or latest relevant values
- investment thesis
- valuation or scenario analysis
- risks
- Devil's Advocate
- action conditions
- monitoring plan
- provenance appendix

## What investiq Does Not Do

`investiq` does not:

- execute trades
- guarantee returns
- replace professional financial, tax, or legal advice
- fetch market data by itself in the current version

It composes and routes a strict investment research workflow to an AI platform. The connected AI platform must perform the fresh web or market-data research.

## Build From Source

```bash
git clone https://github.com/dkstm95/investiq.git
cd investiq
go test ./...
go build -o investiq .
./investiq doctor
```

## Release

```bash
VERSION=0.1.2 sh scripts/build-release.sh
gh release create v0.1.2 dist/*.tar.gz dist/checksums.txt --title "v0.1.2"
```

GitHub Actions templates are available in:

```text
docs/github-actions/
```

They are kept as templates until the repository token has permission to push workflow files.

## Codex Plugin Assets

This repository also includes marketplace/plugin assets:

```text
.agents/plugins/marketplace.json
plugins/investiq/.codex-plugin/plugin.json
plugins/investiq/skills/investiq/SKILL.md
```

## License

MIT
