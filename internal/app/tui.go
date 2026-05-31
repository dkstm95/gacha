package app

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
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
	query   string
	report  string
	choice  *pendingChoice
}

func newTUIModel(version string) tuiModel {
	setThemeStyles(configuredTheme())
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
		m.input.Placeholder = inputPlaceholderForWidth(m.text, msg.Width)
		if m.choice != nil {
			m.view.SetContent(m.choice.RenderWidth(m.text, m.view.Width))
		} else if m.mode == m.text.Auto && !m.busy {
			m.view.SetContent(welcomeContent(m.version, m.text, m.view.Width, m.view.Height))
		}
	case tea.KeyMsg:
		key := msg.String()
		if m.choice != nil {
			switch key {
			case "up", "k":
				m.moveChoice(-1)
				return m, nil
			case "down", "j":
				m.moveChoice(1)
				return m, nil
			}
		}
		switch key {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			value := strings.TrimSpace(m.input.Value())
			if value == "" && m.choice != nil {
				return m.handleChoiceSelection()
			}
			if value == "" {
				return m, nil
			}
			m.input.SetValue("")
			if m.choice != nil {
				m.choice = nil
			}
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
		m.query = msg.query
		if msg.err != nil {
			m.status = m.text.Fallback
			m.report = ""
			m.view.SetContent(errorContent(msg.err, msg.output, m.text))
		} else {
			m.status = m.text.Complete
			if msg.completed && strings.TrimSpace(msg.output) != "" {
				report := strings.TrimSpace(msg.output)
				m.report = report
				m.save = &pendingSave{query: msg.query, report: report}
				m = m.showReportChoice()
			} else {
				m.save = nil
				m.report = strings.TrimSpace(msg.output)
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

func (m tuiModel) View() string {
	width := max(44, m.width)
	outerPadding := 4
	if width < 72 {
		outerPadding = 2
	}
	bodyWidth := max(40, width-outerPadding)
	contentHeight := max(6, m.height-7)
	fullLayout := bodyWidth >= 132 && m.height >= 22
	_, workspaceWidth := splitLayoutWidths(bodyWidth)
	if fullLayout {
		m.view.Width = max(30, workspaceWidth-8)
		m.input.Width = max(16, workspaceWidth-12)
	} else {
		m.view.Width = max(30, bodyWidth-4)
		m.input.Width = max(16, bodyWidth-8)
	}
	m.input.Placeholder = inputPlaceholderForWidth(m.text, m.input.Width)
	workspaceHeight := contentHeight
	if fullLayout {
		workspaceHeight = max(6, contentHeight-4)
	}
	m.view.Height = workspaceHeight

	content := m.view.View()
	if m.mode == m.text.Auto && !m.busy {
		content = welcomeContentWithColumns(m.version, m.text, m.view.Width, workspaceHeight, !fullLayout)
		if !fullLayout {
			homeHeight := lipgloss.Height(content) + 2
			contentHeight = min(contentHeight, max(8, homeHeight))
			workspaceHeight = contentHeight
		}
		m.view.Height = workspaceHeight
	}

	header := renderHeader(bodyWidth, m.version)
	panel := panelStyle.Width(bodyWidth - 2).Height(contentHeight).Render(content)
	if fullLayout {
		panel = m.renderSplitMain(bodyWidth, contentHeight, content, m.input.View())
	}
	status := renderStatus(bodyWidth, m.status, m.runtime, m.mode, m.busy, m.spin.View(), m.text)
	parts := []string{header}
	parts = append(parts, panel)
	if !fullLayout {
		parts = append(parts, renderInput(bodyWidth, m.input.View()))
	}
	parts = append(parts, status)
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m tuiModel) renderSplitMain(bodyWidth int, height int, workspace string, input string) string {
	gap := 2
	railWidth, workspaceWidth := splitLayoutWidths(bodyWidth)
	rail := panelStyle.Width(railWidth - 2).Height(height).Render(m.contextRail(railWidth - 4))
	inputBlock := renderInput(workspaceWidth-8, input)
	innerHeight := max(1, height-4)
	spacerHeight := max(1, innerHeight-lipgloss.Height(workspace)-lipgloss.Height(inputBlock))
	mainContent := lipgloss.JoinVertical(lipgloss.Left, workspace, strings.Repeat("\n", spacerHeight-1), inputBlock)
	main := panelStyle.Width(workspaceWidth - 2).Height(height).Render(mainContent)
	return lipgloss.JoinHorizontal(lipgloss.Top, rail, strings.Repeat(" ", gap), main)
}

func splitLayoutWidths(bodyWidth int) (int, int) {
	gap := 2
	railWidth := max(34, bodyWidth/3)
	workspaceWidth := max(58, bodyWidth-railWidth-gap)
	if railWidth+workspaceWidth+gap > bodyWidth {
		workspaceWidth = max(40, bodyWidth-railWidth-gap)
	}
	return railWidth, workspaceWidth
}

func inputPlaceholderForWidth(text uiText, width int) string {
	if width < 72 {
		return text.InputPlaceholderShort
	}
	return text.InputPlaceholder
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
