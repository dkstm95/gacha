package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var version = "0.1.0"

//go:embed plugins/investiq/platforms/generic/system-prompt.md
//go:embed plugins/investiq/templates/investment-report.md
//go:embed plugins/investiq/workflows/*.md
var embedded embed.FS

var modes = map[string]bool{
	"discover":  true,
	"select":    true,
	"entry":     true,
	"exit":      true,
	"portfolio": true,
	"journal":   true,
}

type PlatformConfig struct {
	Label        string   `json:"label"`
	Command      string   `json:"command"`
	Args         []string `json:"args"`
	PromptMode   string   `json:"promptMode"`
	Subscription string   `json:"subscription"`
	Enabled      bool     `json:"enabled"`
}

type Config struct {
	Version             int                       `json:"version"`
	DefaultPlatform     string                    `json:"defaultPlatform"`
	PlatformPriority    []string                  `json:"platformPriority"`
	RequireFreshData    bool                      `json:"requireFreshData"`
	AllowTradeExecution bool                      `json:"allowTradeExecution"`
	Platforms           map[string]PlatformConfig `json:"platforms"`
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		printUsage()
		return nil
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Println(version)
		return nil
	case "init":
		return initConfig(hasFlag(args[1:], "--yes") || hasFlag(args[1:], "-y"))
	case "doctor":
		return doctor()
	case "platforms":
		return platforms()
	case "prompt":
		if len(args) < 2 {
			return errors.New("missing mode")
		}
		return printPrompt(args[1], args[2:])
	case "run":
		if len(args) < 2 {
			return errors.New("missing mode")
		}
		return runMode(args[1], args[2:])
	default:
		if modes[args[0]] {
			return printPrompt(args[0], args[1:])
		}
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func printUsage() {
	fmt.Println(`investiq

Usage:
  investiq init [--yes]                       Create ~/.investiq/config.json
  investiq doctor                             Check configured AI platforms
  investiq platforms                          Print platform config
  investiq prompt <mode> [request]            Print the composed agent prompt
  investiq run <mode> [request] [options]     Run through a configured platform
  investiq <mode> [request]                   Alias for prompt <mode>

Modes:
  discover, select, entry, exit, portfolio, journal

Options:
  --platform auto|claude|codex|opencode|gemini|manual
  --dry-run

Examples:
  investiq init
  investiq doctor
  investiq run entry "NVDA current entry zone" --platform auto
  investiq run discover "latest opportunities for a 12 month horizon" --platform manual`)
}

func initConfig(yes bool) error {
	cfg := defaultConfig()
	if !yes {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Configure the AI platforms you subscribe to or can run locally.")
		for _, name := range cfg.PlatformPriority {
			if name == "manual" {
				continue
			}
			platform := cfg.Platforms[name]
			installed := hasCommand(platform.Command)
			defaultAnswer := "[y/N]"
			if installed {
				defaultAnswer = "[Y/n]"
			}
			fmt.Printf("Enable %s (%s)? %s ", platform.Label, platform.Command, defaultAnswer)
			answer, _ := reader.ReadString('\n')
			normalized := strings.ToLower(strings.TrimSpace(answer))
			platform.Enabled = (installed && normalized != "n" && normalized != "no") || (!installed && (normalized == "y" || normalized == "yes"))
			if platform.Enabled {
				fmt.Printf("Subscription/account label for %s (optional): ", platform.Label)
				subscription, _ := reader.ReadString('\n')
				platform.Subscription = strings.TrimSpace(subscription)
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
		installed := name == "manual" || hasCommand(platform.Command)
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

func printPrompt(mode string, query []string) error {
	prompt, err := buildPrompt(mode, query)
	if err != nil {
		return err
	}
	fmt.Println(prompt)
	return nil
}

func runMode(mode string, args []string) error {
	if !modes[mode] {
		return fmt.Errorf("invalid mode: %s", mode)
	}
	platformName := "auto"
	dryRun := false
	query := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			dryRun = true
		case "--platform":
			if i+1 >= len(args) {
				return errors.New("--platform requires a value")
			}
			platformName = args[i+1]
			i++
		default:
			query = append(query, args[i])
		}
	}
	cfg := loadConfig()
	selected := selectPlatform(cfg, platformName)
	platform, ok := cfg.Platforms[selected]
	if !ok {
		return fmt.Errorf("unknown platform: %s", selected)
	}
	prompt, err := buildPrompt(mode, query)
	if err != nil {
		return err
	}
	return runPlatform(selected, platform, prompt, dryRun)
}

func buildPrompt(mode string, queryParts []string) (string, error) {
	if !modes[mode] {
		return "", fmt.Errorf("unknown mode: %s", mode)
	}
	system, err := readEmbedded("plugins/investiq/platforms/generic/system-prompt.md")
	if err != nil {
		return "", err
	}
	template, err := readEmbedded("plugins/investiq/templates/investment-report.md")
	if err != nil {
		return "", err
	}
	workflowPath := fmt.Sprintf("plugins/investiq/workflows/%s.md", mode)
	workflow, err := readEmbedded(workflowPath)
	if err != nil {
		return "", err
	}
	query := strings.TrimSpace(strings.Join(queryParts, " "))
	if query == "" {
		query = "(No additional user request supplied.)"
	}
	sections := []string{
		strings.TrimSpace(system),
		strings.TrimSpace(workflow),
		"User request:\n" + query,
		"Report template:\n" + strings.TrimSpace(template),
		`Hard requirements:
- Use current web search or current market-data tools before analysis.
- If fresh data cannot be verified, do not make a recommendation.
- Include data freshness, source links, risks, Devil's Advocate, action conditions, monitoring plan, and provenance.
- Do not execute trades. The final decision remains with the user.`,
	}
	return strings.Join(sections, "\n\n"), nil
}

func readEmbedded(name string) (string, error) {
	data, err := fs.ReadFile(embedded, name)
	if err != nil {
		return "", err
	}
	return string(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))), nil
}

func defaultConfig() Config {
	platforms := map[string]PlatformConfig{
		"claude": {
			Label:      "Claude Code",
			Command:    "claude",
			Args:       []string{"-p", "{{prompt}}"},
			PromptMode: "argument",
		},
		"codex": {
			Label:      "Codex",
			Command:    "codex",
			Args:       []string{"{{prompt}}"},
			PromptMode: "argument",
		},
		"opencode": {
			Label:      "OpenCode / Oh My OpenAgent",
			Command:    "opencode",
			Args:       []string{"run", "{{prompt}}"},
			PromptMode: "argument",
		},
		"gemini": {
			Label:      "Gemini CLI",
			Command:    "gemini",
			Args:       []string{"{{prompt}}"},
			PromptMode: "argument",
		},
		"manual": {
			Label:        "Manual copy/paste",
			PromptMode:   "print",
			Subscription: "manual",
			Enabled:      true,
		},
	}
	for name, platform := range platforms {
		if platform.Command != "" && hasCommand(platform.Command) {
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

func selectPlatform(cfg Config, requested string) string {
	if requested != "" && requested != "auto" {
		return requested
	}
	for _, name := range cfg.PlatformPriority {
		platform := cfg.Platforms[name]
		if !platform.Enabled {
			continue
		}
		if name == "manual" || hasCommand(platform.Command) {
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
	if !hasCommand(platform.Command) {
		return fmt.Errorf("platform command not found: %s\nRun `investiq doctor` or use `--platform manual`.", platform.Command)
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

func hasCommand(command string) bool {
	if command == "" {
		return false
	}
	_, err := exec.LookPath(command)
	return err == nil
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}

func targetTriple() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}
