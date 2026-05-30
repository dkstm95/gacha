# Next Tasks

This file is a handoff note for future sessions.

## Completed

1. Localize the user-facing experience.
   - Detect the user's preferred language from terminal/OS locale.
   - Show the TUI and command-line messages in that language.
   - Ask the AI to write reports in the user's language.
   - If the typed question is clearly Korean, answer in Korean even when the terminal locale is English.

2. Lock the investment report structure so answers are consistent.
   - Keep a fixed 14-section report contract.
   - Require the same section order for every answer.
   - Require explicit "Not applicable" handling instead of silently dropping sections.
   - Require tables for candidates, price zones, scenarios, action conditions, monitoring, and provenance.

3. Save completed reports as Markdown under a local Gacha data directory.
   - Ask before saving and save only when the user agrees.
   - Save only completed AI reports, not dry-run prompts or paste-fallback prompts.
   - Store reports under the user's Gacha data directory.
   - Include the original question at the bottom of each report.

4. Improve running-state messages in the TUI.
   - Show rotating research phases while a report is running.
   - Localize the phase text for Korean and English.

5. Make `gch doctor` more product-facing and less debug-like.
   - Lead with overall readiness, runtime, provider, model, and report path.
   - Move lower-level command/auth/model source details below the summary.
   - Localize the CLI doctor output for Korean terminals.

## Current Priority

1. Add a lightweight settings screen for model mode and language override.

## Suggested Next Product Tasks

1. Add a lightweight settings screen for model mode and language override.

## Recent Context

- Model selection now uses OpenCode model discovery.
- OpenAI auto mode should prefer the newest base frontier model pattern, such as `gpt-N` or `gpt-N.M`, rather than hard-coding `gpt-5.5`.
- `pro`, `fast`, `spark`, `mini`, and coding-specialized variants should rank behind the base frontier model for investment research.
