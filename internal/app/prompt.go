package app

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

//go:embed assets/plugins/investiq/platforms/generic/system-prompt.md
//go:embed assets/plugins/investiq/templates/investment-report.md
//go:embed assets/plugins/investiq/workflows/*.md
var embedded embed.FS

func printPrompt(mode string, query []string) error {
	prompt, err := buildPrompt(mode, query)
	if err != nil {
		return err
	}
	fmt.Println(prompt)
	return nil
}

func runMode(mode string, args []string) error {
	if !modes[mode] {
		return fmt.Errorf("invalid mode: %s", mode)
	}
	dryRun := false
	query := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			dryRun = true
		default:
			query = append(query, args[i])
		}
	}
	cfg := loadConfig()
	selected := selectPlatform(cfg, cfg.DefaultPlatform)
	platform, ok := cfg.Platforms[selected]
	if !ok {
		return fmt.Errorf("unknown platform: %s", selected)
	}
	prompt, err := buildPrompt(mode, query)
	if err != nil {
		return err
	}
	if err := runPlatform(selected, platform, prompt, dryRun); err != nil {
		if selected == "manual" {
			return err
		}
		fmt.Fprintf(os.Stderr, "Could not use %s automatically: %v\n", platform.Label, err)
		fmt.Fprintln(os.Stderr, "Falling back to a prompt you can paste into any web AI.")
		fmt.Fprintln(os.Stderr)
		return runPlatform("manual", cfg.Platforms["manual"], prompt, dryRun)
	}
	return nil
}

func buildPrompt(mode string, queryParts []string) (string, error) {
	if !modes[mode] {
		return "", fmt.Errorf("unknown mode: %s", mode)
	}
	system, err := readEmbedded("assets/plugins/investiq/platforms/generic/system-prompt.md")
	if err != nil {
		return "", err
	}
	template, err := readEmbedded("assets/plugins/investiq/templates/investment-report.md")
	if err != nil {
		return "", err
	}
	workflow := `# investiq auto

Classify the user's request into discover, select, entry, exit, portfolio, or journal, then follow the matching investiq workflow. If the request is ambiguous, choose the safest interpretation and state your assumption.`
	if mode != "auto" {
		workflowPath := fmt.Sprintf("assets/plugins/investiq/workflows/%s.md", mode)
		workflow, err = readEmbedded(workflowPath)
		if err != nil {
			return "", err
		}
	}
	query := strings.TrimSpace(strings.Join(queryParts, " "))
	if query == "" {
		query = "(No additional user request supplied.)"
	}
	sections := []string{
		strings.TrimSpace(system),
		strings.TrimSpace(workflow),
		"User request:\n" + query,
		"Report template:\n" + strings.TrimSpace(template),
		`Hard requirements:
- Always use current web search or current market-data tools before analysis, even if the user does not ask for latest/current/recent data.
- If fresh data cannot be verified, do not make a recommendation.
- Include data freshness, source links, risks, Devil's Advocate, action conditions, monitoring plan, and provenance.
- Do not execute trades. The final decision remains with the user.`,
	}
	return strings.Join(sections, "\n\n"), nil
}

func readEmbedded(name string) (string, error) {
	data, err := fs.ReadFile(embedded, name)
	if err != nil {
		return "", err
	}
	return string(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))), nil
}
