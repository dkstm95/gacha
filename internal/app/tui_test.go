package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"testing"
)

func keyMsg(value string) tea.KeyMsg {
	switch value {
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)}
}

func skipProfileOnboardingForTest(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	if err := saveGachaConfig(gachaConfig{Profile: gachaProfile{Onboarding: profileOnboarding{Skipped: true}}}); err != nil {
		t.Fatal(err)
	}
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
	skipProfileOnboardingForTest(t)
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

func TestWelcomeContentIsPromptFirst(t *testing.T) {
	skipProfileOnboardingForTest(t)
	got := stripANSI(welcomeContent("0.1.27", englishText(), 80, 16))
	for _, expected := range []string{
		"GACHA",
		"Better odds through research.",
		"Ask an investment question.",
		"No research profile set. Type /profile",
		"Discover opportunities",
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("welcome content missing %q:\n%s", expected, got)
		}
	}
	for _, unwanted := range []string{
		"Ask -> Fresh data",
		"[Fresh data required]",
		"[No trading]",
		"Decision desk",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("welcome content kept decorative element %q:\n%s", unwanted, got)
		}
	}
}

func TestWelcomeContentFitsCommonTerminalSizes(t *testing.T) {
	skipProfileOnboardingForTest(t)
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
			skipProfileOnboardingForTest(t)
			model := newTUIModel("0.1.27")
			model.width = tc.width
			model.height = tc.height
			model.view.Width = max(30, tc.width-8)
			model.view.Height = max(6, tc.height-8)
			got := stripANSI(model.View())
			for _, line := range strings.Split(got, "\n") {
				if lipgloss.Width(line) > tc.width {
					t.Fatalf("line width %d exceeds %d: %q\n%s", lipgloss.Width(line), tc.width, line, got)
				}
			}
			if strings.Contains(got, "Checks fresh data") {
				t.Fatalf("status bar repeated safety copy:\n%s", got)
			}
			if tc.name == "full" {
				for _, expected := range []string{"GACHA", "Better odds through research.", "Ask an investment question.", "Ready"} {
					if !strings.Contains(got, expected) {
						t.Fatalf("full layout missing %q:\n%s", expected, got)
					}
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
			got := stripANSI(updated.View())
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
	got := stripANSI(resized.(tuiModel).View())
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
			skipProfileOnboardingForTest(t)
			model := newTUIModel("0.1.27")
			model.text = tc.text
			model.input.Placeholder = tc.text.InputPlaceholder
			model.width = 180
			model.height = 36
			got := stripANSI(model.View())
			if !strings.Contains(got, tc.promptText) {
				t.Fatalf("missing localized prompt %q:\n%s", tc.promptText, got)
			}
		})
	}
}

func TestTUIStatusBarRendersBelowWorkspace(t *testing.T) {
	skipProfileOnboardingForTest(t)
	model := newTUIModel("0.1.27")
	model.width = 180
	model.height = 36
	got := stripANSI(model.View())
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
		t.Fatalf("missing panel bottom border:\n%s", got)
	}
}

func TestTUIOnboardingStartsWhenProfileMissing(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	model := newTUIModel("0.1.27")
	model.width = 180
	model.height = 36

	got := stripANSI(model.View())
	for _, expected := range []string{"Better odds through research.", "Primary markets", "space toggle"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("onboarding missing %q:\n%s", expected, got)
		}
	}
	researchLine := -1
	profileLine := -1
	lines := strings.Split(got, "\n")
	for i, line := range lines {
		if strings.Contains(line, "disciplined research.") {
			researchLine = i
		}
		if strings.Contains(line, "Let's set your research profile.") {
			profileLine = i
		}
	}
	if researchLine < 0 || profileLine != researchLine+2 || strings.Trim(strings.Trim(lines[researchLine+1], "│"), " ") != "" {
		t.Fatalf("onboarding intro should keep separate paragraphs:\n%s", got)
	}
}

func TestTUIOnboardingCanSaveProfile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	model := newTUIModel("0.1.27")
	for _, key := range []string{" ", "enter", " ", "enter", "enter", "enter", " ", "enter"} {
		next, cmd := model.Update(keyMsg(key))
		if cmd != nil {
			t.Fatal("unexpected command")
		}
		model = next.(tuiModel)
	}
	config, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !config.Profile.Onboarding.Completed || config.Profile.Onboarding.Skipped {
		t.Fatalf("unexpected onboarding state: %#v", config.Profile.Onboarding)
	}
	if config.Profile.Markets.Default != profileMarketUSStocksETFs || config.Profile.Horizons.Default != profileHorizonOneToThreeMonths {
		t.Fatalf("unexpected saved profile: %#v", config.Profile)
	}
	if config.Profile.Risk != profileRiskBalanced || config.Profile.ReportStyle != profileReportBasicFirst {
		t.Fatalf("unexpected saved single-select profile: %#v", config.Profile)
	}
}

func TestTUIOnboardingSkipSavesSkippedState(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	model := newTUIModel("0.1.27")
	for range profileOnboardingSteps {
		next, cmd := model.Update(keyMsg("s"))
		if cmd != nil {
			t.Fatal("unexpected command")
		}
		model = next.(tuiModel)
	}
	config, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.Profile.Onboarding.Completed || !config.Profile.Onboarding.Skipped {
		t.Fatalf("expected skipped onboarding state, got %#v", config.Profile.Onboarding)
	}
	model = newTUIModel("0.1.27")
	if model.profile != nil {
		t.Fatalf("skipped onboarding should not reopen, got %#v", model.profile)
	}
}

func TestTUIOnboardingRejectsEmptyMultiSelect(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	model := newTUIModel("0.1.27")
	next, cmd := model.Update(keyMsg("enter"))
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	updated := next.(tuiModel)
	if updated.profile == nil || updated.profile.Category != profileCategoryMarkets {
		t.Fatalf("expected onboarding to stay on markets, got %#v", updated.profile)
	}
	got := stripANSI(updated.View())
	if !strings.Contains(got, "Select at least one option") {
		t.Fatalf("missing inline validation:\n%s", got)
	}
	config, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !profileIsZero(config.Profile) {
		t.Fatalf("invalid category should not save profile: %#v", config.Profile)
	}
}

func TestTUIProfileCommandOpensEditor(t *testing.T) {
	skipProfileOnboardingForTest(t)
	model := newTUIModel("0.1.27")
	next, cmd := model.handleSubmit("/profile")
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	updated := next.(tuiModel)
	if updated.profile == nil || updated.profile.Mode != profileFlowMenu {
		t.Fatalf("expected profile menu, got %#v", updated.profile)
	}
	got := stripANSI(updated.view.View())
	for _, expected := range []string{"Research Profile", "Edit markets", "Reset profile"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("profile editor missing %q:\n%s", expected, got)
		}
	}
}

func TestTUIProfileEditorPersistsEditedCategory(t *testing.T) {
	skipProfileOnboardingForTest(t)
	model := newTUIModel("0.1.27")
	next, cmd := model.handleSubmit("/profile")
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	model = next.(tuiModel)
	for _, key := range []string{"enter", " ", "enter"} {
		next, cmd = model.Update(keyMsg(key))
		if cmd != nil {
			t.Fatal("unexpected command")
		}
		model = next.(tuiModel)
	}
	config, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !config.Profile.Onboarding.Completed || config.Profile.Onboarding.Skipped {
		t.Fatalf("edited profile should become active, got %#v", config.Profile.Onboarding)
	}
	if config.Profile.Markets.Default != profileMarketUSStocksETFs {
		t.Fatalf("expected edited market to persist, got %#v", config.Profile.Markets)
	}
}

func TestTUIProfileEditorResetPreservesSystemSettings(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if err := saveGachaConfig(gachaConfig{
		Model:    modelSettingOpenCodeDefault,
		Language: languageSettingKorean,
		Theme:    themeSettingGacha,
		Profile: gachaProfile{
			Markets:    profileMulti{Selected: []string{profileMarketUSStocksETFs}, Default: profileMarketUSStocksETFs},
			Risk:       profileRiskBalanced,
			Onboarding: profileOnboarding{Completed: true},
		},
	}); err != nil {
		t.Fatal(err)
	}
	model := newTUIModel("0.1.27")
	next, cmd := model.handleSubmit("/profile")
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	model = next.(tuiModel)
	for i := 0; i < len(profileMenuLabels(languageEnglish))-1; i++ {
		next, cmd = model.Update(keyMsg("down"))
		if cmd != nil {
			t.Fatal("unexpected command")
		}
		model = next.(tuiModel)
	}
	next, cmd = model.Update(keyMsg("enter"))
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	config, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.Model != modelSettingOpenCodeDefault || config.Language != languageSettingKorean || config.Theme != themeSettingGacha {
		t.Fatalf("reset changed system settings: %#v", config)
	}
	if !profileIsZero(config.Profile) {
		t.Fatalf("reset should clear profile, got %#v", config.Profile)
	}
}

func TestReportActionsExposeNextChoices(t *testing.T) {
	got := stripANSI(renderReportActions(englishText()))
	normalized := strings.Join(strings.Fields(got), " ")
	for _, expected := range []string{"Next", "Use ↑/↓ and enter", "Detailed analysis", "Save report", "Skip saving", "Shortcut: d", "ask a new question"} {
		if !strings.Contains(normalized, expected) {
			t.Fatalf("report actions missing %q:\n%s", expected, got)
		}
	}
}

func TestReportContextRecognizesKoreanSections(t *testing.T) {
	got := reportContextFromMarkdown("## 쉬운 기본 리포트\n\n### 1. 결론\n대기.\n\n### 3. 행동 기준\n조건.\n\n### 4. 리스크\n높은 밸류에이션.\n\n### 5. 데이터 확인\n출처 포함.", "fallback")
	for _, expected := range []string{"결론", "행동 기준", "리스크", "데이터"} {
		found := false
		for _, item := range got {
			if item == expected {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing Korean report context %q in %#v", expected, got)
		}
	}
}

func TestTUIHelpExposesOnlyUserFacingCommands(t *testing.T) {
	got := stripANSI(helpContent(englishText()))
	for _, expected := range []string{"/home", "/help", "/profile", "/settings", "/model", "/language", "/theme", "/quit"} {
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
	for _, expected := range []string{"/profile", "/settings", "/quit"} {
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
