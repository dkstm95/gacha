package app

import (
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
