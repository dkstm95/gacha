package app

import (
	"charm.land/lipgloss/v2"
	"fmt"
	"os"
	"strings"
)

func runPrompt(query string) promptRunResult {
	prompt, err := buildPrompt([]string{query})
	if err != nil {
		return promptRunResult{err: err}
	}
	commandPath, ok := resolveCommand(openCodeCommand)
	if !ok || !hasOpenCodeAuth() {
		return promptRunResult{output: prompt}
	}
	output, err := runOpenCodeWithResolution(commandPath, prompt, resolveOpenCodeModel(), false)
	if err != nil {
		return promptRunResult{output: output, err: err}
	}
	report := strings.TrimSpace(output)
	return promptRunResult{output: report, completed: report != ""}
}

func runDetailedPrompt(query string, basicReport string) promptRunResult {
	prompt, err := buildDetailedPrompt(query, basicReport)
	if err != nil {
		return promptRunResult{err: err}
	}
	commandPath, ok := resolveCommand(openCodeCommand)
	if !ok || !hasOpenCodeAuth() {
		return promptRunResult{output: prompt}
	}
	output, err := runOpenCodeWithResolution(commandPath, prompt, resolveOpenCodeModel(), false)
	if err != nil {
		return promptRunResult{output: output, err: err}
	}
	report := strings.TrimSpace(output)
	if report == "" {
		return promptRunResult{}
	}
	return promptRunResult{output: strings.TrimSpace(basicReport) + "\n\n" + report, completed: true}
}

func welcomeContent(version string, text uiText, width int, height int) string {
	return welcomeContentWithColumns(version, text, width, height, true)
}

func welcomeContentWithColumns(version string, text uiText, width int, height int, allowColumns bool) string {
	compact := width < 72 || height < 14
	wide := allowColumns && width >= 104 && height >= 16
	actions := text.HomeActions
	outcomes := text.HomeOutcomes
	if compact {
		actions = actions[:min(3, len(actions))]
		outcomes = outcomes[:min(3, len(outcomes))]
	}

	header := renderHomeHero(text, width, compact)
	actionBlock := renderHomeActions(text.HomeActionsTitle, actions, width)
	outcomeBlock := renderHomeSection(text.HomeOutcomesTitle, outcomes, "•", width)
	if wide {
		leftWidth := max(38, (width*54)/100)
		rightWidth := max(30, width-leftWidth-4)
		actionBlock = renderHomeActions(text.HomeActionsTitle, actions, leftWidth)
		outcomeBlock = renderHomeSection(text.HomeOutcomesTitle, outcomes, "•", rightWidth)
	}

	blocks := []string{header}
	if onboarding := onboardingContent(text, width, setupReadiness()); onboarding != "" {
		blocks = append(blocks, "", onboarding)
	}
	if wide {
		leftWidth := max(38, (width*54)/100)
		rightWidth := max(30, width-leftWidth-4)
		blocks = append(blocks, "", lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(leftWidth).Render(actionBlock),
			"    ",
			lipgloss.NewStyle().Width(rightWidth).Render(outcomeBlock),
		))
	} else {
		blocks = append(blocks, "", actionBlock, "", outcomeBlock)
	}
	if !compact {
		blocks = append(blocks, "", renderHomeNote(text.HomeNote, width))
	}
	blocks = append(blocks, "", faintStyle.Render("v"+version))
	return strings.Join(blocks, "\n")
}

func renderHomeHero(text uiText, width int, compact bool) string {
	title := heroStyle.Render(text.HomeTitle)
	bodyWidth := width
	if !compact {
		bodyWidth = max(40, width-6)
	}
	subtitle := mutedStyle.Render(wrapLine(text.HomeSubtitle, bodyWidth))
	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle)
}

func onboardingContent(text uiText, width int, state setupState) string {
	if state == setupReady {
		return ""
	}
	titleIndex := 0
	bodyIndex := 1
	actionIndex := 2
	if state == setupProviderMissing {
		titleIndex = 3
		bodyIndex = 4
		actionIndex = 5
	}
	lines := []string{
		warningStyle.Render(text.Onboarding[titleIndex]),
		wrapLine(text.Onboarding[bodyIndex], max(20, width-4)),
		mutedStyle.Render(wrapLine(text.Onboarding[actionIndex], max(20, width-4))),
	}
	return calloutStyle.Width(max(24, width-2)).Render(strings.Join(lines, "\n"))
}

type setupState int

const (
	setupReady setupState = iota
	setupRuntimeMissing
	setupProviderMissing
)

func setupReadiness() setupState {
	if _, ok := resolveCommand(openCodeCommand); !ok {
		return setupRuntimeMissing
	}
	if !hasOpenCodeAuth() {
		return setupProviderMissing
	}
	return setupReady
}

func researchingContent(query string, text uiText) string {
	lines := text.Research(query)
	return strings.Join([]string{
		titleStyle.Render(lines[0]),
		lines[1],
		lines[2],
		"",
		titleStyle.Render(lines[3]),
		lines[4],
		lines[5],
		lines[6],
		lines[7],
		lines[8],
		"",
		mutedStyle.Render(lines[9]),
	}, "\n")
}

func helpContent(text uiText) string {
	lines := []string{titleStyle.Render(text.HelpLines[0])}
	for _, line := range text.HelpLines[1:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			lines = append(lines, "")
			continue
		}
		command := fields[0]
		description := strings.TrimSpace(strings.TrimPrefix(line, command))
		lines = append(lines, padRight(command, 10)+wrapIndented(description, 56, strings.Repeat(" ", 10)))
	}
	return strings.Join(lines, "\n")
}

func unknownCommandContent(value string, text uiText) string {
	return fmt.Sprintf(text.UnknownCommand, value) + "\n\n" + helpContent(text)
}

func doctorContent(text uiText) string {
	status := text.Missing
	if hasRunnableCommand(openCodeCommand) {
		status = text.Ready
		if !hasOpenCodeAuth() {
			status = text.LoginRequired
		}
	}
	lines := []string{
		titleStyle.Render(text.RuntimeTitle),
		fmt.Sprintf("OpenCode: %s", status),
		fmt.Sprintf("Command:  %s", openCodeCommand),
		fmt.Sprintf("Auth:     %s", openCodeAuthPath()),
	}
	if resolved, ok := resolveCommand(openCodeCommand); ok {
		lines = append(lines, fmt.Sprintf("Resolved: %s", resolved))
	}
	lines = append(lines, fmt.Sprintf("Model:    %s", modelDescription(resolveOpenCodeModel())))
	if hasOpenCodeAuth() {
		if providers, err := openCodeAuthList(); err == nil && strings.TrimSpace(providers) != "" {
			lines = append(lines, "", titleStyle.Render("Providers"), strings.TrimSpace(stripANSI(providers)))
		}
	} else {
		lines = append(lines, "", text.RunSetupHint)
	}
	return strings.Join(lines, "\n")
}

func setupContent(text uiText) string {
	lines := []string{
		titleStyle.Render(text.SetupLines[0]),
		text.SetupLines[1],
		"",
		text.SetupLines[2],
		"",
		text.SetupLines[3],
		text.SetupLines[4],
		text.SetupLines[5],
	}
	return strings.Join(lines, "\n")
}

func settingsContent(text uiText) string {
	config, err := configWithDefaults()
	if err != nil {
		return strings.Join([]string{titleStyle.Render(text.SettingsTitle), err.Error()}, "\n")
	}
	lines := []string{
		titleStyle.Render(text.SettingsTitle),
		fmt.Sprintf("Config:   %s", gachaConfigPath()),
		fmt.Sprintf("Model:    %s", configuredModelSummary(config.Model)),
		fmt.Sprintf("Language: %s", config.Language),
		fmt.Sprintf("Theme:    %s", configuredThemeSummary(config.Theme)),
		fmt.Sprintf("Active:   %s", detectLanguage()),
		"",
		sectionStyle.Render("Commands"),
		"/model auto",
		"/model opencode-default",
		"/model provider/model",
		"/language auto",
		"/language en",
		"/language ko",
		"/theme system",
		"/theme dark",
		"/theme light",
		"/theme gacha",
	}
	if envModel := strings.TrimSpace(os.Getenv("GACHA_OPENCODE_MODEL")); envModel != "" {
		lines = append(lines, "", mutedStyle.Render("GACHA_OPENCODE_MODEL is currently overriding the model setting."))
	}
	if envLang := strings.TrimSpace(os.Getenv("GACHA_LANG")); envLang != "" {
		lines = append(lines, mutedStyle.Render("GACHA_LANG is currently overriding the language setting."))
	}
	return strings.Join(lines, "\n")
}

func errorContent(err error, output string, text uiText) string {
	parts := []string{titleStyle.Render(text.ErrorTitle), err.Error()}
	if strings.TrimSpace(output) != "" {
		parts = append(parts, "", strings.TrimSpace(output))
	}
	return strings.Join(parts, "\n")
}
