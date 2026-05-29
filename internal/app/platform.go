package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const openCodeCommand = "opencode"

func runAgent(prompt string, dryRun bool) error {
	commandPath, ok := resolveCommand(openCodeCommand)
	if !ok || !hasOpenCodeAuth() {
		fmt.Println(prompt)
		return nil
	}
	args := []string{"run", prompt}
	if dryRun {
		parts := []string{commandPath}
		for _, arg := range args {
			parts = append(parts, shellQuote(arg))
		}
		fmt.Println(strings.Join(parts, " "))
		return nil
	}
	cmd := exec.Command(commandPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func hasCommand(command string) bool {
	if command == "" {
		return false
	}
	_, ok := resolveCommand(command)
	return ok
}

func hasRunnableCommand(command string) bool {
	commandPath, ok := resolveCommand(command)
	if !ok {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, commandPath, "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return ctx.Err() == nil
}

func resolveCommand(command string) (string, bool) {
	if command == "" {
		return "", false
	}
	if path, err := exec.LookPath(command); err == nil {
		return path, true
	}
	if command == openCodeCommand {
		if path := managedOpenCodePath(); path != "" {
			if _, err := os.Stat(path); err == nil {
				return path, true
			}
		}
	}
	return "", false
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
