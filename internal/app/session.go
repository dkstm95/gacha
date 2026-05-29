package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (a *App) startSession() error {
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
