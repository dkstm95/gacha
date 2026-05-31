package app

import (
	"os"
	"strings"
	"testing"
)

func TestParseOpenCodeModels(t *testing.T) {
	output := "\x1b[0m\nopenai/gpt-5.5\nopenai/gpt-5.5-pro\nnot a model\nopenai/gpt-5.5\n"
	got := parseOpenCodeModels(output)
	want := []string{"openai/gpt-5.5", "openai/gpt-5.5-pro"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected models: %#v", got)
	}
}

func TestChooseModelPrefersQualityAndLatest(t *testing.T) {
	got := firstModel(rankModels([]string{
		"openai/gpt-5.5-mini",
		"openai/gpt-5.1-codex",
		"openai/gpt-5.5-pro",
		"openai/gpt-5.5",
		"openai/gpt-5",
	}, "openai"))
	if got != "openai/gpt-5.5" {
		t.Fatalf("unexpected model: %s", got)
	}
}

func TestRankModelsPrefersNewestOpenAIBaseFrontier(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dir)
	authDir := dir + "/opencode"
	if err := os.MkdirAll(authDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(authDir+"/auth.json", []byte(`{"openai":{"type":"oauth"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got := rankModels([]string{
		"openai/gpt-5.6-pro",
		"openai/gpt-5.6-fast",
		"openai/gpt-5.6",
		"openai/gpt-5.5-pro",
		"openai/gpt-5.5",
		"openai/gpt-5.1-codex",
	}, "openai")
	if got[0] != "openai/gpt-5.6" {
		t.Fatalf("unexpected ranked models: %#v", got)
	}
}

func TestRankModelsKeepsOpenAICodexBehindFrontierBase(t *testing.T) {
	got := rankModels([]string{
		"openai/gpt-5.3-codex-spark",
		"openai/gpt-5.3-codex",
		"openai/gpt-5.5",
		"openai/gpt-5.5-fast",
	}, "openai")
	if got[0] != "openai/gpt-5.5" {
		t.Fatalf("unexpected ranked models: %#v", got)
	}
}

func TestChooseModelPenalizesSmallModels(t *testing.T) {
	got := chooseModel([]string{
		"google/gemini-3-flash",
		"google/gemini-2.5-pro",
	})
	if got != "google/gemini-2.5-pro" {
		t.Fatalf("unexpected model: %s", got)
	}
}

func TestUnsupportedChatGPTModelErrorDetection(t *testing.T) {
	output := `Bad Request: {"detail":"The 'gpt-5.5-pro' model is not supported when using Codex with a ChatGPT account."}`
	if !isUnsupportedChatGPTModelError(output) {
		t.Fatal("expected unsupported ChatGPT model error")
	}
}
