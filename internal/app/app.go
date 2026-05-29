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
	fmt.Println(`investiq

Usage:
  investiq                                    Open the interactive investiq UI
  iq                                          Open the interactive investiq UI
  iq doctor                                   Check the local AI runtime
  iq setup                                    Install the runtime and connect an AI provider
  iq update                                   Update investiq to the latest release
  iq "question"                               Analyze with automatic request classification

Debug:
  investiq prompt "question"                  Print the composed agent prompt

Examples:
  iq "NVDA 지금 사도 될까?"
  iq "AAPL 현재 매수 구간 분석"
  iq "TSLA 보유 중인데 매도 기준 점검"
  iq "6개월에서 12개월 관점 투자 후보 찾아줘"`)
}
