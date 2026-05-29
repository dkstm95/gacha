package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (a *App) startSession() error {
	reader := bufio.NewReader(os.Stdin)
	interactive := isInteractiveTerminal()
	if interactive {
		enterScreen()
		defer exitScreen()
	}
	printSessionHome(interactive, a.version)

	for {
		fmt.Print("\nAsk > ")
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
			if interactive {
				fmt.Println("\nGoodbye.")
			} else {
				fmt.Println("Goodbye.")
			}
			return nil
		case "/h", "/help", "help":
			printSessionHelp(interactive)
			continue
		case "/home", "home":
			printSessionHome(interactive, a.version)
			continue
		case "/doctor", "doctor":
			if err := doctor(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		case "/update", "update":
			if err := a.updateSelf(); err != nil {
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

func printSessionHome(interactive bool, version string) {
	if interactive {
		clearScreen()
	}
	cfg := loadConfig()
	selected := selectPlatform(cfg, cfg.DefaultPlatform)
	platform := cfg.Platforms[selected]
	platformLabel := platform.Label
	if selected == "manual" {
		platformLabel = "Copy/paste prompt"
	}

	fmt.Println("INVESTIQ")
	fmt.Println("Fresh-data investment research for your AI tools")
	fmt.Println()
	fmt.Println("+------------------------------------------------------------+")
	fmt.Println("| Ask a question. investiq will classify it automatically.   |")
	fmt.Println("| It always asks the AI to use current web or market data.   |")
	fmt.Println("+------------------------------------------------------------+")
	fmt.Println()
	fmt.Printf("Route        %s\n", platformLabel)
	fmt.Printf("Config       %s\n", configPath())
	fmt.Printf("Version      %s\n", version)
	fmt.Println()
	fmt.Println("Try")
	fmt.Println("  NVDA 지금 사도 될까?")
	fmt.Println("  What should I invest in for the next 6 to 12 months?")
	fmt.Println("  I own TSLA. When should I trim, sell, or stop out?")
	fmt.Println()
	fmt.Println("Commands")
	fmt.Println("  /help      show commands")
	fmt.Println("  /doctor    check AI tools")
	fmt.Println("  /update    update investiq")
	fmt.Println("  /quit      exit")
}

func printSessionHelp(interactive bool) {
	if interactive {
		clearScreen()
	}
	fmt.Println()
	fmt.Println("INVESTIQ HELP")
	fmt.Println()
	fmt.Println("Commands")
	fmt.Println("  /home       Show the home screen")
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

func enterScreen() {
	fmt.Print("\033[?1049h")
	clearScreen()
}

func exitScreen() {
	fmt.Print("\033[?1049l")
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func isInteractiveTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
