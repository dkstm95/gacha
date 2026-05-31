package app

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestPreferredOpenCodeModelAllowsEnvOverride(t *testing.T) {
	t.Setenv("GACHA_OPENCODE_MODEL", "anthropic/claude-sonnet-4-5")
	if got := resolveOpenCodeModel(); got.Model != "anthropic/claude-sonnet-4-5" {
		t.Fatalf("unexpected model: %#v", got)
	}
}

func TestOpenAIChatGPTAuthDetection(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dir)
	authDir := dir + "/opencode"
	if err := os.MkdirAll(authDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(authDir+"/auth.json", []byte(`{"openai":{"type":"oauth"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if !hasOpenAIChatGPTAuth() {
		t.Fatal("expected OpenAI OAuth auth to be detected")
	}
}

func TestResolveOpenCodeModelUsesConfigDefault(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("XDG_DATA_HOME", dir)
	t.Setenv("GACHA_OPENCODE_MODEL", "")
	configDir := dir + "/gacha"
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configDir+"/config.json", []byte(`{"model":"opencode-default"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got := resolveOpenCodeModel()
	if got.Model != "" || !strings.Contains(got.Reason, "OpenCode default") {
		t.Fatalf("unexpected model resolution: %#v", got)
	}
}

func TestSaveGachaConfigWritesModelAndLanguage(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	want := gachaConfig{Model: "opencode-default", Language: "ko", Theme: "dark"}
	if err := saveGachaConfig(want); err != nil {
		t.Fatal(err)
	}
	got, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if got.Model != want.Model || got.Language != want.Language || got.Theme != want.Theme {
		t.Fatalf("unexpected config: %#v", got)
	}
}

func TestRunConfigCommandSetsModelAndLanguage(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	if err := runConfigCommand([]string{"set", "model", "OpenCode-Default"}); err != nil {
		t.Fatal(err)
	}
	if err := runConfigCommand([]string{"set", "language", "ko"}); err != nil {
		t.Fatal(err)
	}
	if err := runConfigCommand([]string{"set", "theme", "light"}); err != nil {
		t.Fatal(err)
	}
	got, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if got.Model != "opencode-default" || got.Language != "ko" || got.Theme != "light" {
		t.Fatalf("unexpected config: %#v", got)
	}
}

func TestRunConfigCommandRejectsInvalidValues(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if err := runConfigCommand([]string{"set", "model", "bad-model"}); err == nil {
		t.Fatal("expected invalid model error")
	}
	if err := runConfigCommand([]string{"set", "language", "fr"}); err == nil {
		t.Fatal("expected invalid language error")
	}
	if err := runConfigCommand([]string{"set", "theme", "neon"}); err == nil {
		t.Fatal("expected invalid theme error")
	}
}

func TestConfigWithDefaults(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	config, err := configWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	if config.Model != modelSettingAuto || config.Language != languageSettingAuto || config.Theme != themeSettingSystem {
		t.Fatalf("unexpected default config: %#v", config)
	}
}

func TestExistingConfigWithoutProfileLoads(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	configDir := dir + "/gacha"
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configDir+"/config.json", []byte(`{"model":"auto","language":"en","theme":"system"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !profileIsZero(config.Profile) {
		t.Fatalf("expected missing profile to load as empty, got %#v", config.Profile)
	}
}

func TestConfigJSONOmitsEmptyProfileAndIncludesSavedProfile(t *testing.T) {
	empty, err := json.Marshal(gachaConfig{Model: "auto", Language: "en", Theme: "system"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(empty), "profile") {
		t.Fatalf("empty profile should be omitted: %s", empty)
	}
	withProfile, err := json.Marshal(gachaConfig{
		Model:    "auto",
		Language: "en",
		Theme:    "system",
		Profile: gachaProfile{
			Markets:    profileMulti{Selected: []string{profileMarketUSStocksETFs}, Default: profileMarketUSStocksETFs},
			Onboarding: profileOnboarding{Completed: true},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(withProfile), `"profile"`) || !strings.Contains(string(withProfile), profileMarketUSStocksETFs) {
		t.Fatalf("saved profile should be included: %s", withProfile)
	}
}

func TestSaveGachaConfigPreservesProfile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	want := gachaConfig{
		Model:    "auto",
		Language: "en",
		Theme:    "system",
		Profile: gachaProfile{
			Markets: profileMulti{
				Selected: []string{profileMarketUSStocksETFs, profileMarketKoreanStocksETFs},
				Default:  profileMarketUSStocksETFs,
			},
			Horizons: profileMulti{
				Selected: []string{profileHorizonSixToTwelve},
				Default:  profileHorizonSixToTwelve,
			},
			Risk:        profileRiskBalanced,
			ReportStyle: profileReportBasicFirst,
			Goals: profileMulti{
				Selected: []string{profileGoalDiscover, profileGoalEntry},
				Default:  profileGoalEntry,
			},
			Onboarding: profileOnboarding{Completed: true},
		},
	}
	if err := saveGachaConfig(want); err != nil {
		t.Fatal(err)
	}
	got, err := loadGachaConfig()
	if err != nil {
		t.Fatal(err)
	}
	if got.Profile.Markets.Default != profileMarketUSStocksETFs || len(got.Profile.Markets.Selected) != 2 {
		t.Fatalf("unexpected profile markets: %#v", got.Profile.Markets)
	}
	if got.Profile.Horizons.Default != profileHorizonSixToTwelve || got.Profile.Risk != profileRiskBalanced || got.Profile.ReportStyle != profileReportBasicFirst {
		t.Fatalf("unexpected profile: %#v", got.Profile)
	}
	if !got.Profile.Onboarding.Completed {
		t.Fatalf("expected completed onboarding: %#v", got.Profile.Onboarding)
	}
}

func TestNormalizeProfileMakesNotSureExclusive(t *testing.T) {
	got := normalizeProfile(gachaProfile{
		Markets: profileMulti{
			Selected: []string{profileMarketUSStocksETFs, profileValueNotSure},
			Default:  profileMarketUSStocksETFs,
		},
		Horizons: profileMulti{
			Selected: []string{profileHorizonSixToTwelve, profileValueNotSure},
			Default:  profileHorizonSixToTwelve,
		},
	})
	if len(got.Markets.Selected) != 1 || got.Markets.Selected[0] != profileValueNotSure || got.Markets.Default != "" {
		t.Fatalf("not sure should be exclusive for markets: %#v", got.Markets)
	}
	if len(got.Horizons.Selected) != 1 || got.Horizons.Selected[0] != profileValueNotSure || got.Horizons.Default != "" {
		t.Fatalf("not sure should be exclusive for horizons: %#v", got.Horizons)
	}
}

func TestAppProfileCommandPrintsSummary(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if err := saveGachaConfig(gachaConfig{Profile: gachaProfile{
		Markets: profileMulti{Selected: []string{profileMarketUSStocksETFs}, Default: profileMarketUSStocksETFs},
		Risk:    profileRiskBalanced,
	}}); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	app := New("test")
	app.env.Stdout = &stdout
	if err := app.Run([]string{"profile"}); err != nil {
		t.Fatal(err)
	}
	got := stdout.String()
	if !strings.Contains(got, "Research Profile") || !strings.Contains(got, "US stocks / ETFs") || !strings.Contains(got, "Balanced") {
		t.Fatalf("profile summary missing expected content:\n%s", got)
	}
}
