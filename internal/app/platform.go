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

type openCodeAuthCredential struct {
	Type string `json:"type"`
}

func runAgent(prompt string, dryRun bool) (string, bool, error) {
	commandPath, ok := resolveCommand(openCodeCommand)
	if !ok || !hasOpenCodeAuth() {
		fmt.Println(prompt)
		return "", false, nil
	}
	resolution := resolveOpenCodeModel()
	args := openCodeRunArgs(prompt, resolution.Model)
	if dryRun {
		parts := []string{commandPath}
		for _, arg := range args {
			parts = append(parts, shellQuote(arg))
		}
		fmt.Println(strings.Join(parts, " "))
		return "", false, nil
	}
	output, err := runOpenCodeWithResolution(commandPath, prompt, resolution, true)
	return strings.TrimSpace(output), err == nil && strings.TrimSpace(output) != "", err
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

func runOpenCodeWithResolution(commandPath string, prompt string, resolution modelResolution, stream bool) (string, error) {
	candidates := modelCandidates(resolution)
	if len(candidates) == 0 {
		return runOpenCode(commandPath, openCodeRunArgs(prompt, ""), stream)
	}

	var lastOutput string
	var lastErr error
	for index, model := range candidates {
		if index > 0 && stream {
			fmt.Fprintf(os.Stderr, "\nOpenCode rejected the selected model. Retrying with %s...\n\n", model)
		}
		output, err := runOpenCode(commandPath, openCodeRunArgs(prompt, model), stream)
		if err == nil {
			return output, nil
		}
		lastOutput = output
		lastErr = err
		if !isUnsupportedChatGPTModelError(output) {
			return output, err
		}
	}
	return lastOutput, lastErr
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

func isUnsupportedChatGPTModelError(output string) bool {
	normalized := strings.ToLower(output)
	return strings.Contains(normalized, "model is not supported") &&
		strings.Contains(normalized, "chatgpt account")
}
