package app

import (
	"os"
	"strings"
	"testing"
)

func TestBuildPromptIncludesWorkflowAndRequirements(t *testing.T) {
	prompt, err := buildPrompt([]string{"NVDA"})
	if err != nil {
		t.Fatal(err)
	}

	for _, expected := range []string{
		"# gacha auto",
		"Workflow library:",
		"# gacha entry",
		"User request:\nNVDA",
		"Always use current web search",
		"Provenance Appendix",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt did not contain %q", expected)
		}
	}
}

func TestBuildPromptAutoClassifiesAndRequiresFreshData(t *testing.T) {
	prompt, err := buildPrompt([]string{"NVDA", "지금", "살까?"})
	if err != nil {
		t.Fatal(err)
	}

	for _, expected := range []string{
		"# gacha auto",
		"Classify the user's request",
		"Always use current web search",
		"even if the user does not ask",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt did not contain %q", expected)
		}
	}
}

func TestShellQuote(t *testing.T) {
	got := shellQuote("can't")
	if got != "'can'\\''t'" {
		t.Fatalf("unexpected quote: %s", got)
	}
}

func TestOpenCodeRunArgsUsesExplicitModel(t *testing.T) {
	got := openCodeRunArgs("hello", "openai/gpt-5.1-codex")
	want := []string{"run", "--model", "openai/gpt-5.1-codex", "hello"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected args: %#v", got)
	}
}

func TestPreferredOpenCodeModelAllowsEnvOverride(t *testing.T) {
	t.Setenv("GACHA_OPENCODE_MODEL", "anthropic/claude-sonnet-4-5")
	if got := resolveOpenCodeModel(); got.Model != "anthropic/claude-sonnet-4-5" {
		t.Fatalf("unexpected model: %#v", got)
	}
}

func TestOpenAIChatGPTAuthDetection(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dir)
	t.Setenv("GACHA_OPENCODE_MODEL", "")
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
	if got := resolveOpenCodeModel(); got.Model != defaultOpenAIChatGPTModel {
		t.Fatalf("unexpected model: %#v", got)
	}
}

func TestFallbackOpenCodeModelForUnsupportedChatGPTModel(t *testing.T) {
	output := `Bad Request: {"detail":"The 'gpt-5.5-pro' model is not supported when using Codex with a ChatGPT account."}`
	got := fallbackOpenCodeModel([]string{"run", "--model", "openai/gpt-5.5-pro", "hello"}, output)
	if got != defaultOpenAIChatGPTModel {
		t.Fatalf("unexpected fallback: %s", got)
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

func TestModelFailureCacheSkipsFailedAutoCandidate(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dir)
	t.Setenv("XDG_STATE_HOME", dir)
	t.Setenv("GACHA_OPENCODE_MODEL", "")
	authDir := dir + "/opencode"
	if err := os.MkdirAll(authDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(authDir+"/auth.json", []byte(`{"openai":{"type":"oauth"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	rememberModelFailure(defaultOpenAIChatGPTModel)
	got := resolveOpenCodeModel()
	if got.Model != "openai/gpt-5.1-codex-mini" {
		t.Fatalf("unexpected model after failure cache: %#v", got)
	}
}
