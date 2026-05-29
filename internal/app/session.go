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
	if err := ensureRuntime(interactive); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	printSessionHome(interactive, a.version)

	for {
		printAskPrompt(interactive)
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
			fmt.Println()
			if err := doctor(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		case "/setup", "setup":
			fmt.Println()
			if err := ensureRuntime(true); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		case "/update", "update":
			fmt.Println()
			if err := a.updateSelf(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		fmt.Println()
		if err := runQuery([]string{input}); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func printSessionHome(interactive bool, version string) {
	if interactive {
		clearScreen()
	}
	ui := newTerminalUI(interactive)
	ui.printHeader(version)
	ui.printIntro()
	ui.printStatusGrid(version)
	ui.printExamples()
	ui.printCommandBar()
}

func printSessionHelp(interactive bool) {
	if interactive {
		clearScreen()
	}
	ui := newTerminalUI(interactive)
	ui.printHeader("")
	ui.printPanel("Command palette", []string{
		"/home     return to the dashboard",
		"/help     show this command palette",
		"/doctor   inspect OpenCode runtime and provider auth",
		"/setup    install OpenCode runtime and connect a provider",
		"/update   update gacha to the latest release",
		"/quit     exit",
	})
	ui.printPanel("Example prompts", []string{
		"NVDA 지금 사도 될까?",
		"What should I invest in for the next 6 to 12 months?",
		"I own TSLA. When should I trim, sell, or stop out?",
	})
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

type terminalUI struct {
	interactive bool
	width       int
}

func newTerminalUI(interactive bool) terminalUI {
	return terminalUI{interactive: interactive, width: 88}
}

func (ui terminalUI) printHeader(version string) {
	versionText := " "
	if version != "" {
		versionText = "v" + version
	}
	fmt.Println(ui.color("╭"+repeat("─", ui.width)+"╮", "accent"))
	fmt.Println(ui.color("│", "accent") + "  " + pad("GACHA", ui.width-2-runeLen(versionText)) + versionText + ui.color("│", "accent"))
	fmt.Println(ui.color("│", "accent") + "  " + pad("fresh-data investment research agent", ui.width-2) + ui.color("│", "accent"))
	fmt.Println(ui.color("╰"+repeat("─", ui.width)+"╯", "accent"))
	fmt.Println()
}

func (ui terminalUI) printIntro() {
	ui.printPanel("Ask anything", []string{
		"gacha classifies the request and chooses the right workflow.",
		"It requires current web or market data before analysis.",
		"No fresh data means no investment recommendation.",
	})
}

func (ui terminalUI) printStatusGrid(version string) {
	left := []string{
		"Runtime     " + routeLabel(),
		"Fresh data  required",
		"Trading     disabled",
		"Version     " + version,
	}
	right := []string{
		"1  classify request",
		"2  gather current sources",
		"3  analyze thesis and valuation",
		"4  test risks and exit rules",
	}
	ui.printTwoColumnPanel("Session", left, "Workflow", right)
}

func (ui terminalUI) printExamples() {
	ui.printPanel("Try", []string{
		"NVDA 지금 사도 될까?",
		"What should I invest in for the next 6 to 12 months?",
		"I own TSLA. When should I trim, sell, or stop out?",
	})
}

func (ui terminalUI) printCommandBar() {
	fmt.Println(ui.color("╭─ Commands "+repeat("─", ui.width-11)+"╮", "muted"))
	fmt.Println(ui.color("│", "muted") + "  " + pad("/help   /doctor   /setup   /update   /quit", ui.width-2) + ui.color("│", "muted"))
	fmt.Println(ui.color("╰"+repeat("─", ui.width)+"╯", "muted"))
}

func (ui terminalUI) printPanel(title string, lines []string) {
	fmt.Println(ui.color(topBorder(title, ui.width), "muted"))
	for _, line := range lines {
		fmt.Println(ui.color("│", "muted") + "  " + pad(line, ui.width-2) + ui.color("│", "muted"))
	}
	fmt.Println(ui.color("╰"+repeat("─", ui.width)+"╯", "muted"))
	fmt.Println()
}

func (ui terminalUI) printTwoColumnPanel(leftTitle string, left []string, rightTitle string, right []string) {
	leftWidth := 40
	rightWidth := ui.width - leftWidth - 3
	fmt.Println(ui.color(splitTopBorder(leftTitle, leftWidth, rightTitle, rightWidth), "muted"))
	rows := len(left)
	if len(right) > rows {
		rows = len(right)
	}
	for i := 0; i < rows; i++ {
		leftLine := ""
		rightLine := ""
		if i < len(left) {
			leftLine = left[i]
		}
		if i < len(right) {
			rightLine = right[i]
		}
		fmt.Println(ui.color("│", "muted") + "  " + pad(leftLine, leftWidth-2) + ui.color("│", "muted") + "  " + pad(rightLine, rightWidth-2) + ui.color("│", "muted"))
	}
	fmt.Println(ui.color("╰"+repeat("─", leftWidth)+"┴"+repeat("─", rightWidth)+"╯", "muted"))
	fmt.Println()
}

func (ui terminalUI) color(value string, role string) string {
	if !ui.interactive {
		return value
	}
	switch role {
	case "accent":
		return "\033[38;5;81m" + value + "\033[0m"
	case "muted":
		return "\033[38;5;245m" + value + "\033[0m"
	default:
		return value
	}
}

func printAskPrompt(interactive bool) {
	if interactive {
		fmt.Print("\n\033[38;5;81m╭─\033[0m Ask\n\033[38;5;81m╰─▶\033[0m ")
		return
	}
	fmt.Print("\nAsk > ")
}

func topBorder(title string, width int) string {
	label := "─ " + title + " "
	return "╭" + label + repeat("─", width-runeLen(label)) + "╮"
}

func splitTopBorder(leftTitle string, leftWidth int, rightTitle string, rightWidth int) string {
	leftLabel := "─ " + leftTitle + " "
	rightLabel := "─ " + rightTitle + " "
	return "╭" + leftLabel + repeat("─", leftWidth-runeLen(leftLabel)) + "┬" + rightLabel + repeat("─", rightWidth-runeLen(rightLabel)) + "╮"
}

func repeat(value string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(value, count)
}

func spaces(count int) string {
	return repeat(" ", count)
}

func pad(value string, width int) string {
	current := displayWidth(value)
	if current >= width {
		return truncateDisplay(value, width)
	}
	return value + spaces(width-current)
}

func runeLen(value string) int {
	return len([]rune(value))
}

func displayWidth(value string) int {
	width := 0
	for _, r := range value {
		if r >= 0x1100 {
			width += 2
		} else {
			width++
		}
	}
	return width
}

func truncateDisplay(value string, maxWidth int) string {
	width := 0
	var builder strings.Builder
	for _, r := range value {
		next := 1
		if r >= 0x1100 {
			next = 2
		}
		if width+next > maxWidth {
			break
		}
		builder.WriteRune(r)
		width += next
	}
	return builder.String() + spaces(maxWidth-width)
}
