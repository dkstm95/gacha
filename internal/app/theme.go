package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type uiTheme struct {
	Name    string
	Palette uiPalette
}

type uiPalette struct {
	BrandFg       lipgloss.TerminalColor
	BrandBg       lipgloss.TerminalColor
	Accent        lipgloss.TerminalColor
	Hero          lipgloss.TerminalColor
	SectionFg     lipgloss.TerminalColor
	SectionBg     lipgloss.TerminalColor
	Muted         lipgloss.TerminalColor
	Faint         lipgloss.TerminalColor
	Warning       lipgloss.TerminalColor
	ActionFg      lipgloss.TerminalColor
	ActionBg      lipgloss.TerminalColor
	Border        lipgloss.TerminalColor
	PanelBorder   lipgloss.TerminalColor
	InputBorder   lipgloss.TerminalColor
	CalloutBorder lipgloss.TerminalColor
}

type uiStyles struct {
	brand      lipgloss.Style
	accent     lipgloss.Style
	title      lipgloss.Style
	hero       lipgloss.Style
	section    lipgloss.Style
	muted      lipgloss.Style
	faint      lipgloss.Style
	warning    lipgloss.Style
	bullet     lipgloss.Style
	key        lipgloss.Style
	actionName lipgloss.Style
	note       lipgloss.Style
	status     lipgloss.Style
	panel      lipgloss.Style
	input      lipgloss.Style
	callout    lipgloss.Style
}

func configuredTheme() string {
	config, err := configWithDefaults()
	if err != nil {
		return themeSettingSystem
	}
	theme, ok := normalizeThemeSetting(config.Theme)
	if !ok {
		return themeSettingSystem
	}
	return theme
}

func setThemeStyles(theme string) {
	current := themeByName(theme)
	styles := makeUIStyles(current.Palette)
	brandStyle = styles.brand
	accentStyle = styles.accent
	titleStyle = styles.title
	heroStyle = styles.hero
	sectionStyle = styles.section
	mutedStyle = styles.muted
	faintStyle = styles.faint
	warningStyle = styles.warning
	bulletStyle = styles.bullet
	keyStyle = styles.key
	actionNameStyle = styles.actionName
	noteStyle = styles.note
	statusStyle = styles.status
	panelStyle = styles.panel
	inputStyle = styles.input
	calloutStyle = styles.callout
}

func makeUIStyles(p uiPalette) uiStyles {
	return uiStyles{
		brand: lipgloss.NewStyle().
			Foreground(p.BrandFg).
			Background(p.BrandBg).
			Bold(true),
		accent:  lipgloss.NewStyle().Foreground(p.Accent).Bold(true),
		title:   lipgloss.NewStyle().Foreground(p.Accent).Bold(true),
		hero:    lipgloss.NewStyle().Foreground(p.Hero).Bold(true),
		section: lipgloss.NewStyle().Foreground(p.SectionFg).Background(p.SectionBg).Bold(true).Padding(0, 1),
		muted:   lipgloss.NewStyle().Foreground(p.Muted),
		faint:   lipgloss.NewStyle().Foreground(p.Faint),
		warning: lipgloss.NewStyle().Foreground(p.Warning).Bold(true),
		bullet:  lipgloss.NewStyle().Foreground(p.Accent).Bold(true),
		key:     lipgloss.NewStyle().Foreground(p.SectionFg).Background(p.SectionBg).Bold(true).Padding(0, 1),
		actionName: lipgloss.NewStyle().
			Foreground(p.ActionFg).
			Background(p.ActionBg).
			Bold(true).
			Padding(0, 1),
		note: lipgloss.NewStyle().
			Foreground(p.Faint).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(p.Border).
			PaddingLeft(1),
		status: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.Border).
			Padding(0, 1),
		panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.PanelBorder).
			Padding(1, 2),
		input: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.InputBorder).
			Padding(0, 1),
		callout: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.CalloutBorder).
			Padding(0, 1),
	}
}

func availableThemes() []uiTheme {
	return []uiTheme{
		{
			Name: themeSettingSystem,
			Palette: uiPalette{
				BrandFg:       lipgloss.Color("230"),
				BrandBg:       lipgloss.AdaptiveColor{Light: "25", Dark: "62"},
				Accent:        lipgloss.AdaptiveColor{Light: "25", Dark: "81"},
				Hero:          lipgloss.AdaptiveColor{Light: "24", Dark: "87"},
				SectionFg:     lipgloss.Color("230"),
				SectionBg:     lipgloss.AdaptiveColor{Light: "25", Dark: "62"},
				Muted:         lipgloss.AdaptiveColor{Light: "240", Dark: "244"},
				Faint:         lipgloss.AdaptiveColor{Light: "242", Dark: "245"},
				Warning:       lipgloss.AdaptiveColor{Light: "130", Dark: "230"},
				ActionFg:      lipgloss.Color("230"),
				ActionBg:      lipgloss.AdaptiveColor{Light: "238", Dark: "238"},
				Border:        lipgloss.AdaptiveColor{Light: "245", Dark: "238"},
				PanelBorder:   lipgloss.AdaptiveColor{Light: "25", Dark: "62"},
				InputBorder:   lipgloss.AdaptiveColor{Light: "25", Dark: "81"},
				CalloutBorder: lipgloss.Color("214"),
			},
		},
		{
			Name: themeSettingDark,
			Palette: uiPalette{
				BrandFg:       lipgloss.Color("230"),
				BrandBg:       lipgloss.Color("23"),
				Accent:        lipgloss.Color("80"),
				Hero:          lipgloss.Color("121"),
				SectionFg:     lipgloss.Color("230"),
				SectionBg:     lipgloss.Color("23"),
				Muted:         lipgloss.Color("250"),
				Faint:         lipgloss.Color("244"),
				Warning:       lipgloss.Color("230"),
				ActionFg:      lipgloss.Color("230"),
				ActionBg:      lipgloss.Color("238"),
				Border:        lipgloss.Color("238"),
				PanelBorder:   lipgloss.Color("23"),
				InputBorder:   lipgloss.Color("80"),
				CalloutBorder: lipgloss.Color("214"),
			},
		},
		{
			Name: themeSettingLight,
			Palette: uiPalette{
				BrandFg:       lipgloss.Color("230"),
				BrandBg:       lipgloss.Color("25"),
				Accent:        lipgloss.Color("25"),
				Hero:          lipgloss.Color("24"),
				SectionFg:     lipgloss.Color("230"),
				SectionBg:     lipgloss.Color("25"),
				Muted:         lipgloss.Color("240"),
				Faint:         lipgloss.Color("242"),
				Warning:       lipgloss.Color("130"),
				ActionFg:      lipgloss.Color("230"),
				ActionBg:      lipgloss.Color("238"),
				Border:        lipgloss.Color("245"),
				PanelBorder:   lipgloss.Color("25"),
				InputBorder:   lipgloss.Color("25"),
				CalloutBorder: lipgloss.Color("130"),
			},
		},
		{
			Name: themeSettingGacha,
			Palette: uiPalette{
				BrandFg:       lipgloss.Color("230"),
				BrandBg:       lipgloss.Color("62"),
				Accent:        lipgloss.Color("81"),
				Hero:          lipgloss.Color("87"),
				SectionFg:     lipgloss.Color("230"),
				SectionBg:     lipgloss.Color("62"),
				Muted:         lipgloss.Color("244"),
				Faint:         lipgloss.Color("245"),
				Warning:       lipgloss.Color("230"),
				ActionFg:      lipgloss.Color("230"),
				ActionBg:      lipgloss.Color("238"),
				Border:        lipgloss.Color("238"),
				PanelBorder:   lipgloss.Color("62"),
				InputBorder:   lipgloss.Color("81"),
				CalloutBorder: lipgloss.Color("214"),
			},
		},
	}
}

func themeByName(name string) uiTheme {
	normalized, ok := normalizeThemeSetting(name)
	if !ok {
		normalized = themeSettingSystem
	}
	for _, theme := range availableThemes() {
		if theme.Name == normalized {
			return theme
		}
	}
	return availableThemes()[0]
}

func themeContent(text uiText) string {
	active := configuredTheme()
	choice := pendingChoice{
		Kind:     choiceTheme,
		Title:    text.ThemeTitle,
		Intro:    fmt.Sprintf("%s %s", text.ThemeActive, themeLabel(themeByName(active), text)),
		Options:  themeChoiceOptions(text),
		Selected: selectedChoiceIndex(themeChoiceOptions(text), active),
	}
	return choice.Render(text)
}

func renderThemeChoice(choice pendingChoice, text uiText) string {
	lines := []string{titleStyle.Render(choice.Title)}
	if strings.TrimSpace(choice.Intro) != "" {
		lines = append(lines, wrapLine(choice.Intro, 68))
	}
	lines = append(lines, mutedStyle.Render(text.ChoiceHint), "")
	for i, option := range choice.Options {
		marker := " "
		if i == choice.Selected {
			marker = "›"
		}
		lines = append(lines, bulletStyle.Render(marker)+" "+actionNameStyle.Render(option.Label)+" "+mutedStyle.Render(wrapIndented(option.Description, 58, "    ")))
	}
	lines = append(lines, "", sectionStyle.Render(text.ThemePreviewTitle))
	if choice.Selected >= 0 && choice.Selected < len(choice.Options) {
		selected := themeByName(choice.Options[choice.Selected].Value)
		lines = append(lines, "", renderThemePreview(selected, configuredTheme() == selected.Name, text))
	}
	return strings.Join(lines, "\n")
}

func renderThemePreview(theme uiTheme, active bool, text uiText) string {
	styles := makeUIStyles(theme.Palette)
	marker := " "
	if active {
		marker = "✓"
	}
	actionName := "Buy"
	actionPrompt := "Should I buy NVDA now?"
	if len(text.HomeActions) > 0 {
		actionName = text.HomeActions[0].Name
		actionPrompt = text.HomeActions[0].Prompt
	}
	outcome := "Bottom line"
	if len(text.HomeOutcomes) > 0 {
		outcome = text.HomeOutcomes[0]
	}
	previewWidth := 64
	textWidth := previewWidth - 6
	label := styles.brand.Render(" " + themeLabel(theme, text) + " ")
	descriptionWidth := max(20, textWidth-lipgloss.Width(stripANSI(label))-3)
	header := label + " " + styles.muted.Render(wrapIndented(themeDescription(theme, text), descriptionWidth, ""))
	actionLine := styles.section.Render(text.HomeActionsTitle) + " " + styles.actionName.Render(actionName)
	actionPrompt = wrapIndented(actionPrompt, max(24, textWidth-lipgloss.Width(actionName)-3), strings.Repeat(" ", lipgloss.Width(actionName)+3))
	outcomeLine := styles.bullet.Render("›") + " " + styles.title.Render(outcome)
	note := wrapIndented(text.HomeNote, max(24, textWidth-lipgloss.Width(outcome)-3), strings.Repeat(" ", lipgloss.Width(outcome)+3))
	sample := strings.Join([]string{
		actionLine + " " + actionPrompt,
		outcomeLine + " " + styles.faint.Render(note),
		styles.key.Render("enter") + " " + styles.muted.Render(text.ThemeSelectLabel),
	}, "\n")
	return styles.panel.Width(previewWidth).Render(marker + " " + header + "\n" + sample)
}

func themeLabel(theme uiTheme, text uiText) string {
	if label := text.ThemeLabels[theme.Name]; label != "" {
		return label
	}
	return theme.Name
}

func themeDescription(theme uiTheme, text uiText) string {
	if description := text.ThemeDescriptions[theme.Name]; description != "" {
		return description
	}
	return theme.Name
}

var (
	brandStyle      lipgloss.Style
	accentStyle     lipgloss.Style
	titleStyle      lipgloss.Style
	heroStyle       lipgloss.Style
	sectionStyle    lipgloss.Style
	mutedStyle      lipgloss.Style
	faintStyle      lipgloss.Style
	warningStyle    lipgloss.Style
	bulletStyle     lipgloss.Style
	keyStyle        lipgloss.Style
	actionNameStyle lipgloss.Style
	noteStyle       lipgloss.Style
	statusStyle     lipgloss.Style
	panelStyle      lipgloss.Style
	inputStyle      lipgloss.Style
	calloutStyle    lipgloss.Style
)

func init() {
	setThemeStyles(themeSettingSystem)
}
