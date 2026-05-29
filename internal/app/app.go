package app

import "fmt"

type App struct {
	version string
}

func New(version string) *App {
	return &App{version: version}
}

func (a *App) Run(args []string) error {
	if len(args) == 0 {
		return a.startSession()
	}
	if args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		printUsage()
		return nil
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
	case "prompt":
		return printPrompt(args[1:])
	default:
		return runQuery(args)
	}
}

func printUsage() {
	fmt.Println(`gacha

Usage:
  gacha                                    Open the interactive gacha UI
  gch                                      Open the interactive gacha UI
  gch doctor                               Check the local AI runtime
  gch setup                                Install the runtime and connect an AI provider
  gch update                               Update gacha to the latest release
  gch "question"                           Analyze with automatic request classification

Debug:
  gacha prompt "question"                  Print the composed agent prompt

Examples:
  gch "NVDA 지금 사도 될까?"
  gch "AAPL 현재 매수 구간 분석"
  gch "TSLA 보유 중인데 매도 기준 점검"
  gch "6개월에서 12개월 관점 투자 후보 찾아줘"`)
}
