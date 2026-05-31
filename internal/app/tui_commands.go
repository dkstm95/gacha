package app

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

func (m tuiModel) handleSubmit(value string) (tea.Model, tea.Cmd) {
	m.save = nil
	m.choice = nil
	if isSettingsCommand(value) {
		m.status = m.text.SettingsTitle
		m.mode = m.text.System
		m.view.SetContent(settingsContent(m.text))
		m.view.GotoTop()
		return m, nil
	}
	switch value {
	case "/q", "/quit", "quit", "exit":
		return m, tea.Quit
	case "/h", "/help", "help":
		m.status = m.text.Help
		m.mode = m.text.Command
		m.view.SetContent(helpContent(m.text))
		m.view.GotoTop()
		return m, nil
	case "/theme", "theme", "/themes", "themes":
		return m.showThemeChoice()
	case "/model", "model":
		return m.showModelChoice()
	case "/language", "language", "/lang", "lang":
		return m.showLanguageChoice()
	case "/home", "home":
		m.status = m.text.Ready
		m.mode = m.text.Auto
		m.query = ""
		m.report = ""
		m.view.SetContent(welcomeContent(m.version, m.text, m.view.Width, m.view.Height))
		m.view.GotoTop()
		return m, nil
	case "/doctor", "doctor":
		m.status = "Doctor"
		m.runtime = routeLabelFor(m.lang)
		m.mode = m.text.Runtime
		m.view.SetContent(doctorContent(m.text))
		m.view.GotoTop()
		return m, nil
	case "/setup", "setup":
		m.status = m.text.Setup
		m.mode = m.text.Runtime
		m.view.SetContent(setupContent(m.text))
		m.view.GotoTop()
		return m, nil
	case "/update", "update":
		m.status = m.text.Update
		m.mode = m.text.System
		m.view.SetContent(m.text.UpdateMessage)
		m.view.GotoTop()
		return m, nil
	}
	if strings.HasPrefix(value, "/model ") || strings.HasPrefix(value, "model ") {
		return m.handleModelSetting(value)
	}
	if strings.HasPrefix(value, "/language ") || strings.HasPrefix(value, "language ") || strings.HasPrefix(value, "/lang ") || strings.HasPrefix(value, "lang ") {
		return m.handleLanguageSetting(value)
	}
	if strings.HasPrefix(value, "/theme ") || strings.HasPrefix(value, "theme ") {
		return m.handleThemeSetting(value)
	}
	if isSlashCommand(value) {
		m.status = m.text.Help
		m.mode = m.text.Command
		m.view.SetContent(unknownCommandContent(value, m.text))
		m.view.GotoTop()
		return m, nil
	}

	m.busy = true
	m.phase = 0
	m.query = value
	m.report = ""
	if len(m.text.ResearchPhases) > 0 {
		m.status = m.text.ResearchPhases[0]
	} else {
		m.status = m.text.Researching
	}
	m.mode = m.text.Auto
	m.view.SetContent(researchingContent(value, m.text))
	return m, tea.Batch(m.spin.Tick, researchPhaseTick(), runPromptCmd(value))
}

func (m tuiModel) handleModelSetting(value string) (tea.Model, tea.Cmd) {
	model := settingValue(value)
	if err := updateConfigModel(model); err != nil {
		if !validModelSetting(model) {
			return m.showSettingsError(m.text.SettingsInvalidModel)
		}
		return m.showError(err)
	}
	return m.showSettingsSaved()
}

func (m tuiModel) handleLanguageSetting(value string) (tea.Model, tea.Cmd) {
	if err := updateConfigLanguage(settingValue(value)); err != nil {
		if _, ok := normalizeLanguageSetting(settingValue(value)); !ok {
			return m.showSettingsError(m.text.SettingsInvalidLang)
		}
		return m.showError(err)
	}
	m.lang = detectLanguage()
	m.text = textFor(m.lang)
	m.input.Placeholder = m.text.InputPlaceholder
	return m.showSettingsSaved()
}

func (m tuiModel) handleThemeSetting(value string) (tea.Model, tea.Cmd) {
	theme := settingValue(value)
	normalized, ok := normalizeThemeSetting(theme)
	if !ok {
		return m.showSettingsError(m.text.SettingsInvalidTheme)
	}
	if err := updateConfigTheme(normalized); err != nil {
		return m.showError(err)
	}
	setThemeStyles(normalized)
	m.spin.Style = accentStyle
	options := themeChoiceOptions(m.text)
	m.choice = &pendingChoice{
		Kind:     choiceTheme,
		Title:    m.text.ThemeTitle,
		Intro:    fmt.Sprintf("%s %s", m.text.ThemeActive, themeLabel(themeByName(normalized), m.text)),
		Options:  options,
		Selected: selectedChoiceIndex(options, normalized),
	}
	m.status = m.text.SettingsSaved
	m.mode = m.text.System
	m.view.SetContent(m.text.SettingsSaved + "\n\n" + m.choice.Render(m.text))
	m.view.GotoTop()
	return m, nil
}

func (m *tuiModel) moveChoice(delta int) {
	if m.choice == nil || len(m.choice.Options) == 0 {
		return
	}
	next := (m.choice.Selected + delta) % len(m.choice.Options)
	if next < 0 {
		next += len(m.choice.Options)
	}
	m.choice.Selected = next
	m.view.SetContent(m.choice.Render(m.text))
	m.view.GotoTop()
}

func (m tuiModel) handleChoiceSelection() (tea.Model, tea.Cmd) {
	if m.choice == nil || len(m.choice.Options) == 0 {
		return m, nil
	}
	choice := m.choice
	selected := choice.Options[choice.Selected]
	m.choice = nil
	switch choice.Kind {
	case choiceTheme:
		return m.handleThemeSetting("/theme " + selected.Value)
	case choiceLanguage:
		return m.handleLanguageSetting("/language " + selected.Value)
	case choiceModel:
		return m.handleModelSetting("/model " + selected.Value)
	case choiceReport:
		m.input.SetValue("")
		return m.handleReportAction(selected.Value)
	default:
		return m, nil
	}
}

func (m tuiModel) showModelChoice() (tea.Model, tea.Cmd) {
	config, _ := configWithDefaults()
	options := []choiceOption{
		{Label: "Auto", Value: modelSettingAuto, Description: m.text.ModelDescriptions[modelSettingAuto]},
		{Label: "OpenCode default", Value: modelSettingOpenCodeDefault, Description: m.text.ModelDescriptions[modelSettingOpenCodeDefault]},
	}
	selected := selectedChoiceIndex(options, configuredModelSummary(config.Model))
	m.choice = &pendingChoice{
		Kind:     choiceModel,
		Title:    m.text.ModelTitle,
		Intro:    m.text.ModelIntro,
		Options:  options,
		Selected: selected,
		Footer:   m.text.ModelCustomHint,
	}
	m.status = m.text.ModelTitle
	m.mode = m.text.System
	m.view.SetContent(m.choice.Render(m.text))
	m.view.GotoTop()
	return m, nil
}

func (m tuiModel) showLanguageChoice() (tea.Model, tea.Cmd) {
	config, _ := configWithDefaults()
	active := config.Language
	if active == "" {
		active = languageSettingAuto
	}
	options := []choiceOption{
		{Label: "Auto", Value: languageSettingAuto, Description: m.text.LanguageDescriptions[languageSettingAuto]},
		{Label: "English", Value: languageSettingEnglish, Description: m.text.LanguageDescriptions[languageSettingEnglish]},
		{Label: "한국어", Value: languageSettingKorean, Description: m.text.LanguageDescriptions[languageSettingKorean]},
	}
	m.choice = &pendingChoice{
		Kind:     choiceLanguage,
		Title:    m.text.LanguageTitle,
		Intro:    m.text.LanguageIntro,
		Options:  options,
		Selected: selectedChoiceIndex(options, active),
	}
	m.status = m.text.LanguageTitle
	m.mode = m.text.System
	m.view.SetContent(m.choice.Render(m.text))
	m.view.GotoTop()
	return m, nil
}

func (m tuiModel) showThemeChoice() (tea.Model, tea.Cmd) {
	active := configuredTheme()
	options := themeChoiceOptions(m.text)
	m.choice = &pendingChoice{
		Kind:     choiceTheme,
		Title:    m.text.ThemeTitle,
		Intro:    m.text.ThemeIntro,
		Options:  options,
		Selected: selectedChoiceIndex(options, active),
	}
	m.status = m.text.ThemeTitle
	m.mode = m.text.System
	m.view.SetContent(m.choice.Render(m.text))
	m.view.GotoTop()
	return m, nil
}

func (m tuiModel) showSettingsSaved() (tea.Model, tea.Cmd) {
	m.status = m.text.SettingsSaved
	m.mode = m.text.System
	m.view.SetContent(m.text.SettingsSaved + "\n\n" + settingsContent(m.text))
	m.view.GotoTop()
	return m, nil
}

func (m tuiModel) showSettingsError(message string) (tea.Model, tea.Cmd) {
	m.status = m.text.SettingsTitle
	m.mode = m.text.System
	m.view.SetContent(message + "\n\n" + settingsContent(m.text))
	m.view.GotoTop()
	return m, nil
}

func (m tuiModel) showError(err error) (tea.Model, tea.Cmd) {
	m.status = m.text.Fallback
	m.view.SetContent(errorContent(err, "", m.text))
	m.view.GotoTop()
	return m, nil
}

func (m tuiModel) handleReportAction(value string) (tea.Model, tea.Cmd) {
	pending := m.save
	m.save = nil
	if pending == nil {
		return m, nil
	}
	if wantsSaveReport(value) {
		path, err := saveReport(pending.query, pending.report)
		if err != nil {
			m.status = m.text.Fallback
			m.view.SetContent(pending.report + "\n\n---\n" + errorContent(err, "", m.text))
			m.view.GotoBottom()
			return m, nil
		}
		m.status = m.text.Complete
		m.view.SetContent(pending.report + "\n\n---\n" + m.text.SavedReport + " " + path)
		m.view.GotoBottom()
		return m, nil
	}
	if refusesSaveReport(value) {
		m.status = m.text.Complete
		m.view.SetContent(pending.report + "\n\n---\n" + m.text.SkippedSave)
		m.view.GotoBottom()
		return m, nil
	}
	if wantsDetailedAnalysis(value) {
		m.busy = true
		m.phase = 0
		if len(m.text.ResearchPhases) > 0 {
			m.status = m.text.ResearchPhases[0]
		} else {
			m.status = m.text.Researching
		}
		m.mode = m.text.Report
		m.view.SetContent(researchingContent(pending.query, m.text))
		m.view.GotoTop()
		return m, tea.Batch(m.spin.Tick, researchPhaseTick(), runDetailPromptCmd(pending.query, pending.report))
	}
	return m.handleSubmit(value)
}

func (m tuiModel) showReportChoice() tuiModel {
	m.choice = &pendingChoice{
		Kind:    choiceReport,
		Title:   m.text.ReportActionsTitle,
		Intro:   m.text.ReportChoiceIntro,
		Options: reportChoiceOptions(m.text),
		Footer:  m.text.NewQuestionAction,
	}
	m.view.SetContent(m.report + "\n\n" + m.choice.Render(m.text))
	m.view.GotoBottom()
	return m
}
