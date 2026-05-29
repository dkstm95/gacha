package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	modelSettingAuto            = "auto"
	modelSettingOpenCodeDefault = "opencode-default"
)

type gachaConfig struct {
	Model string `json:"model"`
}

type modelResolution struct {
	Model  string
	Reason string
	Source string
}

func resolveOpenCodeModel() modelResolution {
	if model := strings.TrimSpace(os.Getenv("GACHA_OPENCODE_MODEL")); model != "" {
		return modelResolution{
			Model:  model,
			Reason: "env override",
			Source: "GACHA_OPENCODE_MODEL",
		}
	}

	config, _ := loadGachaConfig()
	configModel := strings.TrimSpace(config.Model)
	switch {
	case configModel == "", strings.EqualFold(configModel, modelSettingAuto):
		return autoOpenCodeModel()
	case strings.EqualFold(configModel, modelSettingOpenCodeDefault):
		return modelResolution{
			Reason: "configured to use OpenCode default",
			Source: gachaConfigPath(),
		}
	default:
		return modelResolution{
			Model:  configModel,
			Reason: "configured custom model",
			Source: gachaConfigPath(),
		}
	}
}

func autoOpenCodeModel() modelResolution {
	provider := preferredOpenCodeProvider()
	if provider == "" {
		return modelResolution{
			Reason: "auto: no connected provider detected",
			Source: "OpenCode default",
		}
	}

	models, err := discoverOpenCodeModels(provider)
	if err != nil || len(models) == 0 {
		reason := "auto: could not read provider model list"
		if err != nil {
			reason += ": " + firstLine(err.Error())
		}
		return modelResolution{
			Reason: reason,
			Source: "OpenCode default",
		}
	}

	selected := chooseModel(models)
	if selected == "" {
		return modelResolution{
			Reason: "auto: no usable model found",
			Source: "OpenCode default",
		}
	}
	return modelResolution{
		Model:  selected,
		Reason: "auto: selected from current OpenCode model list",
		Source: "opencode models " + provider,
	}
}

func firstLine(value string) string {
	value = strings.TrimSpace(value)
	if index := strings.IndexByte(value, '\n'); index >= 0 {
		return strings.TrimSpace(value[:index])
	}
	return value
}

func preferredOpenCodeProvider() string {
	providers, err := openCodeAuthProviders()
	if err != nil || len(providers) == 0 {
		return ""
	}
	for _, preferred := range []string{"openai", "anthropic", "google", "gemini"} {
		if _, ok := providers[preferred]; ok {
			return preferred
		}
	}
	names := make([]string, 0, len(providers))
	for provider := range providers {
		names = append(names, provider)
	}
	sort.Strings(names)
	return names[0]
}

func discoverOpenCodeModels(provider string) ([]string, error) {
	commandPath, ok := resolveCommand(openCodeCommand)
	if !ok {
		return nil, fmt.Errorf("OpenCode runtime not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, commandPath, "models", provider)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, errors.New(strings.TrimSpace(stripANSI(out.String())))
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return parseOpenCodeModels(out.String()), nil
}

func parseOpenCodeModels(output string) []string {
	var models []string
	seen := map[string]bool{}
	for _, line := range strings.Split(stripANSI(output), "\n") {
		model := strings.TrimSpace(line)
		if model == "" || !strings.Contains(model, "/") || seen[model] {
			continue
		}
		seen[model] = true
		models = append(models, model)
	}
	return models
}

func chooseModel(models []string) string {
	if len(models) == 0 {
		return ""
	}
	sorted := append([]string(nil), models...)
	sort.SliceStable(sorted, func(i, j int) bool {
		left := modelScore(sorted[i])
		right := modelScore(sorted[j])
		if left == right {
			return sorted[i] > sorted[j]
		}
		return left > right
	})
	return sorted[0]
}

func modelScore(model string) int {
	name := strings.ToLower(model)
	score := 0

	for _, token := range []string{"pro", "opus", "ultra", "max"} {
		if strings.Contains(name, token) {
			score += 70
		}
	}
	for _, token := range []string{"codex", "sonnet", "gpt", "gemini", "claude"} {
		if strings.Contains(name, token) {
			score += 30
		}
	}
	score += newestVersionScore(name)

	for _, token := range []string{"mini", "nano", "lite", "flash", "haiku"} {
		if strings.Contains(name, token) {
			score -= 80
		}
	}
	return score
}

var modelVersionPattern = regexp.MustCompile(`\d+(?:\.\d+)?`)

func newestVersionScore(name string) int {
	matches := modelVersionPattern.FindAllString(name, -1)
	best := 0
	for _, match := range matches {
		parts := strings.SplitN(match, ".", 2)
		major, _ := strconv.Atoi(parts[0])
		minor := 0
		if len(parts) == 2 {
			minor, _ = strconv.Atoi(parts[1])
		}
		value := major*100 + minor
		if value > best {
			best = value
		}
	}
	return best
}

func modelCandidates(resolution modelResolution) []string {
	if resolution.Model == "" {
		return nil
	}
	return []string{resolution.Model}
}

func modelDescription(resolution modelResolution) string {
	if resolution.Model == "" {
		return fmt.Sprintf("OpenCode default (%s)", resolution.Reason)
	}
	return fmt.Sprintf("%s (%s)", resolution.Model, resolution.Reason)
}

func loadGachaConfig() (gachaConfig, error) {
	path := gachaConfigPath()
	if path == "" {
		return gachaConfig{}, fmt.Errorf("Gacha config path not available")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return gachaConfig{}, nil
		}
		return gachaConfig{}, err
	}
	var config gachaConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return gachaConfig{}, err
	}
	return config, nil
}

func gachaConfigPath() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "gacha", "config.json")
}
