package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func selectPlatform(cfg Config, requested string) string {
	if requested != "" && requested != "auto" {
		return requested
	}
	for _, name := range cfg.PlatformPriority {
		platform := cfg.Platforms[name]
		if !platform.Enabled {
			continue
		}
		if name == "manual" || hasRunnableCommand(platform.Command) {
			return name
		}
	}
	return "manual"
}

func runPlatform(name string, platform PlatformConfig, prompt string, dryRun bool) error {
	if name == "manual" || platform.PromptMode == "print" || platform.Command == "" {
		fmt.Println(prompt)
		return nil
	}
	if !hasRunnableCommand(platform.Command) {
		return fmt.Errorf("platform command not found: %s\nRun `iq doctor` to check platform routing.", platform.Command)
	}
	args := renderArgs(platform.Args, prompt)
	if dryRun {
		parts := []string{platform.Command}
		for _, arg := range args {
			parts = append(parts, shellQuote(arg))
		}
		fmt.Println(strings.Join(parts, " "))
		return nil
	}
	cmd := exec.Command(platform.Command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func renderArgs(args []string, prompt string) []string {
	rendered := make([]string, len(args))
	for i, arg := range args {
		rendered[i] = strings.ReplaceAll(arg, "{{prompt}}", prompt)
	}
	return rendered
}

func hasCommand(command string) bool {
	if command == "" {
		return false
	}
	_, err := exec.LookPath(command)
	return err == nil
}

func hasRunnableCommand(command string) bool {
	if !hasCommand(command) {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return ctx.Err() == nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("cannot open browser on %s", runtime.GOOS)
	}
	return cmd.Start()
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
