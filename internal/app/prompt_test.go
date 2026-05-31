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
		"Easy Basic Report",
		"Source Log",
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

func TestBuildPromptUsesKoreanForKoreanQuestion(t *testing.T) {
	prompt, err := buildPrompt([]string{"NVDA", "지금", "사도", "될까?"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(prompt, "Write the final report in Korean") {
		t.Fatalf("prompt did not request Korean response")
	}
}

func TestBuildPromptLocksReportStructure(t *testing.T) {
	prompt, err := buildPrompt([]string{"What should I buy?"})
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Report structure contract:",
		"The final answer must start with the Easy Basic Report",
		"Write for ordinary users, not investment professionals.",
		"The Basic Report is decision-ready",
		"The Detailed Analysis is the verification layer",
		"Include Detailed Analysis only when",
		"Use simple tables only when they make the answer easier to compare",
		"## Easy Basic Report",
		"### 3. Decision Rules",
		"### 5. Data Check",
		"### 6. More Detail Options",
		"the user knows what they can ask for next",
		"## Detailed Analysis",
		"### G. Source Log",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt did not contain report contract detail %q", expected)
		}
	}
}

func TestBuildDetailedPromptContinuesFromBasicReport(t *testing.T) {
	prompt, err := buildDetailedPrompt("NVDA 지금 사도 될까?", "## Easy Basic Report\n\nWatch")
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"# gacha detailed analysis",
		"User request:\nNVDA 지금 사도 될까?",
		"Existing basic report:",
		"Write the detailed analysis in Korean",
		`Produce only the "Detailed Analysis" section`,
		"Use these headings in order",
		"Use current web search or current market-data tools again",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("detailed prompt did not contain %q", expected)
		}
	}
}
