# gacha

Ask investment questions through a local AI runtime.

`gacha` opens a terminal app for investment research. It uses OpenCode as the local AI runtime, so you can connect ChatGPT, GitHub Copilot, Gemini, OpenAI API, or another supported provider without choosing a platform for every question.

If the runtime is missing, `gacha` can install it for you on first run on macOS and Linux. On Windows, install OpenCode separately first. If no runtime is ready, `gacha` still gives you a complete prompt you can paste into any AI with web browsing.

Korean: [docs/ko/README.md](docs/ko/README.md)

## Install

### macOS and Linux

```bash
curl -fsSL https://raw.githubusercontent.com/dkstm95/gacha/main/install.sh | sh
```

The macOS/Linux installer adds the main command and a short alias:

- `gacha`
- `gch`

Use `gch` for day-to-day use. Use `gacha` when you want the full command name.

No Node, npm, Python, or Go setup is needed for `gacha`.

On first run, `gacha` may ask to install OpenCode runtime and connect an AI provider on macOS and Linux. This runtime runs behind the Gacha UI.

If the installer prints an `export PATH=...` line, run it once in your terminal.

### Windows

Download `gacha-windows-amd64.zip` from the latest release:

```text
https://github.com/dkstm95/gacha/releases/latest
```

Unzip it, move `gacha.exe` into a folder on your `PATH`, then open a new PowerShell or Windows Terminal window and verify it:

```powershell
gacha version
gacha setup
```

If you want the short `gch` command too, copy `gacha.exe` as `gch.exe` in the same folder:

```powershell
Copy-Item gacha.exe gch.exe
```

Windows automatic OpenCode setup is not supported yet. Install OpenCode from its Windows documentation, make sure `opencode` is on your `PATH`, then run `gacha setup` to connect a provider.

## Start

```bash
gch
```

You will see a terminal workspace with a left-side context rail and a main research area:

![Gacha TUI home screen](docs/assets/gacha-home.png)

Type a question:

```text
Ask > Should I buy NVDA now?
```

You do not need to pick a model or choose an AI platform. `gacha` handles routing through OpenCode. On wide terminals, the context rail shows recent state and decision types; on smaller terminals, it collapses so the workspace and prompt stay focused.

## Example Questions

```text
Ask > What should I invest in for the next 6 to 12 months?
Ask > I want exposure to AI infrastructure. Which stocks or ETFs should I compare?
Ask > I own TSLA. When should I trim, sell, or stop out?
Ask > Review my portfolio: AAPL 35%, NVDA 30%, SGOV 35%.
```

You can also ask one question without opening the app:

```bash
gch "Should I buy NVDA now?"
```

## Setup Check

Run this if you want to see whether the local AI runtime is ready:

```bash
gch doctor
```

The check shows overall readiness, connected provider status, the model Gacha will ask OpenCode to use, and where reports are saved. The default model mode is `auto`.

To view settings in the app, type:

```text
/settings
```

Inside the app, you can set model mode and language:

```text
/model auto
/model opencode-default
/model provider/model
/language auto
/language en
/language ko
```

For scripts, use non-interactive config commands:

```bash
gch config get
gch config set model auto
gch config set model opencode-default
gch config set model provider/model
gch config set language ko
```

In auto mode, Gacha asks OpenCode for the connected provider's current model list, then chooses a strong research model from that list. For OpenAI, it prefers the newest base frontier model pattern, such as `gpt-N` or `gpt-N.M`, instead of hard-coding a specific model name. It pushes small/fast variants such as `mini`, `nano`, `lite`, `flash`, `fast`, `spark`, and coding-specialized variants behind the frontier base model.

For OpenAI OAuth, the ChatGPT subscription route, Gacha treats `pro` models as lower priority because OpenCode can list models that the current ChatGPT account cannot actually run. If OpenCode still rejects a selected model as unsupported for the ChatGPT account, Gacha retries the next discovered candidate.

If Gacha cannot read the model list, it does not guess a hard-coded model. It runs OpenCode without `--model` and lets OpenCode use its default.

Advanced users can override this with:

```bash
GACHA_OPENCODE_MODEL=provider/model gch
```

Or create `~/.config/gacha/config.json`:

```json
{
  "model": "auto",
  "language": "auto"
}
```

Supported values are `auto`, `opencode-default`, or a custom `provider/model`.

## Saved Reports

When the AI runtime completes a report, `gacha` asks whether you want to save it as Markdown.

Default location:

```text
~/.local/share/gacha/reports
```

If `XDG_DATA_HOME` is set, reports are saved under:

```text
$XDG_DATA_HOME/gacha/reports
```

Reports are saved only when you answer yes. Paste-fallback prompts and dry runs are not saved as reports.

## Language

`gacha` detects your terminal language from `GACHA_LANG`, `LANGUAGE`, `LC_ALL`, `LC_MESSAGES`, or `LANG`.

If the language is Korean, the interactive UI is shown in Korean. Reports are also requested in the detected language. If your question contains Korean text, Gacha asks the AI to answer in Korean even when your terminal locale is English.

Set `"language"` in `~/.config/gacha/config.json` to `auto`, `en`, or `ko` to override the terminal language. `GACHA_LANG` still takes precedence.

`gacha` uses this route:

```text
OpenCode runtime -> copy/paste prompt
```

If OpenCode is missing or no provider is connected, run:

```bash
gch setup
```

`gch setup` installs the runtime if needed, then starts provider login. You can connect ChatGPT, GitHub Copilot, Gemini, OpenAI API, or another OpenCode-supported provider.

On Windows, `gch setup` does not install OpenCode automatically. Install OpenCode separately, make sure `opencode` is on your `PATH`, then run `gacha setup` or `gch setup`.

The interactive home screen also shows a setup callout when OpenCode or provider login is missing. After setup, return to `gch` and ask your first investment question.

After setup, `gacha` keeps the investment workflow and results inside the Gacha UI.

If the runtime fails, `gacha` falls back to a prompt you can paste into a web AI.

## Update

```bash
gch update
```

On macOS and Linux, this downloads the right binary for your computer and replaces the old one.

On Windows, self-update is disabled to avoid replacing a running `.exe`. Download the latest `gacha-windows-amd64.zip` or `gacha-windows-arm64.zip`, replace `gacha.exe` in your `PATH`, and open a new terminal.

## Fresh Data

Investment information changes quickly. `gacha` always tells the AI to check current web or market data, even if you do not ask for "latest" data.

If current data cannot be checked, the AI should not make a recommendation.

A good answer starts with a short, plain-language basic report that is complete enough for a first decision. Detailed analysis is included only when it helps verify the decision.

The basic report should include:

- data date and time
- source links
- current price or latest numbers
- plain bottom line
- simple plan
- time horizon, action trigger, thesis-break trigger, and review timing
- biggest risks
- strongest opposite view when relevant
- buy, hold, sell, or watch conditions
- what to monitor next
- a short note that detailed valuation, scenarios, portfolio fit, or source-level evidence can be requested

## Important Limits

`gacha` does not:

- place trades
- promise returns
- replace professional financial, tax, or legal advice
- fetch market data by itself yet

It prepares a strict research workflow and sends it to an AI tool. The AI tool must do the current web or market-data research.

## Developers

Development notes: [docs/development.md](docs/development.md)

## License

MIT
