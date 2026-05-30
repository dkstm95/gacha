package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type runResultMsg struct {
	query     string
	output    string
	completed bool
	err       error
}

type researchPhaseMsg struct{}

type tuiModel struct {
	version string
	lang    language
	text    uiText
	input   textinput.Model
	view    viewport.Model
	spin    spinner.Model
	width   int
	height  int
	busy    bool
	phase   int
	save    *pendingSave
	status  string
	runtime string
	mode    string
}

func newTUIModel(version string) tuiModel {
	lang := detectLanguage()
	text := textFor(lang)
	input := textinput.New()
	input.Placeholder = text.InputPlaceholder
	input.Focus()
	input.CharLimit = 2000
	input.Prompt = " "

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = accentStyle

	view := viewport.New(80, 16)
	view.SetContent(welcomeContent(version, text, 80, 16))

	return tuiModel{
		version: version,
		lang:    lang,
		text:    text,
		input:   input,
		view:    view,
		spin:    spin,
		width:   100,
		height:  34,
		status:  text.Ready,
		runtime: routeLabelFor(lang),
		mode:    text.Auto,
	}
}

func (m tuiModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.view.Width = max(32, msg.Width-4)
		m.view.Height = max(6, msg.Height-8)
		m.input.Width = max(16, msg.Width-8)
		if msg.Width < 72 {
			m.input.Placeholder = m.text.InputPlaceholderShort
		} else {
			m.input.Placeholder = m.text.InputPlaceholder
		}
		if m.mode == m.text.Auto && !m.busy {
			m.view.SetContent(welcomeContent(m.version, m.text, m.view.Width, m.view.Height))
		}
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			value := strings.TrimSpace(m.input.Value())
			if value == "" {
				return m, nil
			}
			m.input.SetValue("")
			if m.save != nil {
				return m.handleReportAction(value)
			}
			return m.handleSubmit(value)
		}
	case runResultMsg:
		m.busy = false
		m.phase = 0
		m.runtime = routeLabelFor(m.lang)
		m.mode = m.text.Report
		if msg.err != nil {
			m.status = m.text.Fallback
			m.view.SetContent(errorContent(msg.err, msg.output, m.text))
		} else {
			m.status = m.text.Complete
			if msg.completed && strings.TrimSpace(msg.output) != "" {
				m.save = &pendingSave{query: msg.query, report: strings.TrimSpace(msg.output)}
				m.view.SetContent(msg.output + "\n\n" + renderReportActions(m.text))
			} else {
				m.save = nil
				m.view.SetContent(msg.output)
			}
		}
		m.view.GotoTop()
	case researchPhaseMsg:
		if m.busy && len(m.text.ResearchPhases) > 0 {
			m.phase++
			m.status = m.text.ResearchPhases[m.phase%len(m.text.ResearchPhases)]
			cmds = append(cmds, researchPhaseTick())
		}
	}

	if m.busy {
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		cmds = append(cmds, cmd)
	}
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)
	var viewCmd tea.Cmd
	m.view, viewCmd = m.view.Update(msg)
	cmds = append(cmds, viewCmd)
	return m, tea.Batch(cmds...)
}

func (m tuiModel) handleSubmit(value string) (tea.Model, tea.Cmd) {
	m.save = nil
	switch value {
	case "/q", "/quit", "quit", "exit":
		return m, tea.Quit
	case "/h", "/help", "help":
		m.status = m.text.Help
		m.mode = m.text.Command
		m.view.SetContent(helpContent(m.text))
		m.view.GotoTop()
		return m, nil
	case "/settings", "settings":
		m.status = m.text.SettingsTitle
		m.mode = m.text.System
		m.view.SetContent(settingsContent(m.text))
		m.view.GotoTop()
		return m, nil
	case "/home", "home":
		m.status = m.text.Ready
		m.mode = m.text.Auto
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

	m.busy = true
	m.phase = 0
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

func (m tuiModel) View() string {
	width := max(44, m.width)
	outerPadding := 4
	if width < 72 {
		outerPadding = 2
	}
	bodyWidth := max(40, width-outerPadding)
	contentHeight := max(6, m.height-8)
	m.view.Width = max(30, bodyWidth-4)
	m.view.Height = contentHeight

	content := m.view.View()
	if m.mode == m.text.Auto && !m.busy {
		content = welcomeContent(m.version, m.text, m.view.Width, contentHeight)
		homeHeight := lipgloss.Height(content) + 2
		contentHeight = min(contentHeight, max(8, homeHeight))
		m.view.Height = contentHeight
	}

	header := renderHeader(bodyWidth, m.version)
	status := renderStatus(bodyWidth, m.status, m.runtime, m.mode, m.busy, m.spin.View(), m.text)
	panel := panelStyle.Width(bodyWidth - 2).Height(contentHeight).Render(content)
	input := renderInput(bodyWidth, m.input.View())
	footer := renderFooter(bodyWidth, m.text)
	return lipgloss.JoinVertical(lipgloss.Left, header, status, panel, input, footer)
}

func runPromptCmd(query string) tea.Cmd {
	return func() tea.Msg {
		result := runPrompt(query)
		return runResultMsg{query: query, output: result.output, completed: result.completed, err: result.err}
	}
}

func runDetailPromptCmd(query string, basicReport string) tea.Cmd {
	return func() tea.Msg {
		result := runDetailedPrompt(query, basicReport)
		return runResultMsg{query: query, output: result.output, completed: result.completed, err: result.err}
	}
}

func researchPhaseTick() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return researchPhaseMsg{}
	})
}

type promptRunResult struct {
	output    string
	completed bool
	err       error
}

type pendingSave struct {
	query  string
	report string
}

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
	compact := width < 72 || height < 14
	wide := width >= 104 && height >= 16
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
	lines := append([]string(nil), text.HelpLines...)
	lines[0] = titleStyle.Render(lines[0])
	return strings.Join(lines, "\n")
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
		fmt.Sprintf("Active:   %s", detectLanguage()),
		"",
		sectionStyle.Render("Commands"),
		"/model auto",
		"/model opencode-default",
		"/model provider/model",
		"/language auto",
		"/language en",
		"/language ko",
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

func renderReportActions(text uiText) string {
	var parts []string
	for _, action := range text.ReportActions {
		parts = append(parts, keyStyle.Render(action.Key)+" "+action.Label)
	}
	parts = append(parts, mutedStyle.Render(text.NewQuestionAction))
	return strings.Join([]string{
		sectionStyle.Render(text.ReportActionsTitle),
		strings.Join(parts, "   "),
	}, "\n")
}

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
	return statusStyle.Width(width - 2).Render(strings.Join(items, "   "))
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

func settingValue(value string) string {
	fields := strings.Fields(value)
	if len(fields) < 2 {
		return ""
	}
	return strings.TrimSpace(fields[1])
}

func configuredModelSummary(model string) string {
	if envModel := strings.TrimSpace(os.Getenv("GACHA_OPENCODE_MODEL")); envModel != "" {
		return envModel + " (from GACHA_OPENCODE_MODEL)"
	}
	switch {
	case model == "":
		return modelSettingAuto
	case strings.EqualFold(model, modelSettingAuto):
		return modelSettingAuto
	case strings.EqualFold(model, modelSettingOpenCodeDefault):
		return modelSettingOpenCodeDefault
	default:
		return model
	}
}

func wrapLine(value string, width int) string {
	value = strings.TrimSpace(value)
	if width <= 0 || lipgloss.Width(value) <= width {
		return value
	}
	words := strings.Fields(value)
	if len(words) == 0 {
		return value
	}
	var lines []string
	current := words[0]
	for _, word := range words[1:] {
		next := current + " " + word
		if lipgloss.Width(next) > width {
			lines = append(lines, current)
			current = word
			continue
		}
		current = next
	}
	lines = append(lines, current)
	return strings.Join(lines, "\n")
}

func wrapIndented(value string, width int, indent string) string {
	wrapped := wrapLine(value, width)
	lines := strings.Split(wrapped, "\n")
	for i := 1; i < len(lines); i++ {
		lines[i] = indent + lines[i]
	}
	return strings.Join(lines, "\n")
}

func padRight(value string, width int) string {
	padding := width - lipgloss.Width(value)
	if padding <= 0 {
		return value
	}
	return value + strings.Repeat(" ", padding)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func stripANSI(value string) string {
	var builder strings.Builder
	inEscape := false
	for _, r := range value {
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		if r == '\x1b' {
			inEscape = true
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

var (
	brandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Bold(true)
	accentStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	heroStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("87")).Bold(true)
	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Bold(true).
			Padding(0, 1)
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	faintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Bold(true)
	bulletStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	keyStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62")).Bold(true).Padding(0, 1)
	actionNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("238")).Bold(true).Padding(0, 1)
	noteStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("238")).PaddingLeft(1)
	statusStyle     = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(0, 1)
	calloutStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("214")).
			Padding(0, 1)
)
