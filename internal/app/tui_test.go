package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"strings"
	"testing"
)

func keyMsg(value string) tea.KeyPressMsg {
	switch value {
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	case "up":
		return tea.KeyPressMsg{Code: tea.KeyUp}
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	}
	code := rune(0)
	for _, r := range value {
		code = r
		break
	}
	return tea.KeyPressMsg{Code: code, Text: value}
}

func TestTUILanguageAndModelCommandsOpenChoices(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	for _, tc := range []struct {
		value string
		kind  choiceKind
		title string
	}{
		{value: "/language", kind: choiceLanguage, title: "Language"},
		{value: "/model", kind: choiceModel, title: "Model"},
	} {
		t.Run(tc.value, func(t *testing.T) {
			model := newTUIModel("0.1.27")
			next, cmd := model.handleSubmit(tc.value)
			if cmd != nil {
				t.Fatal("unexpected command")
			}
			updated := next.(tuiModel)
			if updated.choice == nil || updated.choice.Kind != tc.kind {
				t.Fatalf("expected %s choice, got %#v", tc.kind, updated.choice)
			}
			if !strings.Contains(stripANSI(updated.view.View()), tc.title) {
				t.Fatalf("choice view missing title %q:\n%s", tc.title, stripANSI(updated.view.View()))
			}
		})
	}
}

func TestTUIUnknownSlashCommandDoesNotRunPrompt(t *testing.T) {
	model := newTUIModel("0.1.27")
	next, cmd := model.handleSubmit("/setting")
	if cmd != nil {
		t.Fatal("unknown slash command should not run a prompt command")
	}
	updated := next.(tuiModel)
	got := stripANSI(updated.view.View())
	for _, expected := range []string{"Unknown command: /setting", "Command palette"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("unknown command view missing %q:\n%s", expected, got)
		}
	}
}

func TestTUIShowsPromptOutputWhileBusy(t *testing.T) {
	model := newTUIModel("0.1.27")
	model.busy = true
	model.query = "NVDA"
	model.view.SetContent(researchingContent(model.query, model.text))

	next, cmd := model.Update(promptOutputMsg{query: "NVDA", chunk: "\x1b[32mpartial report\rnext line"})
	if cmd != nil {
		t.Fatal("unexpected command without stream")
	}
	updated := next.(tuiModel)
	got := stripANSI(updated.view.View())
	if !strings.Contains(got, "partial report") || !strings.Contains(got, "next line") {
		t.Fatalf("streamed output not visible:\n%s", got)
	}
	if updated.mode != updated.text.Report {
		t.Fatalf("expected report mode while output is streaming, got %q", updated.mode)
	}
}

func TestOnboardingContentReflectsSetupState(t *testing.T) {
	text := englishText()
	if got := onboardingContent(text, 80, setupReady); got != "" {
		t.Fatalf("expected no onboarding for ready state, got %q", got)
	}
	if got := onboardingContent(text, 80, setupRuntimeMissing); !strings.Contains(stripANSI(got), "OpenCode is not installed yet") {
		t.Fatalf("unexpected runtime onboarding: %q", got)
	}
	if got := onboardingContent(text, 80, setupProviderMissing); !strings.Contains(stripANSI(got), "no AI provider") {
		t.Fatalf("unexpected provider onboarding: %q", got)
	}
}

func TestWelcomeContentIsDecisionDesk(t *testing.T) {
	got := stripANSI(welcomeContent("0.1.27", englishText(), 80, 16))
	for _, expected := range []string{
		"What are you deciding?",
		"Decision desk",
		"Buy",
		"Find",
		"Hold",
		"Exit",
		"Portfolio",
		"You'll get",
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("welcome content missing %q:\n%s", expected, got)
		}
	}
	for _, unwanted := range []string{
		"Ask -> Fresh data",
		"[Fresh data required]",
		"[No trading]",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("welcome content kept decorative element %q:\n%s", unwanted, got)
		}
	}
}

func TestWelcomeContentFitsCommonTerminalSizes(t *testing.T) {
	for _, tc := range []struct {
		name   string
		width  int
		height int
		text   uiText
	}{
		{name: "quarter english", width: 80, height: 16, text: englishText()},
		{name: "half english", width: 170, height: 24, text: englishText()},
		{name: "quarter korean", width: 80, height: 16, text: koreanText()},
		{name: "half korean", width: 170, height: 24, text: koreanText()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := stripANSI(welcomeContent("0.1.27", tc.text, tc.width, tc.height))
			for _, line := range strings.Split(got, "\n") {
				if lipgloss.Width(line) > tc.width {
					t.Fatalf("line width %d exceeds %d: %q\n%s", lipgloss.Width(line), tc.width, line, got)
				}
			}
		})
	}
}

func TestTUIViewFitsCommonTerminalSizes(t *testing.T) {
	for _, tc := range []struct {
		name   string
		width  int
		height int
	}{
		{name: "quarter", width: 80, height: 24},
		{name: "half", width: 120, height: 30},
		{name: "full", width: 180, height: 36},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model := newTUIModel("0.1.27")
			model.width = tc.width
			model.height = tc.height
			model.view.SetWidth(max(30, tc.width-8))
			model.view.SetHeight(max(6, tc.height-8))
			got := stripANSI(model.render())
			for _, line := range strings.Split(got, "\n") {
				if lipgloss.Width(line) > tc.width {
					t.Fatalf("line width %d exceeds %d: %q\n%s", lipgloss.Width(line), tc.width, line, got)
				}
			}
			if strings.Contains(got, "Checks fresh data") {
				t.Fatalf("status bar repeated safety copy:\n%s", got)
			}
			if tc.name == "full" {
				for _, expected := range []string{"Context", "Decision types", "Decision desk", "Ready"} {
					if !strings.Contains(got, expected) {
						t.Fatalf("full layout missing %q:\n%s", expected, got)
					}
				}
				if strings.Contains(got, "No saved reports yet") {
					t.Fatalf("home context showed empty recent state:\n%s", got)
				}
				railWidth, _ := splitLayoutWidths(max(40, tc.width-4))
				foundPrompt := false
				for _, line := range strings.Split(got, "\n") {
					promptColumn := strings.Index(line, "Ask:")
					if promptColumn < 0 {
						continue
					}
					foundPrompt = true
					if promptColumn < railWidth {
						t.Fatalf("full layout prompt crossed into the context rail at column %d before rail width %d:\n%s", promptColumn, railWidth, got)
					}
				}
				if !foundPrompt {
					t.Fatalf("full layout missing prompt input:\n%s", got)
				}
			}
		})
	}
}

func TestTUICommandViewsFitQuarterTerminal(t *testing.T) {
	for _, command := range []string{"/theme", "/help"} {
		t.Run(command, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			model := newTUIModel("0.1.27")
			next, cmd := model.handleSubmit(command)
			if cmd != nil {
				t.Fatal("unexpected command")
			}
			updated := next.(tuiModel)
			updated.width = 80
			updated.height = 24
			got := stripANSI(updated.render())
			for _, line := range strings.Split(got, "\n") {
				if lipgloss.Width(line) > 80 {
					t.Fatalf("line width %d exceeds 80: %q\n%s", lipgloss.Width(line), line, got)
				}
			}
			for _, fragment := range []string{
				"│  pro                                                                     │",
				"│  /languag                                                                │",
			} {
				if strings.Contains(got, fragment) {
					t.Fatalf("command view contains wrapped fragment %q:\n%s", fragment, got)
				}
			}
		})
	}
}

func TestTUIChoiceViewsRerenderAfterNarrowResize(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	model := newTUIModel("0.1.27")
	next, cmd := model.handleSubmit("/theme")
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	resized, cmd := next.(tuiModel).Update(tea.WindowSizeMsg{Width: 60, Height: 20})
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	got := stripANSI(resized.(tuiModel).render())
	for _, line := range strings.Split(got, "\n") {
		if lipgloss.Width(line) > 60 {
			t.Fatalf("line width %d exceeds 60: %q\n%s", lipgloss.Width(line), line, got)
		}
	}
	for _, fragment := range []string{
		"│  backgroun",
		"│  lig",
		"│  Gach",
		"──                                                    │",
	} {
		if strings.Contains(got, fragment) {
			t.Fatalf("narrow theme view contains clipped fragment %q:\n%s", fragment, got)
		}
	}
}

func TestTUIFullLayoutKeepsLocalizedPromptInRightPanel(t *testing.T) {
	for _, tc := range []struct {
		name       string
		text       uiText
		promptText string
	}{
		{name: "english", text: englishText(), promptText: "Ask:"},
		{name: "korean", text: koreanText(), promptText: "예:"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model := newTUIModel("0.1.27")
			model.text = tc.text
			model.input.Placeholder = tc.text.InputPlaceholder
			model.width = 180
			model.height = 36
			got := stripANSI(model.render())
			railWidth, _ := splitLayoutWidths(max(40, model.width-4))
			foundPrompt := false
			for _, line := range strings.Split(got, "\n") {
				promptColumn := strings.Index(line, tc.promptText)
				if promptColumn < 0 {
					continue
				}
				foundPrompt = true
				if promptColumn < railWidth {
					t.Fatalf("localized prompt crossed into the context rail at column %d before rail width %d:\n%s", promptColumn, railWidth, got)
				}
			}
			if !foundPrompt {
				t.Fatalf("missing localized prompt %q:\n%s", tc.promptText, got)
			}
		})
	}
}

func TestTUIStatusBarRendersBelowWorkspace(t *testing.T) {
	model := newTUIModel("0.1.27")
	model.width = 180
	model.height = 36
	got := stripANSI(model.render())
	promptLine := -1
	statusLine := -1
	panelBottomLine := -1
	for i, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "Ask:") {
			promptLine = i
		}
		if strings.Contains(line, "Ready") && strings.Contains(line, "Mode") {
			statusLine = i
		}
		if strings.Contains(line, "╰") && statusLine < 0 {
			panelBottomLine = i
		}
	}
	if promptLine < 0 {
		t.Fatalf("missing prompt input:\n%s", got)
	}
	if statusLine <= promptLine {
		t.Fatalf("status bar did not render below workspace; prompt line %d status line %d:\n%s", promptLine, statusLine, got)
	}
	if panelBottomLine < 0 {
		t.Fatalf("missing split panel bottom border:\n%s", got)
	}
	if panelBottomLine-promptLine > 5 {
		t.Fatalf("prompt input is too far above right panel bottom; prompt line %d panel bottom line %d:\n%s", promptLine, panelBottomLine, got)
	}
}

func TestTUIViewPlacesRealCursorOnInputLine(t *testing.T) {
	model := newTUIModel("0.1.27")
	model.width = 100
	model.height = 30
	model.input.SetValue("한")

	view := model.View()
	if view.Cursor == nil {
		t.Fatal("expected real input cursor")
	}
	if model.input.VirtualCursor() {
		t.Fatal("text input should use Bubble Tea's real cursor path")
	}

	lines := strings.Split(stripANSI(view.Content), "\n")
	inputLine := -1
	inputColumn := -1
	for i, line := range lines {
		if column := strings.Index(line, "한"); column >= 0 {
			inputLine = i
			inputColumn = column
			break
		}
	}
	if inputLine < 0 {
		t.Fatalf("rendered view missing input text:\n%s", stripANSI(view.Content))
	}
	if view.Cursor.Position.Y != inputLine {
		t.Fatalf("cursor row %d did not match input row %d:\n%s", view.Cursor.Position.Y, inputLine, stripANSI(view.Content))
	}
	inputDisplayColumn := lipgloss.Width(lines[inputLine][:inputColumn])
	if view.Cursor.Position.X <= inputDisplayColumn {
		t.Fatalf("cursor column %d did not move past input column %d", view.Cursor.Position.X, inputDisplayColumn)
	}
	if wantMin := inputDisplayColumn + lipgloss.Width("한"); view.Cursor.Position.X < wantMin {
		t.Fatalf("cursor column %d did not account for wide Korean text; want at least %d", view.Cursor.Position.X, wantMin)
	}
}

func TestContextRailReflectsTUIState(t *testing.T) {
	model := newTUIModel("0.1.27")
	model.width = 180
	model.height = 36
	model.view.SetWidth(172)
	model.view.SetHeight(28)

	home := stripANSI(model.render())
	for _, expected := range []string{"Context", "Decision types"} {
		if !strings.Contains(home, expected) {
			t.Fatalf("home context missing %q:\n%s", expected, home)
		}
	}
	if strings.Contains(home, "No saved reports yet") {
		t.Fatalf("home context showed empty recent state:\n%s", home)
	}

	model.busy = true
	model.query = "Should I buy NVDA now?"
	model.phase = 1
	model.status = model.text.ResearchPhases[1]
	model.view.SetContent(researchingContent(model.query, model.text))
	research := stripANSI(model.render())
	for _, expected := range []string{"Current request", "Should I buy NVDA now?", "Research", "Sources"} {
		if !strings.Contains(research, expected) {
			t.Fatalf("research context missing %q:\n%s", expected, research)
		}
	}

	model.busy = false
	model.mode = model.text.Report
	model.status = model.text.Complete
	model.report = "## Easy Basic Report\n\n### 1. Bottom Line\nWait.\n\n### 3. Decision Rules\nBuy if...\n\n### 4. Biggest Risks\nValuation.\n\n### 5. Data Check\nCurrent price checked."
	model.view.SetContent(model.report)
	report := stripANSI(model.render())
	for _, expected := range []string{"Current request", "Report", "Bottom line", "Decision rules", "Risks", "Data check"} {
		if !strings.Contains(report, expected) {
			t.Fatalf("report context missing %q:\n%s", expected, report)
		}
	}
}

func TestReportActionsExposeNextChoices(t *testing.T) {
	got := stripANSI(renderReportActions(englishText()))
	normalized := strings.Join(strings.Fields(got), " ")
	for _, expected := range []string{"Next", "Use ↑/↓ and enter", "d Detailed analysis", "y Save", "n Skip", "ask a new question"} {
		if !strings.Contains(normalized, expected) {
			t.Fatalf("report actions missing %q:\n%s", expected, got)
		}
	}
}

func TestTUIHelpExposesOnlyUserFacingCommands(t *testing.T) {
	got := stripANSI(helpContent(englishText()))
	for _, expected := range []string{"/home", "/help", "/settings", "/model", "/language", "/theme", "/quit"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("help content missing %q:\n%s", expected, got)
		}
	}
	for _, hidden := range []string{"/doctor", "/setup", "/update"} {
		if strings.Contains(got, hidden) {
			t.Fatalf("help content exposed operational command %q:\n%s", hidden, got)
		}
	}
}

func TestFooterKeepsPrimaryCommandsVisible(t *testing.T) {
	got := stripANSI(renderFooter(120, englishText()))
	for _, expected := range []string{"/help", "/settings", "/theme", "/quit"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("footer missing %q:\n%s", expected, got)
		}
	}
	for _, hidden := range []string{"/doctor", "/setup", "/update"} {
		if strings.Contains(got, hidden) {
			t.Fatalf("footer exposed operational command %q:\n%s", hidden, got)
		}
	}
}
