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
	text := textFor(detectLanguage())
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("gacha")
	if detectLanguage() == languageKorean {
		fmt.Println("질문을 입력하세요. 종료하려면 /quit을 입력하세요.")
	} else {
		fmt.Println("Type a question, or /quit to exit.")
	}
	for {
		fmt.Print("\n" + strings.TrimSuffix(text.InputPlaceholder, "...") + " > ")
		line, err := reader.ReadString('\n')
		if err != nil && strings.TrimSpace(line) == "" {
			fmt.Println()
			return nil
		}
		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}
		if isSettingsCommand(input) {
			fmt.Println(settingsContent(text))
			continue
		}
		if isSlashCommand(input) {
			fmt.Println(fmt.Sprintf(text.UnknownCommand, input))
			fmt.Println()
			fmt.Println(helpContent(text))
			continue
		}
		switch input {
		case "/q", "/quit", "quit", "exit":
			if detectLanguage() == languageKorean {
				fmt.Println("종료합니다.")
			} else {
				fmt.Println("Goodbye.")
			}
			return nil
		case "/doctor", "doctor":
			if err := doctor(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		case "/setup", "setup":
			if err := ensureRuntime(true); err != nil {
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
