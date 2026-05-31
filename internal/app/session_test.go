package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestLineSessionAcceptsDocumentedSlashQuit(t *testing.T) {
	output := runLineSessionWithInput(t, "/quit\n")
	if strings.Contains(output, "Unknown command") {
		t.Fatalf("line session treated /quit as unknown:\n%s", output)
	}
	if !strings.Contains(output, "Goodbye.") {
		t.Fatalf("line session did not exit cleanly:\n%s", output)
	}
}

func TestLineSessionShowsHelpForUnknownSlashCommand(t *testing.T) {
	output := runLineSessionWithInput(t, "/unknown\n/quit\n")
	if !strings.Contains(output, "Unknown command") {
		t.Fatalf("line session did not explain unknown slash command:\n%s", output)
	}
	if !strings.Contains(output, "Command palette") {
		t.Fatalf("line session did not show help for unknown slash command:\n%s", output)
	}
}

func runLineSessionWithInput(t *testing.T, input string) string {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GACHA_LANG", "en")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app := New("test")
	app.env.Stdin = strings.NewReader(input)
	app.env.Stdout = &stdout
	app.env.Stderr = &stderr
	if err := app.startLineSession(); err != nil {
		t.Fatal(err)
	}
	return stdout.String()
}
