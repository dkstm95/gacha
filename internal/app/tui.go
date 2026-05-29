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
	input := textinput.New()
	input.Placeholder = "Ask about an investment..."
	input.Focus()
	input.CharLimit = 2000
	input.Prompt = " "

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = accentStyle

	view := viewport.New(80, 16)
	view.SetContent(welcomeContent(version))

	return tuiModel{
		version: version,
		input:   input,
		view:    view,
		spin:    spin,
		width:   100,
		height:  34,
		status:  "Ready",
		runtime: routeLabel(),
		mode:    "Auto",
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
		m.runtime = routeLabel()
		m.mode = "Report"
		if msg.err != nil {
			m.status = "Fallback"
			m.view.SetContent(errorContent(msg.err, msg.output))
		} else {
			m.status = "Complete"
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
		m.status = "Help"
		m.mode = "Command"
		m.view.SetContent(helpContent())
		m.view.GotoTop()
		return m, nil
	case "/home", "home":
		m.status = "Ready"
		m.mode = "Auto"
		m.view.SetContent(welcomeContent(m.version))
		m.view.GotoTop()
		return m, nil
	case "/doctor", "doctor":
		m.status = "Doctor"
		m.runtime = routeLabel()
		m.mode = "Runtime"
		m.view.SetContent(doctorContent())
		m.view.GotoTop()
		return m, nil
	case "/setup", "setup":
		m.status = "Setup"
		m.mode = "Runtime"
		m.view.SetContent(setupContent())
		m.view.GotoTop()
		return m, nil
	case "/update", "update":
		m.status = "Update"
		m.mode = "System"
		m.view.SetContent("Run `gacha update` outside the interactive UI to update the binary.")
		m.view.GotoTop()
		return m, nil
	}

	m.busy = true
	m.status = "Researching"
	m.mode = "Auto"
	m.view.SetContent(researchingContent(value))
	return m, tea.Batch(m.spin.Tick, runPromptCmd(value))
}

func (m tuiModel) View() string {
	width := max(72, m.width)
	bodyWidth := max(64, width-4)
	contentHeight := max(8, m.height-17)
	m.view.Width = bodyWidth - 4
	m.view.Height = contentHeight

	header := renderHeader(bodyWidth, m.version)
	status := renderStatus(bodyWidth, m.status, m.runtime, m.mode, m.busy, m.spin.View())
	panel := panelStyle.Width(bodyWidth - 2).Height(contentHeight).Render(m.view.View())
	input := renderInput(bodyWidth, m.input.View())
	footer := mutedStyle.Render(" /help  /doctor  /setup  /update  /quit   •   enter to run   •   esc to exit")
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

func welcomeContent(version string) string {
	return strings.Join([]string{
		titleStyle.Render("Investment research cockpit"),
		"Ask one question. Gacha routes it through the right research workflow.",
		"Every workflow requires current web or market data before analysis.",
		"",
		titleStyle.Render("Workflow rail"),
		chipStyle.Render("Discover") + "  find opportunities when you do not know what to buy",
		chipStyle.Render("Select") + "    rank concrete assets inside a sector or theme",
		chipStyle.Render("Entry") + "     decide whether the current price is attractive",
		chipStyle.Render("Exit") + "      define trim, sell, stop-loss, and thesis-break zones",
		chipStyle.Render("Portfolio") + " review concentration, exposure, and rebalancing risks",
		chipStyle.Render("Journal") + "   record thesis, decision rules, and postmortems",
		"",
		titleStyle.Render("Try"),
		"NVDA 지금 사도 될까?",
		"What should I invest in for the next 6 to 12 months?",
		"I own TSLA. When should I trim, sell, or stop out?",
		"",
		titleStyle.Render("Report contract"),
		"Data freshness • Sources • Thesis • Valuation • Risks • Devil's Advocate • Action conditions",
		"",
		mutedStyle.Render("No fresh data, no recommendation. Trading is disabled. Version " + version),
	}, "\n")
}

func researchingContent(query string) string {
	return strings.Join([]string{
		titleStyle.Render("Research run"),
		"Query:",
		"  " + query,
		"",
		titleStyle.Render("Pipeline"),
		"1. Classify request: discover, select, entry, exit, portfolio, or journal",
		"2. Require current web or market data",
		"3. Build thesis, valuation, and scenario analysis",
		"4. Run risk review and Devil's Advocate",
		"5. Produce action conditions and provenance",
		"",
		mutedStyle.Render("Waiting for the local AI runtime..."),
	}, "\n")
}

func helpContent() string {
	return strings.Join([]string{
		titleStyle.Render("Command palette"),
		"/home     return to the dashboard",
		"/help     show this command palette",
		"/doctor   inspect OpenCode runtime and provider auth",
		"/setup    show setup instructions",
		"/update   show update instructions",
		"/quit     exit",
	}, "\n")
}

func doctorContent() string {
	status := "missing"
	if hasRunnableCommand(openCodeCommand) {
		status = "ready"
		if !hasOpenCodeAuth() {
			status = "login required"
		}
	}
	lines := []string{
		titleStyle.Render("Runtime"),
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
		lines = append(lines, "", "Run `gch setup` outside this screen to connect ChatGPT, Copilot, Gemini, or an API provider.")
	}
	return strings.Join(lines, "\n")
}

func setupContent() string {
	return strings.Join([]string{
		titleStyle.Render("Setup"),
		"Run this command in your shell:",
		"",
		"  gch setup",
		"",
		"That flow installs OpenCode if needed and starts provider login.",
		"Interactive provider login is intentionally handled outside this screen so your terminal can hand control to OpenCode safely.",
	}, "\n")
}

func errorContent(err error, output string) string {
	parts := []string{titleStyle.Render("OpenCode failed"), err.Error()}
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

func renderStatus(width int, status string, runtime string, mode string, busy bool, spin string) string {
	indicator := "●"
	if busy {
		indicator = spin
	}
	items := []string{
		accentStyle.Render(indicator + " " + status),
		"Mode " + mode,
		"Runtime " + runtime,
		"Fresh data required",
		"No trading",
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
