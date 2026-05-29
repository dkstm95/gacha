package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

const openCodeCommand = "opencode"
const defaultOpenAIChatGPTCodexModel = "openai/gpt-5.1-codex"

type openCodeAuthCredential struct {
	Type string `json:"type"`
}

func runAgent(prompt string, dryRun bool) error {
	commandPath, ok := resolveCommand(openCodeCommand)
	if !ok || !hasOpenCodeAuth() {
		fmt.Println(prompt)
		return nil
	}
	args := openCodeRunArgs(prompt, preferredOpenCodeModel())
	if dryRun {
		parts := []string{commandPath}
		for _, arg := range args {
			parts = append(parts, shellQuote(arg))
		}
		fmt.Println(strings.Join(parts, " "))
		return nil
	}
	output, err := runOpenCode(commandPath, args, true)
	if err == nil {
		return nil
	}
	if fallback := fallbackOpenCodeModel(args, output); fallback != "" {
		fmt.Fprintf(os.Stderr, "\nOpenCode rejected the selected model. Retrying with %s...\n\n", fallback)
		_, retryErr := runOpenCode(commandPath, openCodeRunArgs(prompt, fallback), true)
		return retryErr
	}
	return err
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

func openCodeRunArgs(prompt string, model string) []string {
	args := []string{"run"}
	if model != "" {
		args = append(args, "--model", model)
	}
	return append(args, prompt)
}

func preferredOpenCodeModel() string {
	if model := strings.TrimSpace(os.Getenv("GACHA_OPENCODE_MODEL")); model != "" {
		return model
	}
	if hasOpenAIChatGPTAuth() {
		return defaultOpenAIChatGPTCodexModel
	}
	return ""
}

func hasOpenAIChatGPTAuth() bool {
	providers, err := openCodeAuthProviders()
	if err != nil {
		return false
	}
	credential, ok := providers["openai"]
	return ok && strings.EqualFold(credential.Type, "oauth")
}

func openCodeAuthProviders() (map[string]openCodeAuthCredential, error) {
	path := openCodeAuthPath()
	if path == "" {
		return nil, fmt.Errorf("OpenCode auth path not available")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var providers map[string]openCodeAuthCredential
	if err := json.Unmarshal(data, &providers); err != nil {
		return nil, err
	}
	return providers, nil
}

func runOpenCode(commandPath string, args []string, stream bool) (string, error) {
	cmd := exec.Command(commandPath, args...)
	cmd.Stdin = os.Stdin
	if stream {
		var out bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &out)
		cmd.Stderr = io.MultiWriter(os.Stderr, &out)
		err := cmd.Run()
		return out.String(), err
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func fallbackOpenCodeModel(args []string, output string) string {
	if !isUnsupportedChatGPTCodexModelError(output) {
		return ""
	}
	current := modelFromOpenCodeArgs(args)
	for _, candidate := range []string{
		defaultOpenAIChatGPTCodexModel,
		"openai/gpt-5.1-codex-mini",
		"openai/gpt-5-codex",
	} {
		if candidate != current {
			return candidate
		}
	}
	return ""
}

func isUnsupportedChatGPTCodexModelError(output string) bool {
	normalized := strings.ToLower(output)
	return strings.Contains(normalized, "model is not supported") &&
		strings.Contains(normalized, "chatgpt account")
}

func modelFromOpenCodeArgs(args []string) string {
	for i, arg := range args {
		if (arg == "--model" || arg == "-m") && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(arg, "--model=") {
			return strings.TrimPrefix(arg, "--model=")
		}
	}
	return ""
}
