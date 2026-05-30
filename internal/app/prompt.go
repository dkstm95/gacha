package app

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

//go:embed assets/plugins/gacha/platforms/generic/system-prompt.md
//go:embed assets/plugins/gacha/templates/investment-report.md
//go:embed assets/plugins/gacha/workflows/*.md
var embedded embed.FS

func printPrompt(query []string) error {
	prompt, err := buildPrompt(query)
	if err != nil {
		return err
	}
	fmt.Println(prompt)
	return nil
}

func runQuery(args []string) error {
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
	prompt, err := buildPrompt(query)
	if err != nil {
		return err
	}
	if err := runAgent(prompt, dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "Could not use OpenCode automatically: %v\n", err)
		fmt.Fprintln(os.Stderr, "Falling back to a prompt you can paste into any web AI.")
		fmt.Fprintln(os.Stderr)
		fmt.Println(prompt)
	}
	return nil
}

func buildPrompt(queryParts []string) (string, error) {
	system, err := readEmbedded("assets/plugins/gacha/platforms/generic/system-prompt.md")
	if err != nil {
		return "", err
	}
	template, err := readEmbedded("assets/plugins/gacha/templates/investment-report.md")
	if err != nil {
		return "", err
	}
	workflows, err := readWorkflows()
	if err != nil {
		return "", err
	}
	workflow := `# gacha auto

Classify the user's request into discover, select, entry, exit, portfolio, or journal, then follow the matching gacha workflow. If the request is ambiguous, choose the safest interpretation and state your assumption.`
	query := strings.TrimSpace(strings.Join(queryParts, " "))
	if query == "" {
		query = "(No additional user request supplied.)"
	}
	lang := responseLanguage(query)
	sections := []string{
		strings.TrimSpace(system),
		strings.TrimSpace(workflow),
		"Workflow library:\n" + strings.TrimSpace(workflows),
		"User request:\n" + query,
		"Response language:\nWrite the final report in " + string(lang) + ". Keep source names, ticker symbols, numbers, and URLs unchanged.",
		"Report template:\n" + strings.TrimSpace(template),
		`Hard requirements:
- Always use current web search or current market-data tools before analysis, even if the user does not ask for latest/current/recent data.
- If fresh data cannot be verified, do not make a recommendation.
- Include data freshness, source links, risks, Devil's Advocate, action conditions, monitoring plan, and provenance.
- Write user-facing explanations, headings, and action conditions in the response language above.
- Do not execute trades. The final decision remains with the user.`,
	}
	return strings.Join(sections, "\n\n"), nil
}

func readWorkflows() (string, error) {
	entries, err := fs.ReadDir(embedded, "assets/plugins/gacha/workflows")
	if err != nil {
		return "", err
	}
	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		content, err := readEmbedded("assets/plugins/gacha/workflows/" + entry.Name())
		if err != nil {
			return "", err
		}
		parts = append(parts, strings.TrimSpace(content))
	}
	return strings.Join(parts, "\n\n"), nil
}

func readEmbedded(name string) (string, error) {
	data, err := fs.ReadFile(embedded, name)
	if err != nil {
		return "", err
	}
	return string(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))), nil
}
