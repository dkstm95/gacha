package app

import "fmt"

var modes = map[string]bool{
	"auto":      true,
	"discover":  true,
	"select":    true,
	"entry":     true,
	"exit":      true,
	"portfolio": true,
	"journal":   true,
}

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
	case "init":
		return initConfig(hasFlag(args[1:], "--yes") || hasFlag(args[1:], "-y"))
	case "doctor":
		return doctor()
	case "update":
		return a.updateSelf()
	case "platforms":
		return platforms()
	case "prompt":
		if len(args) >= 2 && modes[args[1]] {
			return printPrompt(args[1], args[2:])
		}
		return printPrompt("auto", args[1:])
	case "run":
		if len(args) >= 2 && modes[args[1]] {
			return runMode(args[1], args[2:])
		}
		return runMode("auto", args[1:])
	default:
		if modes[args[0]] {
			return runMode(args[0], args[1:])
		}
		return runMode("auto", args)
	}
}

func printUsage() {
	fmt.Println(`investiq

Usage:
  investiq                                    Open the interactive investiq UI
  iq                                          Open the interactive investiq UI
  iq init                                     Set up AI platform routing
  iq doctor                                   Check detected AI platforms
  iq update                                   Update investiq to the latest release
  iq "question"                               Analyze with automatic request classification

Debug:
  investiq prompt "question"                  Print the composed agent prompt
  investiq platforms                          Print platform config

Examples:
  iq "NVDA 지금 사도 될까?"
  iq "AAPL 현재 매수 구간 분석"
  iq "TSLA 보유 중인데 매도 기준 점검"
  iq "6개월에서 12개월 관점 투자 후보 찾아줘"`)
}
