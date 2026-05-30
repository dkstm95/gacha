package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	modelSettingAuto            = "auto"
	modelSettingOpenCodeDefault = "opencode-default"
	languageSettingAuto         = "auto"
	languageSettingEnglish      = "en"
	languageSettingKorean       = "ko"
	themeSettingSystem          = "system"
	themeSettingDark            = "dark"
	themeSettingLight           = "light"
	themeSettingGacha           = "gacha"
)

type gachaConfig struct {
	Model    string `json:"model"`
	Language string `json:"language"`
	Theme    string `json:"theme"`
}

func runConfigCommand(args []string) error {
	if len(args) == 0 || args[0] == "get" {
		return printConfig()
	}
	if args[0] != "set" {
		return fmt.Errorf("usage: gch config get | gch config set model <value> | gch config set language <auto|en|ko> | gch config set theme <system|dark|light|gacha>")
	}
	if len(args) < 3 {
		return fmt.Errorf("usage: gch config set model <value> | gch config set language <auto|en|ko> | gch config set theme <system|dark|light|gacha>")
	}
	key := strings.ToLower(strings.TrimSpace(args[1]))
	value := strings.TrimSpace(strings.Join(args[2:], " "))
	switch key {
	case "model":
		if err := updateConfigModel(value); err != nil {
			return err
		}
	case "language", "lang":
		if err := updateConfigLanguage(value); err != nil {
			return err
		}
	case "theme":
		if err := updateConfigTheme(value); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown config key %q", key)
	}
	fmt.Printf("Saved %s to %s\n", key, gachaConfigPath())
	return nil
}

func printConfig() error {
	config, err := configWithDefaults()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func updateConfigModel(value string) error {
	model, ok := normalizeModelSetting(value)
	if !ok {
		return fmt.Errorf("model must be auto, opencode-default, or provider/model")
	}
	config, err := loadGachaConfig()
	if err != nil {
		return err
	}
	config.Model = model
	return saveGachaConfig(config)
}

func updateConfigLanguage(value string) error {
	lang, ok := normalizeLanguageSetting(value)
	if !ok {
		return fmt.Errorf("language must be auto, en, or ko")
	}
	config, err := loadGachaConfig()
	if err != nil {
		return err
	}
	config.Language = lang
	return saveGachaConfig(config)
}

func updateConfigTheme(value string) error {
	theme, ok := normalizeThemeSetting(value)
	if !ok {
		return fmt.Errorf("theme must be system, dark, light, or gacha")
	}
	config, err := loadGachaConfig()
	if err != nil {
		return err
	}
	config.Theme = theme
	return saveGachaConfig(config)
}

func configWithDefaults() (gachaConfig, error) {
	config, err := loadGachaConfig()
	if err != nil {
		return gachaConfig{}, err
	}
	if strings.TrimSpace(config.Model) == "" {
		config.Model = modelSettingAuto
	}
	if strings.TrimSpace(config.Language) == "" {
		config.Language = languageSettingAuto
	}
	if strings.TrimSpace(config.Theme) == "" {
		config.Theme = themeSettingSystem
	}
	return config, nil
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

func saveGachaConfig(config gachaConfig) error {
	path := gachaConfigPath()
	if path == "" {
		return fmt.Errorf("Gacha config path not available")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o600)
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

func normalizeModelSetting(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}
	if strings.EqualFold(value, modelSettingAuto) {
		return modelSettingAuto, true
	}
	if strings.EqualFold(value, modelSettingOpenCodeDefault) {
		return modelSettingOpenCodeDefault, true
	}
	if strings.ContainsAny(value, " \t\n\r") || !strings.Contains(value, "/") {
		return "", false
	}
	return value, true
}

func validModelSetting(value string) bool {
	_, ok := normalizeModelSetting(value)
	return ok
}

func normalizeThemeSetting(value string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", themeSettingSystem, "auto":
		return themeSettingSystem, true
	case themeSettingDark:
		return themeSettingDark, true
	case themeSettingLight:
		return themeSettingLight, true
	case themeSettingGacha, "classic", "original":
		return themeSettingGacha, true
	default:
		return "", false
	}
}
