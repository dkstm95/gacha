package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func initConfig(yes bool) error {
	cfg := defaultConfig()
	if !yes {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("investiq will use the best available AI CLI on this machine.")
		fmt.Println("If none are ready, it falls back to printing a prompt you can paste into any web AI.")
		fmt.Println()
		for _, name := range cfg.PlatformPriority {
			if name == "manual" {
				continue
			}
			platform := cfg.Platforms[name]
			installed := hasRunnableCommand(platform.Command)
			if installed {
				fmt.Printf("Use %s when available? [Y/n] ", platform.Label)
				answer, _ := reader.ReadString('\n')
				normalized := strings.ToLower(strings.TrimSpace(answer))
				platform.Enabled = normalized != "n" && normalized != "no"
			} else {
				platform.Enabled = false
				if platform.SetupURL != "" {
					fmt.Printf("%s was not found. Open setup page? [y/N] ", platform.Label)
					answer, _ := reader.ReadString('\n')
					normalized := strings.ToLower(strings.TrimSpace(answer))
					if normalized == "y" || normalized == "yes" {
						_ = openBrowser(platform.SetupURL)
					}
				} else {
					fmt.Printf("%s was not found. Skipping.\n", platform.Label)
				}
			}
			cfg.Platforms[name] = platform
		}
	}

	if err := os.MkdirAll(configDir(), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(configPath(), append(data, '\n'), 0o644); err != nil {
		return err
	}
	fmt.Printf("Wrote %s\n", configPath())
	fmt.Println()
	fmt.Println("Try:")
	fmt.Println(`  iq "NVDA 지금 사도 될까?"`)
	return nil
}

func doctor() error {
	cfg := loadConfig()
	if _, err := os.Stat(configPath()); err == nil {
		fmt.Printf("Config: %s\n\n", configPath())
	} else {
		fmt.Println("Config: (not created; using defaults)")
		fmt.Println()
	}
	for _, name := range cfg.PlatformPriority {
		platform := cfg.Platforms[name]
		installed := name == "manual" || hasRunnableCommand(platform.Command)
		status := "disabled"
		if platform.Enabled && installed {
			status = "ready"
		} else if platform.Enabled {
			status = "missing"
		}
		fmt.Printf("%-9s %-8s %s\n", name, status, platform.Label)
		if platform.Command != "" {
			fmt.Printf("           command: %s\n", platform.Command)
		}
		if platform.Subscription != "" {
			fmt.Printf("           subscription: %s\n", platform.Subscription)
		}
	}
	return nil
}

func platforms() error {
	cfg := loadConfig()
	data, err := json.MarshalIndent(cfg.Platforms, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func defaultConfig() Config {
	platforms := map[string]PlatformConfig{
		"claude": {
			Label:      "Claude Code",
			Command:    "claude",
			Args:       []string{"-p", "{{prompt}}"},
			PromptMode: "argument",
			SetupURL:   "https://docs.anthropic.com/en/docs/claude-code/setup",
		},
		"codex": {
			Label:      "Codex",
			Command:    "codex",
			Args:       []string{"{{prompt}}"},
			PromptMode: "argument",
			SetupURL:   "https://developers.openai.com/codex/cli",
		},
		"opencode": {
			Label:      "OpenCode / Oh My OpenAgent",
			Command:    "opencode",
			Args:       []string{"run", "{{prompt}}"},
			PromptMode: "argument",
			SetupURL:   "https://opencode.ai/",
		},
		"gemini": {
			Label:      "Gemini CLI",
			Command:    "gemini",
			Args:       []string{"{{prompt}}"},
			PromptMode: "argument",
			SetupURL:   "https://github.com/google-gemini/gemini-cli",
		},
		"manual": {
			Label:        "Manual copy/paste",
			PromptMode:   "print",
			Subscription: "manual",
			Enabled:      true,
		},
	}
	for name, platform := range platforms {
		if platform.Command != "" && hasRunnableCommand(platform.Command) {
			platform.Enabled = true
			platforms[name] = platform
		}
	}
	return Config{
		Version:             1,
		DefaultPlatform:     "auto",
		PlatformPriority:    []string{"claude", "codex", "opencode", "gemini", "manual"},
		RequireFreshData:    true,
		AllowTradeExecution: false,
		Platforms:           platforms,
	}
}

func loadConfig() Config {
	cfg := defaultConfig()
	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}
	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		return cfg
	}
	if loaded.Version != 0 {
		cfg.Version = loaded.Version
	}
	if loaded.DefaultPlatform != "" {
		cfg.DefaultPlatform = loaded.DefaultPlatform
	}
	if len(loaded.PlatformPriority) > 0 {
		cfg.PlatformPriority = loaded.PlatformPriority
	}
	cfg.RequireFreshData = loaded.RequireFreshData
	cfg.AllowTradeExecution = loaded.AllowTradeExecution
	for name, platform := range loaded.Platforms {
		cfg.Platforms[name] = platform
	}
	return cfg
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".investiq"
	}
	return filepath.Join(home, ".investiq")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}
