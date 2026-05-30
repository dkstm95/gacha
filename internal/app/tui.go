package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type runResultMsg struct {
	output string
	err    error
}

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
	view.SetContent(welcomeContent(version, text))

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
		m.view.Width = max(40, msg.Width-4)
		m.view.Height = max(8, msg.Height-17)
		m.input.Width = max(20, msg.Width-8)
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
			return m.handleSubmit(value)
		}
	case runResultMsg:
		m.busy = false
		m.runtime = routeLabelFor(m.lang)
		m.mode = m.text.Report
		if msg.err != nil {
			m.status = m.text.Fallback
			m.view.SetContent(errorContent(msg.err, msg.output, m.text))
		} else {
			m.status = m.text.Complete
			m.view.SetContent(msg.output)
		}
		m.view.GotoTop()
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
	switch value {
	case "/q", "/quit", "quit", "exit":
		return m, tea.Quit
	case "/h", "/help", "help":
		m.status = m.text.Help
		m.mode = m.text.Command
		m.view.SetContent(helpContent(m.text))
		m.view.GotoTop()
		return m, nil
	case "/home", "home":
		m.status = m.text.Ready
		m.mode = m.text.Auto
		m.view.SetContent(welcomeContent(m.version, m.text))
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

	m.busy = true
	m.status = m.text.Researching
	m.mode = m.text.Auto
	m.view.SetContent(researchingContent(value, m.text))
	return m, tea.Batch(m.spin.Tick, runPromptCmd(value))
}

func (m tuiModel) View() string {
	width := max(72, m.width)
	bodyWidth := max(64, width-4)
	contentHeight := max(8, m.height-17)
	m.view.Width = bodyWidth - 4
	m.view.Height = contentHeight

	header := renderHeader(bodyWidth, m.version)
	status := renderStatus(bodyWidth, m.status, m.runtime, m.mode, m.busy, m.spin.View(), m.text)
	panel := panelStyle.Width(bodyWidth - 2).Height(contentHeight).Render(m.view.View())
	input := renderInput(bodyWidth, m.input.View())
	footer := mutedStyle.Render(m.text.Footer)
	return lipgloss.JoinVertical(lipgloss.Left, header, status, panel, input, footer)
}

func runPromptCmd(query string) tea.Cmd {
	return func() tea.Msg {
		output, err := runPrompt(query)
		return runResultMsg{output: output, err: err}
	}
}

func runPrompt(query string) (string, error) {
	prompt, err := buildPrompt([]string{query})
	if err != nil {
		return "", err
	}
	commandPath, ok := resolveCommand(openCodeCommand)
	if !ok || !hasOpenCodeAuth() {
		return prompt, nil
	}
	output, err := runOpenCodeWithResolution(commandPath, prompt, resolveOpenCodeModel(), false)
	if err != nil {
		return output, err
	}
	return strings.TrimSpace(output), nil
}

func welcomeContent(version string, text uiText) string {
	lines := []string{
		titleStyle.Render(text.Welcome[0]),
		text.Welcome[1],
		text.Welcome[2],
		"",
		titleStyle.Render(text.Welcome[3]),
	}
	for _, item := range text.Welcome[4:10] {
		parts := strings.SplitN(item, "|", 2)
		lines = append(lines, chipStyle.Render(parts[0])+"  "+parts[1])
	}
	lines = append(lines,
		"",
		titleStyle.Render(text.Welcome[10]),
		text.Welcome[11],
		text.Welcome[12],
		text.Welcome[13],
		"",
		titleStyle.Render(text.Welcome[14]),
		text.Welcome[15],
		"",
		mutedStyle.Render(text.Welcome[16]+" Version "+version),
	)
	return strings.Join(lines, "\n")
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
	return strings.Join([]string{
		titleStyle.Render(text.SetupLines[0]),
		text.SetupLines[1],
		"",
		text.SetupLines[2],
		"",
		text.SetupLines[3],
		text.SetupLines[4],
	}, "\n")
}

func errorContent(err error, output string, text uiText) string {
	parts := []string{titleStyle.Render(text.ErrorTitle), err.Error()}
	if strings.TrimSpace(output) != "" {
		parts = append(parts, "", strings.TrimSpace(output))
	}
	return strings.Join(parts, "\n")
}

func renderHeader(width int, version string) string {
	left := brandStyle.Render(" GACHA ")
	right := mutedStyle.Render("v" + version)
	line := strings.Repeat("─", max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-2))
	return lipgloss.JoinHorizontal(lipgloss.Center, left, mutedStyle.Render(line), right)
}

func renderStatus(width int, status string, runtime string, mode string, busy bool, spin string, text uiText) string {
	indicator := "●"
	if busy {
		indicator = spin
	}
	items := []string{
		accentStyle.Render(indicator + " " + status),
		text.StatusMode + mode,
		text.StatusRuntime + runtime,
		text.StatusFreshData,
		text.StatusNoTrading,
	}
	return statusStyle.Width(width - 2).Render(strings.Join(items, "   "))
}

func renderInput(width int, input string) string {
	return inputStyle.Width(width - 2).Render(input)
}

func max(a, b int) int {
	if a > b {
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
	accentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	chipStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("238")).Padding(0, 1)
	statusStyle = lipgloss.NewStyle().
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
)
