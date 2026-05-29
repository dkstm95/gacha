package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (a *App) startSession() error {
	if !isInteractiveTerminal() {
		return a.startLineSession()
	}
	program := tea.NewProgram(newTUIModel(a.version), tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func (a *App) startLineSession() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("gacha")
	fmt.Println("Type a question, or /quit to exit.")
	for {
		fmt.Print("\nAsk > ")
		line, err := reader.ReadString('\n')
		if err != nil && strings.TrimSpace(line) == "" {
			fmt.Println()
			return nil
		}
		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}
		switch input {
		case "/q", "/quit", "quit", "exit":
			fmt.Println("Goodbye.")
			return nil
		case "/doctor", "doctor":
			if err := doctor(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		default:
			if err := runQuery([]string{input}); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func isInteractiveTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
