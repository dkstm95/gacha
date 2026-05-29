package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	modelSettingAuto            = "auto"
	modelSettingOpenCodeDefault = "opencode-default"
	defaultOpenAIChatGPTModel   = "openai/gpt-5.1-codex"
	modelFailureTTL             = 24 * time.Hour
)

type gachaConfig struct {
	Model string `json:"model"`
}

type modelResolution struct {
	Model     string
	Reason    string
	Source    string
	Fallbacks []string
}

type modelFailureCache struct {
	Failures map[string]time.Time `json:"failures"`
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
	if hasOpenAIChatGPTAuth() {
		known := []string{
			defaultOpenAIChatGPTModel,
			"openai/gpt-5.1-codex-mini",
		}
		candidates := filterFailedModels(known)
		reason := "auto: OpenAI OAuth ChatGPT account"
		if len(candidates) == 0 {
			candidates = known
			reason = "auto: all known OpenAI OAuth Codex models recently failed; retrying best known model"
		}
		return modelResolution{
			Model:     candidates[0],
			Fallbacks: candidates[1:],
			Reason:    reason,
			Source:    "OpenCode auth",
		}
	}
	return modelResolution{
		Reason: "auto: provider-specific default",
		Source: "OpenCode",
	}
}

func modelCandidates(resolution modelResolution) []string {
	if resolution.Model == "" {
		return nil
	}
	candidates := []string{resolution.Model}
	candidates = append(candidates, resolution.Fallbacks...)
	return candidates
}

func modelDescription(resolution modelResolution) string {
	if resolution.Model == "" {
		return fmt.Sprintf("OpenCode default (%s)", resolution.Reason)
	}
	if len(resolution.Fallbacks) == 0 {
		return fmt.Sprintf("%s (%s)", resolution.Model, resolution.Reason)
	}
	return fmt.Sprintf("%s (%s; fallback: %s)", resolution.Model, resolution.Reason, strings.Join(resolution.Fallbacks, ", "))
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

func gachaStateDir() string {
	base := os.Getenv("XDG_STATE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(base, "gacha")
}

func modelFailurePath() string {
	dir := gachaStateDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, "model-failures.json")
}

func filterFailedModels(models []string) []string {
	cache := loadModelFailureCache(time.Now())
	filtered := make([]string, 0, len(models))
	for _, model := range models {
		if _, failed := cache.Failures[model]; !failed {
			filtered = append(filtered, model)
		}
	}
	return filtered
}

func rememberModelFailure(model string) {
	model = strings.TrimSpace(model)
	if model == "" {
		return
	}
	now := time.Now()
	cache := loadModelFailureCache(now)
	if cache.Failures == nil {
		cache.Failures = map[string]time.Time{}
	}
	cache.Failures[model] = now.Add(modelFailureTTL)
	_ = saveModelFailureCache(cache)
}

func loadModelFailureCache(now time.Time) modelFailureCache {
	path := modelFailurePath()
	if path == "" {
		return modelFailureCache{Failures: map[string]time.Time{}}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return modelFailureCache{Failures: map[string]time.Time{}}
	}
	var cache modelFailureCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return modelFailureCache{Failures: map[string]time.Time{}}
	}
	if cache.Failures == nil {
		cache.Failures = map[string]time.Time{}
	}
	for model, expiresAt := range cache.Failures {
		if !expiresAt.After(now) {
			delete(cache.Failures, model)
		}
	}
	return cache
}

func saveModelFailureCache(cache modelFailureCache) error {
	path := modelFailurePath()
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}
