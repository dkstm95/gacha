# Next Tasks

This file is a handoff note for future sessions.

## Current Priority

1. Localize the user-facing experience.
   - Detect the user's preferred language from terminal/OS locale.
   - Show the TUI and command-line messages in that language.
   - Ask the AI to write reports in the user's language.
   - If the typed question is clearly Korean, answer in Korean even when the terminal locale is English.

## Suggested Next Product Tasks

1. Lock the investment report structure so answers are consistent.
2. Save completed reports as Markdown under a local Gacha data directory.
3. Improve running-state messages in the TUI.
4. Make `gch doctor` more product-facing and less debug-like.
5. Add a lightweight settings screen for model mode and language override.

## Recent Context

- Model selection now uses OpenCode model discovery.
- OpenAI auto mode should prefer the newest base frontier model pattern, such as `gpt-N` or `gpt-N.M`, rather than hard-coding `gpt-5.5`.
- `pro`, `fast`, `spark`, `mini`, and coding-specialized variants should rank behind the base frontier model for investment research.
