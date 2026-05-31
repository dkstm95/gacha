package app

import (
	"strings"
	"testing"
)

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
