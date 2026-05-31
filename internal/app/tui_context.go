package app

import (
	"strings"
)

func (m tuiModel) contextRail(width int) string {
	switch {
	case m.busy:
		return m.researchContext(width)
	case m.mode == m.text.Report:
		return m.reportContext(width)
	default:
		return m.homeContext(width)
	}
}

func (m tuiModel) homeContext(width int) string {
	actions := actionNames(m.text.HomeActions)
	lines := []string{
		sectionStyle.Render(m.text.ContextTitle),
		titleStyle.Render(m.text.ContextTypesTitle),
	}
	for _, action := range actions {
		lines = append(lines, bulletStyle.Render("›")+" "+wrapLine(action, max(12, width-3)))
	}
	return strings.Join(lines, "\n")
}

func (m tuiModel) researchContext(width int) string {
	lines := []string{
		sectionStyle.Render(m.text.ContextTitle),
		titleStyle.Render(m.text.ContextRequestTitle),
		wrapLine(m.query, max(16, width)),
		"",
		titleStyle.Render(m.text.ContextResearchTitle),
	}
	for i, phase := range m.text.ResearchPhases {
		marker := " "
		switch {
		case i < m.phase:
			marker = "✓"
		case i == m.phase:
			marker = "•"
		}
		lines = append(lines, mutedStyle.Render(marker)+" "+wrapLine(phase, max(12, width-3)))
		if i >= 4 {
			break
		}
	}
	lines = append(lines, "", titleStyle.Render(m.text.ContextSourcesTitle), mutedStyle.Render(m.text.ContextSourcesPending))
	return strings.Join(lines, "\n")
}

func (m tuiModel) reportContext(width int) string {
	context := reportContextFromMarkdown(m.report, m.text.ContextReportFallback)
	lines := []string{
		sectionStyle.Render(m.text.ContextTitle),
		titleStyle.Render(m.text.ContextRequestTitle),
		wrapLine(m.query, max(16, width)),
		"",
		titleStyle.Render(m.text.ContextReportTitle),
	}
	for _, item := range context {
		lines = append(lines, bulletStyle.Render("•")+" "+wrapLine(item, max(12, width-3)))
	}
	return strings.Join(lines, "\n")
}

func actionNames(actions []homeAction) []string {
	names := make([]string, 0, len(actions))
	for _, action := range actions {
		names = append(names, action.Name)
	}
	return names
}

func reportContextFromMarkdown(report string, fallback string) []string {
	if strings.TrimSpace(report) == "" {
		return []string{fallback}
	}
	candidates := []struct {
		match string
		label string
	}{
		{match: "bottom line", label: "Bottom line"},
		{match: "decision rules", label: "Decision rules"},
		{match: "biggest risks", label: "Risks"},
		{match: "risks", label: "Risks"},
		{match: "data check", label: "Data check"},
		{match: "source", label: "Sources"},
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
