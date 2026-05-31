package app

import "strings"

func reportContextFromMarkdown(report string, fallback string) []string {
	if strings.TrimSpace(report) == "" {
		return []string{fallback}
	}
	candidates := []struct {
		match string
		label string
	}{
		{match: "bottom line", label: "Bottom line"},
		{match: "결론", label: "결론"},
		{match: "decision rules", label: "Decision rules"},
		{match: "행동 기준", label: "행동 기준"},
		{match: "biggest risks", label: "Risks"},
		{match: "risks", label: "Risks"},
		{match: "리스크", label: "리스크"},
		{match: "data check", label: "Data check"},
		{match: "데이터", label: "데이터"},
		{match: "source", label: "Sources"},
		{match: "출처", label: "출처"},
		{match: "반대 의견", label: "반대 의견"},
	}
	var found []string
	seen := map[string]bool{}
	for _, line := range strings.Split(stripANSI(report), "\n") {
		normalized := strings.ToLower(strings.TrimSpace(strings.TrimLeft(line, "#0123456789. ")))
		for _, candidate := range candidates {
			if strings.Contains(normalized, candidate.match) && !seen[candidate.label] {
				found = append(found, candidate.label)
				seen[candidate.label] = true
			}
		}
		if len(found) >= 4 {
			break
		}
	}
	if len(found) == 0 {
		return []string{fallback}
	}
	return found
}
