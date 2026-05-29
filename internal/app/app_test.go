package app

import (
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
