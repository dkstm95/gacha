package app

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func renderHeader(width int, version string) string {
	left := brandStyle.Render(" GACHA ")
	right := mutedStyle.Render("v" + version)
	if width < 58 {
		right = ""
	}
	line := strings.Repeat("─", max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-2))
	return lipgloss.JoinHorizontal(lipgloss.Center, left, mutedStyle.Render(line), right)
}

func renderStatus(width int, status string, runtime string, mode string, busy bool, spin string, text uiText) string {
	indicator := "●"
	if busy {
		indicator = spin
	}
	items := []string{accentStyle.Render(indicator + " " + status), mutedStyle.Render(text.StatusMode + mode)}
	if width >= 92 {
		items = append(items, mutedStyle.Render(text.StatusRuntime+runtime))
	}
	left := strings.Join(items, "   ")
	right := renderFooter(width, text)
	gap := width - 4 - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 3 {
		right = renderFooter(72, text)
		gap = width - 4 - lipgloss.Width(left) - lipgloss.Width(right)
	}
	if gap < 3 {
		return statusStyle.Width(width - 2).Render(left)
	}
	return statusStyle.Width(width - 2).Render(left + strings.Repeat(" ", gap) + right)
}

func renderInput(width int, input string) string {
	return inputStyle.Width(width - 2).Render(input)
}

func renderFooter(width int, text uiText) string {
	footer := text.Footer
	if width < 86 {
		footer = text.FooterShort
	}
	return faintStyle.Render(footer)
}

func renderHomeSection(title string, items []string, marker string, width int) string {
	lines := []string{sectionStyle.Render(title)}
	for _, item := range items {
		wrapped := wrapIndented(item, max(16, width-4), "  ")
		lines = append(lines, bulletStyle.Render(marker)+" "+wrapped)
	}
	return strings.Join(lines, "\n")
}

func renderHomeActions(title string, actions []homeAction, width int) string {
	lines := []string{sectionStyle.Render(title)}
	nameWidth := 0
	for _, action := range actions {
		nameWidth = max(nameWidth, lipgloss.Width(action.Name))
	}
	for _, action := range actions {
		name := actionNameStyle.Render(padRight(action.Name, nameWidth))
		promptWidth := max(18, width-nameWidth-4)
		prompt := wrapIndented(action.Prompt, promptWidth, strings.Repeat(" ", nameWidth+3))
		lines = append(lines, bulletStyle.Render("›")+" "+name+" "+prompt)
	}
	return strings.Join(lines, "\n")
}

func renderHomeNote(value string, width int) string {
	return noteStyle.Width(max(24, width-2)).Render(wrapLine(value, max(20, width-6)))
}
