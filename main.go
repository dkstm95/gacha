package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var version = "0.1.0"

//go:embed plugins/investiq/platforms/generic/system-prompt.md
//go:embed plugins/investiq/templates/investment-report.md
//go:embed plugins/investiq/workflows/*.md
var embedded embed.FS

var modes = map[string]bool{
	"auto":      true,
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
	SetupURL     string   `json:"setupUrl"`
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

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return startSession()
	}
	if args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
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
	case "update":
		return updateSelf()
	case "platforms":
		return platforms()
	case "prompt":
		if len(args) >= 2 && modes[args[1]] {
			return printPrompt(args[1], args[2:])
		}
		return printPrompt("auto", args[1:])
	case "run":
		if len(args) >= 2 && modes[args[1]] {
			return runMode(args[1], args[2:])
		}
		return runMode("auto", args[1:])
	default:
		if modes[args[0]] {
			return runMode(args[0], args[1:])
		}
		return runMode("auto", args)
	}
}

func printUsage() {
	fmt.Println(`investiq

Usage:
  investiq                                    Open the interactive investiq UI
  iq                                          Open the interactive investiq UI
  iq init                                     Set up AI platform routing
  iq doctor                                   Check detected AI platforms
  iq update                                   Update investiq to the latest release
  iq "question"                               Analyze with automatic request classification

Debug:
  investiq prompt "question"                  Print the composed agent prompt
  investiq platforms                          Print platform config

Examples:
  iq "NVDA 지금 사도 될까?"
  iq "AAPL 현재 매수 구간 분석"
  iq "TSLA 보유 중인데 매도 기준 점검"
  iq "6개월에서 12개월 관점 투자 후보 찾아줘"`)
}

func startSession() error {
	reader := bufio.NewReader(os.Stdin)
	printSessionHeader()

	for {
		fmt.Print("\niq> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if strings.TrimSpace(line) == "" {
				fmt.Println()
				return nil
			}
		}
		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}
		switch input {
		case "/q", "/quit", "quit", "exit":
			fmt.Println("Goodbye.")
			return nil
		case "/h", "/help", "help":
			printSessionHelp()
			continue
		case "/doctor", "doctor":
			if err := doctor(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		case "/update", "update":
			if err := updateSelf(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		case "/init", "init":
			if err := initConfig(false); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		case "/platforms":
			if err := platforms(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		fmt.Println()
		if err := runMode("auto", []string{input}); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func printSessionHeader() {
	fmt.Println("investiq")
	fmt.Println("Fresh-data investment research agent")
	fmt.Println()
	fmt.Println("Ask an investment question. investiq will classify it and route it automatically.")
	fmt.Println("Type /help for commands, /doctor to check AI platforms, /quit to exit.")
}

func printSessionHelp() {
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  /help       Show this help")
	fmt.Println("  /doctor     Check detected AI platforms")
	fmt.Println("  /update     Update investiq to the latest release")
	fmt.Println("  /init       Configure AI platform routing")
	fmt.Println("  /platforms  Print platform config")
	fmt.Println("  /quit       Exit")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println(`  NVDA 지금 사도 될까?`)
	fmt.Println(`  What should I invest in for the next 6 to 12 months?`)
	fmt.Println(`  I own TSLA. When should I trim, sell, or stop out?`)
}

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

func updateSelf() error {
	latest, err := latestReleaseTag()
	if err != nil {
		return err
	}
	current := normalizeVersion(version)
	target := normalizeVersion(latest)
	if current == target {
		fmt.Printf("investiq is already up to date (%s).\n", version)
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return err
	}

	fmt.Printf("Updating investiq %s -> %s\n", version, target)
	tmpDir, err := os.MkdirTemp("", "investiq-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, "investiq.tar.gz")
	url := releaseAssetURL(latest)
	if err := downloadFile(url, archivePath); err != nil {
		return err
	}
	if err := extractTarGz(archivePath, tmpDir); err != nil {
		return err
	}
	newBinary := filepath.Join(tmpDir, "investiq")
	if err := os.Chmod(newBinary, 0o755); err != nil {
		return err
	}

	backup := exe + ".old"
	_ = os.Remove(backup)
	if err := os.Rename(exe, backup); err != nil {
		return fmt.Errorf("cannot replace %s: %w", exe, err)
	}
	if err := copyFile(newBinary, exe, 0o755); err != nil {
		_ = os.Rename(backup, exe)
		return err
	}
	_ = os.Remove(backup)

	aliasPath := filepath.Join(filepath.Dir(exe), "iq")
	if _, err := os.Lstat(aliasPath); err == nil {
		_ = os.Remove(aliasPath)
		_ = os.Symlink(filepath.Base(exe), aliasPath)
	}

	fmt.Printf("Updated %s to %s.\n", exe, target)
	return nil
}

func latestReleaseTag() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/dkstm95/investiq/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "investiq/"+version)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("GitHub release check failed: %s", resp.Status)
	}
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	if release.TagName == "" {
		return "", fmt.Errorf("GitHub latest release did not include a tag")
	}
	return release.TagName, nil
}

func normalizeVersion(value string) string {
	return strings.TrimPrefix(strings.TrimSpace(value), "v")
}

func releaseAssetURL(tag string) string {
	return fmt.Sprintf("https://github.com/dkstm95/investiq/releases/download/%s/investiq-%s.tar.gz", tag, targetTriple())
}

func downloadFile(url string, destination string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "investiq/"+version)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.ReadFrom(resp.Body)
	return err
}

func extractTarGz(archivePath string, destinationDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		name := filepath.Clean(header.Name)
		if name != "investiq" {
			continue
		}
		target := filepath.Join(destinationDir, "investiq")
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, tarReader)
		closeErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		return nil
	}
}

func copyFile(source string, destination string, mode fs.FileMode) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Chmod(destination, mode)
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
	dryRun := false
	query := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			dryRun = true
		default:
			query = append(query, args[i])
		}
	}
	cfg := loadConfig()
	selected := selectPlatform(cfg, cfg.DefaultPlatform)
	platform, ok := cfg.Platforms[selected]
	if !ok {
		return fmt.Errorf("unknown platform: %s", selected)
	}
	prompt, err := buildPrompt(mode, query)
	if err != nil {
		return err
	}
	if err := runPlatform(selected, platform, prompt, dryRun); err != nil {
		if selected == "manual" {
			return err
		}
		fmt.Fprintf(os.Stderr, "Could not use %s automatically: %v\n", platform.Label, err)
		fmt.Fprintln(os.Stderr, "Falling back to a prompt you can paste into any web AI.")
		fmt.Fprintln(os.Stderr)
		return runPlatform("manual", cfg.Platforms["manual"], prompt, dryRun)
	}
	return nil
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
	workflow := `# investiq auto

Classify the user's request into discover, select, entry, exit, portfolio, or journal, then follow the matching investiq workflow. If the request is ambiguous, choose the safest interpretation and state your assumption.`
	if mode != "auto" {
		workflowPath := fmt.Sprintf("plugins/investiq/workflows/%s.md", mode)
		workflow, err = readEmbedded(workflowPath)
		if err != nil {
			return "", err
		}
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
- Always use current web search or current market-data tools before analysis, even if the user does not ask for latest/current/recent data.
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
