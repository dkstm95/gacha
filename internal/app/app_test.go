package app

import (
	"strings"
	"testing"
)

func TestBuildPromptIncludesWorkflowAndRequirements(t *testing.T) {
	prompt, err := buildPrompt("entry", []string{"NVDA"})
	if err != nil {
		t.Fatal(err)
	}

	for _, expected := range []string{
		"# investiq entry",
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
	prompt, err := buildPrompt("auto", []string{"NVDA", "지금", "살까?"})
	if err != nil {
		t.Fatal(err)
	}

	for _, expected := range []string{
		"# investiq auto",
		"Classify the user's request",
		"Always use current web search",
		"even if the user does not ask",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt did not contain %q", expected)
		}
	}
}

func TestSelectPlatformFallsBackToManual(t *testing.T) {
	cfg := Config{
		Version:          1,
		DefaultPlatform:  "auto",
		PlatformPriority: []string{"missing", "manual"},
		Platforms: map[string]PlatformConfig{
			"missing": {
				Label:      "Missing",
				Command:    "definitely-not-installed-investiq-test",
				Args:       []string{"{{prompt}}"},
				PromptMode: "argument",
				Enabled:    true,
			},
			"manual": {
				Label:      "Manual",
				PromptMode: "print",
				Enabled:    true,
			},
		},
	}

	if got := selectPlatform(cfg, "auto"); got != "manual" {
		t.Fatalf("expected manual fallback, got %q", got)
	}
}

func TestRenderArgs(t *testing.T) {
	got := renderArgs([]string{"-p", "{{prompt}}"}, "hello")
	if len(got) != 2 || got[0] != "-p" || got[1] != "hello" {
		t.Fatalf("unexpected args: %#v", got)
	}
}
