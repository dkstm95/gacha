package app

import (
	"fmt"
	"strings"
)

type App struct {
	version string
	env     Env
}

func New(version string) *App {
	return &App{version: version, env: defaultEnv()}
}

func (a *App) Run(args []string) error {
	if len(args) == 0 {
		return a.startSession()
	}
	if args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		printUsage()
		return nil
	}
	if isSettingsCommand(args[0]) {
		fmt.Println(settingsContent(textFor(detectLanguage())))
		return nil
	}
	if isSlashCommand(args[0]) {
		return fmt.Errorf(textFor(detectLanguage()).UnknownCommand, args[0])
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Println(a.version)
		return nil
	case "doctor":
		return doctor()
	case "setup":
		return ensureRuntime(true)
	case "update":
		return a.updateSelf()
	case "config":
		return runConfigCommand(args[1:])
	case "prompt":
		return printPrompt(args[1:])
	default:
		return a.runQuery(args)
	}
}

func isSettingsCommand(value string) bool {
	switch value {
	case "/settings", "settings":
		return true
	default:
		return false
	}
}

func isSlashCommand(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), "/")
}

func printUsage() {
	fmt.Println(`gacha

Usage:
  gacha                                    Open the interactive gacha UI
  gch                                      Open the interactive gacha UI
  gch doctor                               Check the local AI runtime
  gch setup                                Install the runtime and connect an AI provider
  gch settings                             Show model and language settings
  gch config get                           Print the current config JSON
  gch config set model auto                Set model mode or provider/model
  gch config set language ko               Set language: auto, en, or ko
  gch config set theme system              Set theme: system, dark, light, or gacha
  gch update                               Update gacha to the latest release
  gch "question"                           Analyze with automatic request classification

Examples:
  gch "NVDA 지금 사도 될까?"
  gch "AAPL 현재 매수 구간 분석"
  gch "TSLA 보유 중인데 매도 기준 점검"
  gch "6개월에서 12개월 관점 투자 후보 찾아줘"`)
}
