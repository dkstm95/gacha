package app

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
	"testing"
)

func TestThemeContentShowsInteractivePreviews(t *testing.T) {
	for _, tc := range []struct {
		name     string
		text     uiText
		expected []string
	}{
		{
			name: "english",
			text: englishText(),
			expected: []string{
				"Themes",
				"Active: System",
				"Use ↑/↓ and enter",
				"System",
				"Dark",
				"Light",
				"Gacha",
				"Preview",
				"select this theme",
				"Better odds through research.",
				"Profile: US stocks / ETFs",
				"Ask an investment question.",
			},
		},
		{
			name: "korean",
			text: koreanText(),
			expected: []string{
				"테마",
				"현재: 시스템",
				"↑/↓로 이동하고 enter로 선택",
				"시스템",
				"다크",
				"라이트",
				"Gacha",
				"예시",
				"이 테마 선택",
				"Better odds through research.",
				"Profile: US stocks / ETFs",
				"Ask an investment question.",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			got := stripANSI(themeContent(tc.text))
			for _, expected := range tc.expected {
				if !strings.Contains(got, expected) {
					t.Fatalf("theme content missing %q:\n%s", expected, got)
				}
			}
			for _, line := range strings.Split(got, "\n") {
				if lipgloss.Width(line) > 80 {
					t.Fatalf("line width %d exceeds 80: %q\n%s", lipgloss.Width(line), line, got)
				}
			}
			if strings.Contains(got, "Decision desk") {
				t.Fatalf("theme preview should match the prompt-first home:\n%s", got)
			}
		})
	}
}

func TestThemeContentReflectsSavedTheme(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if err := updateConfigTheme(themeSettingLight); err != nil {
		t.Fatal(err)
	}
	got := stripANSI(themeContent(englishText()))
	if !strings.Contains(got, "Active: Light") {
		t.Fatalf("theme content did not reflect saved theme:\n%s", got)
	}
}

func TestTUIThemeCommandSavesAndAppliesTheme(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	model := newTUIModel("0.1.27")
	next, cmd := model.handleThemeSetting("/theme light")
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	updated := next.(tuiModel)
	got, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if got.Theme != themeSettingLight {
		t.Fatalf("unexpected theme config: %#v", got)
	}
	if !strings.Contains(stripANSI(updated.view.View()), "Active: Light") {
		t.Fatalf("theme view did not show previews:\n%s", stripANSI(updated.view.View()))
	}
}

func TestTUIThemeCommandOpensSelectableThemeChoice(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	model := newTUIModel("0.1.27")

	next, cmd := model.handleSubmit("/theme")
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	updated := next.(tuiModel)
	if updated.choice == nil || updated.choice.Kind != choiceTheme {
		t.Fatalf("expected theme choice, got %#v", updated.choice)
	}

	next, cmd = updated.Update(keyMsg("down"))
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	updated = next.(tuiModel)
	if updated.choice.Selected != 1 {
		t.Fatalf("expected second theme selected, got %d", updated.choice.Selected)
	}

	next, cmd = updated.Update(keyMsg("enter"))
	if cmd != nil {
		t.Fatal("unexpected command")
	}
	updated = next.(tuiModel)
	config, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.Theme != themeSettingDark {
		t.Fatalf("unexpected theme config: %#v", config)
	}
	if updated.choice == nil || updated.choice.Selected != 1 {
		t.Fatalf("expected theme choice to stay active after selection: %#v", updated.choice)
	}
}
