package agent

import "embed"

// FS contains the prompt assets used by the gacha runtime.
//
//go:embed system-prompt.md
//go:embed templates/investment-report.md
//go:embed workflows/*.md
var FS embed.FS
