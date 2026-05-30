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
	if lang, ok := languageFromSetting(os.Getenv("GACHA_LANG")); ok {
		return lang
	}
	if config, err := loadGachaConfig(); err == nil {
		if lang, ok := languageFromSetting(config.Language); ok {
			return lang
		}
	}
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

func languageFromSetting(value string) (language, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case languageSettingKorean, "kr", "kor", "korean", "한국어":
		return languageKorean, true
	case languageSettingEnglish, "eng", "english":
		return languageEnglish, true
	case "", languageSettingAuto:
		return "", false
	default:
		return "", false
	}
}

func normalizeLanguageSetting(value string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", languageSettingAuto:
		return languageSettingAuto, true
	case languageSettingEnglish, "eng", "english":
		return languageSettingEnglish, true
	case languageSettingKorean, "kr", "kor", "korean", "한국어":
		return languageSettingKorean, true
	default:
		return "", false
	}
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
	InputPlaceholder      string
	InputPlaceholderShort string
	Ready                 string
	Auto                  string
	Report                string
	Fallback              string
	Complete              string
	Help                  string
	Command               string
	Runtime               string
	Setup                 string
	Update                string
	System                string
	Researching           string
	ResearchPhases        []string
	Footer                string
	HomeTitle             string
	HomeSubtitle          string
	HomeActionsTitle      string
	HomeActions           []homeAction
	HomeOutcomesTitle     string
	HomeOutcomes          []string
	HomeNote              string
	ContextTitle          string
	ContextRecentTitle    string
	ContextNoRecent       string
	ContextTypesTitle     string
	ContextRequestTitle   string
	ContextResearchTitle  string
	ContextSourcesTitle   string
	ContextSourcesPending string
	ContextReportTitle    string
	ContextReportFallback string
	Onboarding            []string
	Research              func(string) []string
	HelpLines             []string
	SetupLines            []string
	UpdateMessage         string
	ErrorTitle            string
	RuntimeTitle          string
	SettingsTitle         string
	LoginRequired         string
	Missing               string
	RunSetupHint          string
	StatusMode            string
	StatusRuntime         string
	FooterShort           string
	SavePrompt            string
	ReportActionsTitle    string
	ReportActions         []reportChoice
	NewQuestionAction     string
	SavedReport           string
	SkippedSave           string
	SettingsSaved         string
	SettingsInvalidModel  string
	SettingsInvalidLang   string
	SettingsInvalidTheme  string
	ThemeTitle            string
	ThemeIntro            string
	ThemeActive           string
	ThemeCommandHint      string
	ThemePreviewTitle     string
	ThemeSelectLabel      string
	ThemeLabels           map[string]string
	ThemeDescriptions     map[string]string
}

type homeAction struct {
	Name   string
	Prompt string
}

type reportChoice struct {
	Key   string
	Label string
}

func textFor(lang language) uiText {
	if lang == languageKorean {
		return koreanText()
	}
	return englishText()
}

func englishText() uiText {
	return uiText{
		InputPlaceholder:      `Ask: "Should I buy NVDA now?"`,
		InputPlaceholderShort: "Ask a question...",
		Ready:                 "Ready",
		Auto:                  "Auto",
		Report:                "Report",
		Fallback:              "Fallback",
		Complete:              "Complete",
		Help:                  "Help",
		Command:               "Command",
		Runtime:               "Runtime",
		Setup:                 "Setup",
		Update:                "Update",
		System:                "System",
		Researching:           "Researching",
		ResearchPhases: []string{
			"Classifying request",
			"Checking fresh data",
			"Building thesis",
			"Testing valuation",
			"Reviewing risks",
			"Writing report",
		},
		Footer:           " /help  /settings  /theme  /quit   •   enter to run   •   esc to exit",
		FooterShort:      " /help  /quit   •   enter run   •   esc exit",
		HomeTitle:        "What are you deciding?",
		HomeSubtitle:     "Choose a starting point or ask in plain language.",
		HomeActionsTitle: "Decision desk",
		HomeActions: []homeAction{
			{Name: "Buy", Prompt: "Should I buy NVDA now?"},
			{Name: "Find", Prompt: "What should I invest in for the next 6 to 12 months?"},
			{Name: "Hold", Prompt: "I own TSLA. Should I trim, hold, or sell?"},
			{Name: "Exit", Prompt: "Where should I stop out or sell?"},
			{Name: "Portfolio", Prompt: "Is my portfolio too concentrated?"},
		},
		HomeOutcomesTitle: "You'll get",
		HomeOutcomes: []string{
			"Bottom line",
			"Decision rules",
			"Biggest risks",
			"Checked data",
			"Optional detailed analysis",
		},
		HomeNote:              "Fresh data before recommendations. No automatic trading.",
		ContextTitle:          "Context",
		ContextRecentTitle:    "Recent",
		ContextNoRecent:       "No saved reports yet",
		ContextTypesTitle:     "Decision types",
		ContextRequestTitle:   "Current request",
		ContextResearchTitle:  "Research",
		ContextSourcesTitle:   "Sources",
		ContextSourcesPending: "Collected during research",
		ContextReportTitle:    "Report",
		ContextReportFallback: "Report output on the right",
		Onboarding: []string{
			"Setup needed",
			"OpenCode is not installed yet.",
			"Run `gch setup`, connect a provider, then ask your first question.",
			"Provider login needed",
			"OpenCode is installed, but no AI provider is connected yet.",
			"Run `gch setup`, finish provider login, then come back here.",
			"Ready to research",
			"Ask a question below. Gacha will check current data before making a recommendation.",
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
				"4. Run risk review and strongest opposite-view check",
				"5. Produce action conditions and source notes",
				"Waiting for the local AI runtime...",
			}
		},
		HelpLines: []string{
			"Command palette",
			"/home     return to the dashboard",
			"/help     show this command palette",
			"/settings show model and language settings",
			"/model    set model: /model auto, /model opencode-default, or /model provider/model",
			"/language set UI/report language: /language auto, /language en, /language ko",
			"/theme    choose a theme and preview each option",
			"/quit     exit",
		},
		SetupLines: []string{
			"Setup",
			"Use this one-time setup flow:",
			"  gch setup",
			"1. Install OpenCode if it is missing.",
			"2. Connect ChatGPT, Copilot, Gemini, OpenAI API, or another provider.",
			"3. Return here and ask your first investment question.",
		},
		UpdateMessage:      "Run `gacha update` outside the interactive UI to update the binary.",
		ErrorTitle:         "OpenCode failed",
		RuntimeTitle:       "Runtime",
		SettingsTitle:      "Settings",
		LoginRequired:      "login required",
		Missing:            "missing",
		RunSetupHint:       "Run `gch setup` outside this screen to connect ChatGPT, Copilot, Gemini, or an API provider.",
		StatusMode:         "Mode ",
		StatusRuntime:      "Runtime ",
		SavePrompt:         "Next: type d for detailed analysis, y to save, n to skip, or ask a new question.",
		ReportActionsTitle: "Next",
		ReportActions: []reportChoice{
			{Key: "d", Label: "detailed analysis"},
			{Key: "y", Label: "save"},
			{Key: "n", Label: "skip"},
		},
		NewQuestionAction:    "or ask a new question",
		SavedReport:          "Saved report:",
		SkippedSave:          "Report was not saved.",
		SettingsSaved:        "Settings saved.",
		SettingsInvalidModel: "Use `/model auto`, `/model opencode-default`, or `/model provider/model`.",
		SettingsInvalidLang:  "Use `/language auto`, `/language en`, or `/language ko`.",
		SettingsInvalidTheme: "Use `/theme system`, `/theme dark`, `/theme light`, or `/theme gacha`.",
		ThemeTitle:           "Themes",
		ThemeIntro:           "Choose how Gacha should look in your terminal.",
		ThemeActive:          "Active:",
		ThemeCommandHint:     "Use /theme <name> to save and apply a theme.",
		ThemePreviewTitle:    "Preview",
		ThemeSelectLabel:     "select this theme",
		ThemeLabels: map[string]string{
			themeSettingSystem: "System",
			themeSettingDark:   "Dark",
			themeSettingLight:  "Light",
			themeSettingGacha:  "Gacha",
		},
		ThemeDescriptions: map[string]string{
			themeSettingSystem: "Adapts to light and dark terminal backgrounds.",
			themeSettingDark:   "Quiet teal accents for dark terminals.",
			themeSettingLight:  "Deeper accents and darker muted text for light terminals.",
			themeSettingGacha:  "The original fixed palette from earlier Gacha releases.",
		},
	}
}

func koreanText() uiText {
	return uiText{
		InputPlaceholder:      `예: "NVDA 지금 사도 될까?"`,
		InputPlaceholderShort: "질문 입력...",
		Ready:                 "준비됨",
		Auto:                  "자동",
		Report:                "리포트",
		Fallback:              "대체",
		Complete:              "완료",
		Help:                  "도움말",
		Command:               "명령",
		Runtime:               "런타임",
		Setup:                 "설정",
		Update:                "업데이트",
		System:                "시스템",
		Researching:           "조사 중",
		ResearchPhases: []string{
			"요청 분류 중",
			"최신 데이터 확인 중",
			"투자 thesis 구성 중",
			"밸류에이션 점검 중",
			"리스크 검토 중",
			"리포트 작성 중",
		},
		Footer:           " /help  /settings  /theme  /quit   •   enter 실행   •   esc 종료",
		FooterShort:      " /help  /quit   •   enter 실행   •   esc 종료",
		HomeTitle:        "어떤 결정을 도와드릴까요?",
		HomeSubtitle:     "아래에서 시작하거나 평소 말처럼 질문하세요.",
		HomeActionsTitle: "결정 데스크",
		HomeActions: []homeAction{
			{Name: "매수", Prompt: "NVDA 지금 사도 될까?"},
			{Name: "탐색", Prompt: "앞으로 6~12개월 관점에서 무엇에 투자하면 좋을까?"},
			{Name: "보유", Prompt: "TSLA를 보유 중인데 줄일까, 유지할까, 팔까?"},
			{Name: "매도", Prompt: "어디서 손절하거나 매도해야 할까?"},
			{Name: "포트폴리오", Prompt: "내 포트폴리오가 너무 집중되어 있을까?"},
		},
		HomeOutcomesTitle: "받게 되는 답변",
		HomeOutcomes: []string{
			"쉬운 결론",
			"행동 기준",
			"가장 큰 리스크",
			"확인한 데이터",
			"선택 상세 분석",
		},
		HomeNote:              "추천 전 최신 데이터를 확인합니다. 거래는 실행하지 않습니다.",
		ContextTitle:          "맥락",
		ContextRecentTitle:    "최근",
		ContextNoRecent:       "저장된 리포트 없음",
		ContextTypesTitle:     "결정 유형",
		ContextRequestTitle:   "현재 질문",
		ContextResearchTitle:  "조사",
		ContextSourcesTitle:   "출처",
		ContextSourcesPending: "조사 중 수집",
		ContextReportTitle:    "리포트",
		ContextReportFallback: "오른쪽 리포트 참고",
		Onboarding: []string{
			"설정 필요",
			"OpenCode가 아직 설치되어 있지 않습니다.",
			"`gch setup`을 실행해 provider를 연결한 뒤 첫 질문을 입력하세요.",
			"provider 로그인 필요",
			"OpenCode는 설치되어 있지만 AI provider가 아직 연결되지 않았습니다.",
			"`gch setup`에서 provider 로그인을 마친 뒤 돌아오세요.",
			"리서치 준비 완료",
			"아래에 질문을 입력하세요. Gacha는 추천 전에 최신 데이터를 확인합니다.",
		},
		Research: func(query string) []string {
			return []string{
				"리서치 실행",
				"질문:",
				"  " + query,
				"진행 단계",
				"1. 요청 분류: discover, select, entry, exit, portfolio, journal",
				"2. 최신 웹 또는 시장 데이터 요구",
				"3. 쉬운 기본 리포트와 필요한 상세 분석 구성",
				"4. 리스크 검토와 반대 논리 점검",
				"5. 행동 조건과 출처 정리",
				"로컬 AI 런타임을 기다리는 중...",
			}
		},
		HelpLines: []string{
			"명령 팔레트",
			"/home     대시보드로 돌아가기",
			"/help     명령 팔레트 보기",
			"/settings 모델과 언어 설정 보기",
			"/model    모델 설정: /model auto, /model opencode-default, /model provider/model",
			"/language UI/리포트 언어: /language auto, /language en, /language ko",
			"/theme    테마 선택과 예시 보기",
			"/quit     종료",
		},
		SetupLines: []string{
			"설정",
			"처음 한 번 다음 설정 흐름을 실행하세요:",
			"  gch setup",
			"1. OpenCode가 없으면 설치합니다.",
			"2. ChatGPT, Copilot, Gemini, OpenAI API 또는 다른 provider를 연결합니다.",
			"3. 다시 돌아와 첫 투자 질문을 입력합니다.",
		},
		UpdateMessage:      "바이너리를 업데이트하려면 인터랙티브 UI 밖에서 `gacha update`를 실행하세요.",
		ErrorTitle:         "OpenCode 실행 실패",
		RuntimeTitle:       "런타임",
		SettingsTitle:      "설정",
		LoginRequired:      "로그인 필요",
		Missing:            "없음",
		RunSetupHint:       "ChatGPT, Copilot, Gemini 또는 API provider를 연결하려면 이 화면 밖에서 `gch setup`을 실행하세요.",
		StatusMode:         "모드 ",
		StatusRuntime:      "런타임 ",
		SavePrompt:         "다음: d=상세 분석, y=저장, n=건너뛰기, 또는 새 질문을 입력하세요.",
		ReportActionsTitle: "다음",
		ReportActions: []reportChoice{
			{Key: "d", Label: "상세 분석"},
			{Key: "y", Label: "저장"},
			{Key: "n", Label: "건너뛰기"},
		},
		NewQuestionAction:    "또는 새 질문 입력",
		SavedReport:          "리포트 저장:",
		SkippedSave:          "리포트를 저장하지 않았습니다.",
		SettingsSaved:        "설정을 저장했습니다.",
		SettingsInvalidModel: "`/model auto`, `/model opencode-default`, 또는 `/model provider/model` 형식으로 입력하세요.",
		SettingsInvalidLang:  "`/language auto`, `/language en`, 또는 `/language ko`를 입력하세요.",
		SettingsInvalidTheme: "`/theme system`, `/theme dark`, `/theme light`, 또는 `/theme gacha`를 입력하세요.",
		ThemeTitle:           "테마",
		ThemeIntro:           "터미널에 맞는 Gacha 화면 스타일을 선택하세요.",
		ThemeActive:          "현재:",
		ThemeCommandHint:     "/theme <이름>을 입력하면 저장 후 바로 적용됩니다.",
		ThemePreviewTitle:    "예시",
		ThemeSelectLabel:     "이 테마 선택",
		ThemeLabels: map[string]string{
			themeSettingSystem: "시스템",
			themeSettingDark:   "다크",
			themeSettingLight:  "라이트",
			themeSettingGacha:  "Gacha",
		},
		ThemeDescriptions: map[string]string{
			themeSettingSystem: "터미널의 밝은/어두운 배경에 맞춰 조정됩니다.",
			themeSettingDark:   "어두운 터미널에 맞춘 차분한 틸 계열 강조색입니다.",
			themeSettingLight:  "밝은 터미널용으로 강조색과 보조 텍스트 대비를 높입니다.",
			themeSettingGacha:  "이전 Gacha 릴리스의 고정 팔레트입니다.",
		},
	}
}
