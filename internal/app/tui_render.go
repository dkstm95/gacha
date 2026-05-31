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
