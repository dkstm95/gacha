package app

import (
	"strings"
)

type choiceKind string

const (
	choiceLanguage choiceKind = "language"
	choiceSettings choiceKind = "settings"
	choiceTheme    choiceKind = "theme"
	choiceReport   choiceKind = "report"
)

type pendingChoice struct {
	Kind     choiceKind
	Title    string
	Intro    string
	Options  []choiceOption
	Selected int
	Footer   string
}

type choiceOption struct {
	Label       string
	Value       string
	Description string
}

func (c pendingChoice) Render(text uiText) string {
	return c.RenderWidth(text, 78)
}

func (c pendingChoice) RenderWidth(text uiText, width int) string {
	if c.Kind == choiceTheme {
		return renderThemeChoice(c, text, width)
	}
	lines := []string{titleStyle.Render(c.Title)}
	if strings.TrimSpace(c.Intro) != "" {
		lines = append(lines, wrapLine(c.Intro, max(24, width-4)))
	}
	lines = append(lines, mutedStyle.Render(text.ChoiceHint), "")
	for i, option := range c.Options {
		marker := " "
		if i == c.Selected {
			marker = "›"
		}
		line := bulletStyle.Render(marker) + " " + actionNameStyle.Render(option.Label)
		if strings.TrimSpace(option.Description) != "" {
			line += " " + mutedStyle.Render(wrapIndented(option.Description, max(18, width-20), "    "))
		}
		lines = append(lines, line)
	}
	if strings.TrimSpace(c.Footer) != "" {
		lines = append(lines, "", mutedStyle.Render(c.Footer))
	}
	return strings.Join(lines, "\n")
}

func renderReportActions(text uiText) string {
	choice := pendingChoice{
		Kind:    choiceReport,
		Title:   text.ReportActionsTitle,
		Intro:   text.ReportChoiceIntro,
		Options: reportChoiceOptions(text),
		Footer:  text.NewQuestionAction,
	}
	return choice.Render(text)
}

func reportChoiceOptions(text uiText) []choiceOption {
	options := make([]choiceOption, 0, len(text.ReportActions))
	for _, action := range text.ReportActions {
		options = append(options, choiceOption{
			Label:       action.Label,
			Value:       action.Key,
			Description: action.Description,
		})
	}
	return options
}

func themeChoiceOptions(text uiText) []choiceOption {
	themes := availableThemes()
	options := make([]choiceOption, 0, len(themes))
	for _, theme := range themes {
		options = append(options, choiceOption{
			Label:       themeLabel(theme, text),
			Value:       theme.Name,
			Description: themeDescription(theme, text),
		})
	}
	return options
}

func selectedChoiceIndex(options []choiceOption, value string) int {
	for i, option := range options {
		if strings.EqualFold(option.Value, value) || strings.EqualFold(option.Label, value) {
			return i
		}
	}
	return 0
}
