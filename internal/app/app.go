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
  gacha doctor                             Check the local AI runtime
  gacha setup                              Install the runtime and connect an AI provider
  gacha update                             Update gacha to the latest release
  gacha "question"                         Analyze with automatic request classification

Debug:
  gacha prompt "question"                  Print the composed agent prompt

Examples:
  gacha "NVDA 지금 사도 될까?"
  gacha "AAPL 현재 매수 구간 분석"
  gacha "TSLA 보유 중인데 매도 기준 점검"
  gacha "6개월에서 12개월 관점 투자 후보 찾아줘"`)
}
