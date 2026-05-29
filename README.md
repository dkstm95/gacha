# investiq

Ask investment questions through a local AI runtime.

`investiq` opens a terminal app for investment research. It uses OpenCode as the local AI runtime, so you can connect ChatGPT, GitHub Copilot, Gemini, OpenAI API, or another supported provider without choosing a platform for every question.

If the runtime is missing, `investiq` can install it for you on first run. If no runtime is ready, it still gives you a complete prompt you can paste into any AI with web browsing.

Korean: [docs/ko/README.md](docs/ko/README.md)

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/dkstm95/investiq/main/install.sh | sh
```

The installer adds two commands:

- `investiq`
- `iq`

Use `investiq` to open the app. Use `iq` if you want a shorter command.

No Node, npm, Python, or Go setup is needed for `investiq`.

On first run, `investiq` may ask to install OpenCode runtime and connect an AI provider. This runtime runs behind the InvestIQ UI.

If the installer prints an `export PATH=...` line, run it once in your terminal.

## Start

```bash
investiq
```

You will see:

```text
INVESTIQ
Fresh-data investment research for your AI tools

+------------------------------------------------------------+
| Ask a question. investiq will classify it automatically.   |
| It always asks the AI to use current web or market data.   |
+------------------------------------------------------------+

Ask >
```

Type a question:

```text
Ask > Should I buy NVDA now?
```

You do not need to pick a mode or choose an AI platform. `investiq` handles that for you.

## Example Questions

```text
Ask > What should I invest in for the next 6 to 12 months?
Ask > I want exposure to AI infrastructure. Which stocks or ETFs should I compare?
Ask > I own TSLA. When should I trim, sell, or stop out?
Ask > Review my portfolio: AAPL 35%, NVDA 30%, SGOV 35%.
```

You can also ask one question without opening the app:

```bash
iq "Should I buy NVDA now?"
```

## Setup Check

Run this if you want to see whether the local AI runtime is ready:

```bash
iq doctor
```

`investiq` uses this route:

```text
OpenCode runtime -> copy/paste prompt
```

If OpenCode is missing or no provider is connected, run:

```bash
iq setup
```

`iq setup` installs the runtime if needed, then starts provider login. You can connect ChatGPT, GitHub Copilot, Gemini, OpenAI API, or another OpenCode-supported provider.

After setup, `investiq` keeps the investment workflow and results inside the InvestIQ UI.

If the runtime fails, `investiq` falls back to a prompt you can paste into a web AI.

## Update

```bash
iq update
```

This downloads the right binary for your computer and replaces the old one.

## Fresh Data

Investment information changes quickly. `investiq` always tells the AI to check current web or market data, even if you do not ask for "latest" data.

If current data cannot be checked, the AI should not make a recommendation.

A good answer should include:

- data date and time
- source links
- current price or latest numbers
- main idea
- risks
- opposite view
- buy, hold, sell, or watch conditions
- what to monitor next

## Important Limits

`investiq` does not:

- place trades
- promise returns
- replace professional financial, tax, or legal advice
- fetch market data by itself yet

It prepares a strict research workflow and sends it to an AI tool. The AI tool must do the current web or market-data research.

## Developers

Development notes: [docs/development.md](docs/development.md)

## License

MIT
