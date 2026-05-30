package app

import (
	"fmt"
	"strings"
)

func doctor() error {
	status := "missing"
	if hasRunnableCommand(openCodeCommand) {
		status = "ready"
		if !hasOpenCodeAuth() {
			status = "login?"
		}
	}

	fmt.Printf("%-9s %-8s %s\n", "opencode", status, "OpenCode runtime")
	fmt.Printf("           command: %s\n", openCodeCommand)
	if resolved, ok := resolveCommand(openCodeCommand); ok && resolved != openCodeCommand {
		fmt.Printf("           resolved: %s\n", resolved)
	}
	fmt.Printf("           auth: %s\n", openCodeAuthPath())
	resolution := resolveOpenCodeModel()
	fmt.Printf("           model: %s\n", modelDescription(resolution))
	if resolution.Source != "" {
		fmt.Printf("           model source: %s\n", resolution.Source)
	}
	if hasOpenCodeAuth() {
		if providers, err := openCodeAuthList(); err == nil && strings.TrimSpace(providers) != "" {
			for _, line := range strings.Split(strings.TrimSpace(providers), "\n") {
				fmt.Printf("           provider: %s\n", strings.TrimSpace(line))
			}
		}
	} else {
		fmt.Println("           next: run `gch setup` to connect ChatGPT, Copilot, Gemini, or an API provider")
	}

	fmt.Printf("%-9s %-8s %s\n", "manual", "ready", "Copy/paste prompt fallback")
	return nil
}

func routeLabel() string {
	if hasRunnableCommand(openCodeCommand) && hasOpenCodeAuth() {
		return "OpenCode runtime"
	}
	return "Copy/paste prompt"
}

func routeLabelFor(lang language) string {
	if lang != languageKorean {
		return routeLabel()
	}
	if hasRunnableCommand(openCodeCommand) && hasOpenCodeAuth() {
		return "OpenCode 런타임"
	}
	return "복사/붙여넣기 프롬프트"
}
