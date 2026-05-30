package app

import (
	"encoding/json"
	"fmt"
	"strings"
)

func runConfigCommand(args []string) error {
	if len(args) == 0 || args[0] == "get" {
		return printConfig()
	}
	if args[0] != "set" {
		return fmt.Errorf("usage: gch config get | gch config set model <value> | gch config set language <auto|en|ko>")
	}
	if len(args) < 3 {
		return fmt.Errorf("usage: gch config set model <value> | gch config set language <auto|en|ko>")
	}
	key := strings.ToLower(strings.TrimSpace(args[1]))
	value := strings.TrimSpace(strings.Join(args[2:], " "))
	config, err := loadGachaConfig()
	if err != nil {
		return err
	}
	switch key {
	case "model":
		if !validModelSetting(value) {
			return fmt.Errorf("model must be auto, opencode-default, or provider/model")
		}
		config.Model = value
	case "language", "lang":
		lang, ok := normalizeLanguageSetting(value)
		if !ok {
			return fmt.Errorf("language must be auto, en, or ko")
		}
		config.Language = lang
	default:
		return fmt.Errorf("unknown config key %q", key)
	}
	if err := saveGachaConfig(config); err != nil {
		return err
	}
	fmt.Printf("Saved %s to %s\n", key, gachaConfigPath())
	return nil
}

func printConfig() error {
	config, err := loadGachaConfig()
	if err != nil {
		return err
	}
	if strings.TrimSpace(config.Model) == "" {
		config.Model = modelSettingAuto
	}
	if strings.TrimSpace(config.Language) == "" {
		config.Language = languageSettingAuto
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
