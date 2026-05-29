package app

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const openCodeInstallCommand = "curl -fsSL https://opencode.ai/install | bash"

func ensureRuntime(interactive bool) error {
	if hasRunnableCommand("opencode") {
		if interactive {
			return ensureProviderConnected()
		}
		return nil
	}
	if !interactive {
		return nil
	}
	if runtime.GOOS == "windows" {
		fmt.Println("OpenCode runtime was not found.")
		fmt.Println("Automatic OpenCode setup is currently supported on macOS and Linux.")
		fmt.Println("See https://opencode.ai/docs/ for Windows setup.")
		return nil
	}
	if !hasCommand("curl") || !hasCommand("bash") {
		fmt.Println("OpenCode runtime was not found.")
		fmt.Println("Install curl and bash, then run /setup again.")
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Gacha needs a local AI runtime to run research inside this UI.")
	fmt.Print("Install OpenCode runtime now? [Y/n] ")
	answer, _ := reader.ReadString('\n')
	normalized := strings.ToLower(strings.TrimSpace(answer))
	if normalized == "n" || normalized == "no" {
		fmt.Println("Skipping runtime setup. Gacha will print a copy/paste prompt when no AI runtime is ready.")
		return nil
	}
	if err := installOpenCode(); err != nil {
		return err
	}
	return ensureProviderConnected()
}

func installOpenCode() error {
	fmt.Println("Installing OpenCode runtime...")
	cmd := exec.Command("bash", "-lc", openCodeInstallCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("OpenCode install failed: %w", err)
	}
	if hasRunnableCommand("opencode") {
		fmt.Println("OpenCode runtime is ready.")
		return nil
	}
	if path := managedOpenCodePath(); path != "" {
		fmt.Println("OpenCode was installed, but its bin directory is not on PATH yet.")
		fmt.Printf("Run this now:\n  export PATH=\"%s:$PATH\"\n", filepath.Dir(path))
	}
	return nil
}

func ensureProviderConnected() error {
	if hasOpenCodeAuth() {
		printConnectedProviders()
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("OpenCode runtime is installed, but no AI provider is connected yet.")
	fmt.Print("Connect ChatGPT, GitHub Copilot, Gemini, OpenAI API, or another provider now? [Y/n] ")
	answer, _ := reader.ReadString('\n')
	normalized := strings.ToLower(strings.TrimSpace(answer))
	if normalized == "n" || normalized == "no" {
		fmt.Println("Skipping provider connection. Gacha will print a copy/paste prompt until a provider is connected.")
		return nil
	}

	commandPath, ok := resolveCommand("opencode")
	if !ok {
		return fmt.Errorf("OpenCode runtime was not found after install")
	}

	fmt.Println("Starting OpenCode provider login...")
	fmt.Println("Choose a provider and complete the prompts. Credentials are stored by OpenCode.")
	cmd := exec.Command(commandPath, "auth", "login")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("OpenCode provider login failed: %w", err)
	}

	if hasOpenCodeAuth() {
		fmt.Println("AI provider connection is ready.")
		printConnectedProviders()
		return nil
	}
	fmt.Println("Provider login finished, but Gacha could not confirm stored credentials yet.")
	fmt.Printf("Check OpenCode credentials at: %s\n", openCodeAuthPath())
	return nil
}

func printConnectedProviders() {
	providers, err := openCodeAuthList()
	if err != nil || strings.TrimSpace(providers) == "" {
		fmt.Println("OpenCode credentials were found.")
		return
	}
	fmt.Println("OpenCode credentials:")
	fmt.Println(strings.TrimSpace(providers))
}

func openCodeAuthList() (string, error) {
	commandPath, ok := resolveCommand("opencode")
	if !ok {
		return "", fmt.Errorf("OpenCode runtime not found")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, commandPath, "auth", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func managedOpenCodePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".opencode", "bin", "opencode")
}

func openCodeAuthPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "opencode", "auth.json")
	}
	return filepath.Join(home, ".local", "share", "opencode", "auth.json")
}

func hasOpenCodeAuth() bool {
	path := openCodeAuthPath()
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.Size() > 2
}
