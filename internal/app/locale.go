package app

import (
	"os"
	"strings"
	"unicode"
)

type language string

const (
	languageEnglish language = "English"
	languageKorean  language = "Korean"
)

func detectLanguage() language {
	return detectLanguageFromEnv(os.Getenv)
}

func detectLanguageFromEnv(getenv func(string) string) language {
	for _, key := range []string{"GACHA_LANG", "LANGUAGE", "LC_ALL", "LC_MESSAGES", "LANG"} {
		value := strings.ToLower(strings.TrimSpace(getenv(key)))
		if value == "" {
			continue
		}
		if strings.HasPrefix(value, "ko") || strings.Contains(value, ":ko") {
			return languageKorean
		}
		if strings.HasPrefix(value, "en") || strings.Contains(value, ":en") {
			return languageEnglish
		}
	}
	return languageEnglish
}

func responseLanguage(query string) language {
	if containsKorean(query) {
		return languageKorean
	}
	return detectLanguage()
}

func containsKorean(value string) bool {
	for _, r := range value {
		if unicode.In(r, unicode.Hangul) {
			return true
		}
	}
	return false
}

type uiText struct {
	InputPlaceholder string
	Ready            string
	Auto             string
	Report           string
	Fallback         string
	Complete         string
	Help             string
	Command          string
	Runtime          string
	Setup            string
	Update           string
	System           string
	Researching      string
	ResearchPhases   []string
	Footer           string
	Welcome          []string
	Research         func(string) []string
	HelpLines        []string
	SetupLines       []string
	UpdateMessage    string
	ErrorTitle       string
	RuntimeTitle     string
	LoginRequired    string
	Missing          string
	RunSetupHint     string
	StatusMode       string
	StatusRuntime    string
	StatusFreshData  string
	StatusNoTrading  string
	SavePrompt       string
	SavedReport      string
	SkippedSave      string
}

func textFor(lang language) uiText {
	if lang == languageKorean {
		return koreanText()
	}
	return englishText()
}

func englishText() uiText {
	return uiText{
		InputPlaceholder: "Ask about an investment...",
		Ready:            "Ready",
		Auto:             "Auto",
		Report:           "Report",
		Fallback:         "Fallback",
		Complete:         "Complete",
		Help:             "Help",
		Command:          "Command",
		Runtime:          "Runtime",
		Setup:            "Setup",
		Update:           "Update",
		System:           "System",
		Researching:      "Researching",
		ResearchPhases: []string{
			"Classifying request",
			"Checking fresh data",
			"Building thesis",
			"Testing valuation",
			"Reviewing risks",
			"Writing report",
		},
		Footer: " /help  /doctor  /setup  /update  /quit   •   enter to run   •   esc to exit",
		Welcome: []string{
			"Investment research cockpit",
			"Ask one question. Gacha routes it through the right research workflow.",
			"Every workflow requires current web or market data before analysis.",
			"Workflow rail",
			"Discover|find opportunities when you do not know what to buy",
			"Select|rank concrete assets inside a sector or theme",
			"Entry|decide whether the current price is attractive",
			"Exit|define trim, sell, stop-loss, and thesis-break zones",
			"Portfolio|review concentration, exposure, and rebalancing risks",
			"Journal|record thesis, decision rules, and postmortems",
			"Try",
			"NVDA 지금 사도 될까?",
			"What should I invest in for the next 6 to 12 months?",
			"I own TSLA. When should I trim, sell, or stop out?",
			"Report contract",
			"Data freshness • Sources • Thesis • Valuation • Risks • Devil's Advocate • Action conditions",
			"No fresh data, no recommendation. Trading is disabled.",
		},
		Research: func(query string) []string {
			return []string{
				"Research run",
				"Query:",
				"  " + query,
				"Pipeline",
				"1. Classify request: discover, select, entry, exit, portfolio, or journal",
				"2. Require current web or market data",
				"3. Build thesis, valuation, and scenario analysis",
				"4. Run risk review and Devil's Advocate",
				"5. Produce action conditions and provenance",
				"Waiting for the local AI runtime...",
			}
		},
		HelpLines: []string{
			"Command palette",
			"/home     return to the dashboard",
			"/help     show this command palette",
			"/doctor   inspect OpenCode runtime and provider auth",
			"/setup    show setup instructions",
			"/update   show update instructions",
			"/quit     exit",
		},
		SetupLines: []string{
			"Setup",
			"Run this command in your shell:",
			"  gch setup",
			"That flow installs OpenCode if needed and starts provider login.",
			"Interactive provider login is intentionally handled outside this screen so your terminal can hand control to OpenCode safely.",
		},
		UpdateMessage:   "Run `gacha update` outside the interactive UI to update the binary.",
		ErrorTitle:      "OpenCode failed",
		RuntimeTitle:    "Runtime",
		LoginRequired:   "login required",
		Missing:         "missing",
		RunSetupHint:    "Run `gch setup` outside this screen to connect ChatGPT, Copilot, Gemini, or an API provider.",
		StatusMode:      "Mode ",
		StatusRuntime:   "Runtime ",
		StatusFreshData: "Fresh data required",
		StatusNoTrading: "No trading",
		SavePrompt:      "Save this report as Markdown? Type y or n.",
		SavedReport:     "Saved report:",
		SkippedSave:     "Report was not saved.",
	}
}

func koreanText() uiText {
	return uiText{
		InputPlaceholder: "투자 질문을 입력하세요...",
		Ready:            "준비됨",
		Auto:             "자동",
		Report:           "리포트",
		Fallback:         "대체",
		Complete:         "완료",
		Help:             "도움말",
		Command:          "명령",
		Runtime:          "런타임",
		Setup:            "설정",
		Update:           "업데이트",
		System:           "시스템",
		Researching:      "조사 중",
		ResearchPhases: []string{
			"요청 분류 중",
			"최신 데이터 확인 중",
			"투자 thesis 구성 중",
			"밸류에이션 점검 중",
			"리스크 검토 중",
			"리포트 작성 중",
		},
		Footer: " /help  /doctor  /setup  /update  /quit   •   enter 실행   •   esc 종료",
		Welcome: []string{
			"투자 리서치 조종석",
			"질문 하나를 입력하면 Gacha가 알맞은 리서치 흐름으로 보냅니다.",
			"모든 흐름은 분석 전에 최신 웹 또는 시장 데이터를 요구합니다.",
			"워크플로우",
			"Discover|무엇을 살지 모를 때 투자 기회 찾기",
			"Select|섹터나 테마 안에서 구체적 후보 비교",
			"Entry|현재 가격이 매수하기 좋은지 판단",
			"Exit|분할매도, 매도, 손절, thesis break 구간 정의",
			"Portfolio|집중도, 노출, 리밸런싱 리스크 점검",
			"Journal|투자 thesis, 의사결정 기준, 회고 기록",
			"예시",
			"NVDA 지금 사도 될까?",
			"앞으로 6~12개월 관점에서 무엇에 투자하면 좋을까?",
			"TSLA를 보유 중인데 언제 줄이거나 팔아야 할까?",
			"리포트 원칙",
			"데이터 시점 • 출처 • thesis • 밸류에이션 • 리스크 • 반대 논리 • 행동 조건",
			"최신 데이터가 없으면 추천도 없습니다. Gacha는 거래를 실행하지 않습니다.",
		},
		Research: func(query string) []string {
			return []string{
				"리서치 실행",
				"질문:",
				"  " + query,
				"진행 단계",
				"1. 요청 분류: discover, select, entry, exit, portfolio, journal",
				"2. 최신 웹 또는 시장 데이터 요구",
				"3. thesis, 밸류에이션, 시나리오 분석 구성",
				"4. 리스크 검토와 반대 논리 점검",
				"5. 행동 조건과 출처 정리",
				"로컬 AI 런타임을 기다리는 중...",
			}
		},
		HelpLines: []string{
			"명령 팔레트",
			"/home     대시보드로 돌아가기",
			"/help     명령 팔레트 보기",
			"/doctor   OpenCode 런타임과 provider 인증 점검",
			"/setup    설정 안내 보기",
			"/update   업데이트 안내 보기",
			"/quit     종료",
		},
		SetupLines: []string{
			"설정",
			"셸에서 다음 명령을 실행하세요:",
			"  gch setup",
			"필요하면 OpenCode를 설치하고 provider 로그인을 시작합니다.",
			"provider 로그인은 터미널 제어권을 OpenCode에 안전하게 넘기기 위해 이 화면 밖에서 처리합니다.",
		},
		UpdateMessage:   "바이너리를 업데이트하려면 인터랙티브 UI 밖에서 `gacha update`를 실행하세요.",
		ErrorTitle:      "OpenCode 실행 실패",
		RuntimeTitle:    "런타임",
		LoginRequired:   "로그인 필요",
		Missing:         "없음",
		RunSetupHint:    "ChatGPT, Copilot, Gemini 또는 API provider를 연결하려면 이 화면 밖에서 `gch setup`을 실행하세요.",
		StatusMode:      "모드 ",
		StatusRuntime:   "런타임 ",
		StatusFreshData: "최신 데이터 필수",
		StatusNoTrading: "거래 실행 안 함",
		SavePrompt:      "이 리포트를 Markdown으로 저장할까요? y 또는 n을 입력하세요.",
		SavedReport:     "리포트 저장:",
		SkippedSave:     "리포트를 저장하지 않았습니다.",
	}
}
