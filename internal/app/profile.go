package app

import (
	"fmt"
	"io"
	"strings"
)

const (
	profileMarketUSStocksETFs     = "us-stocks-etfs"
	profileMarketKoreanStocksETFs = "korean-stocks-etfs"
	profileMarketGlobalExUS       = "global-ex-us"
	profileMarketCrypto           = "crypto"
	profileMarketBondsCash        = "bonds-cash"
	profileValueNotSure           = "not-sure"

	profileHorizonOneToThreeMonths = "1-3m"
	profileHorizonSixToTwelve      = "6-12m"
	profileHorizonOneToThreeYears  = "1-3y"
	profileHorizonFiveYearsPlus    = "5y-plus"

	profileRiskConservative = "conservative"
	profileRiskBalanced     = "balanced"
	profileRiskAggressive   = "aggressive"

	profileReportBasicFirst = "basic-first"
	profileReportDetailed   = "detailed"
	profileReportConcise    = "concise"

	profileGoalDiscover  = "discover"
	profileGoalTheme     = "theme"
	profileGoalEntry     = "entry"
	profileGoalHolding   = "holding"
	profileGoalPortfolio = "portfolio"
)

type gachaProfile struct {
	Markets     profileMulti      `json:"markets,omitempty"`
	Horizons    profileMulti      `json:"horizons,omitempty"`
	Risk        string            `json:"risk,omitempty"`
	ReportStyle string            `json:"reportStyle,omitempty"`
	Goals       profileMulti      `json:"goals,omitempty"`
	Onboarding  profileOnboarding `json:"onboarding,omitempty"`
}

type profileMulti struct {
	Selected []string `json:"selected,omitempty"`
	Default  string   `json:"default,omitempty"`
}

type profileOnboarding struct {
	Completed bool `json:"completed,omitempty"`
	Skipped   bool `json:"skipped,omitempty"`
}

type profileOption struct {
	Value string
	EN    string
	KO    string
}

type profileCategory string

const (
	profileCategoryMarkets     profileCategory = "markets"
	profileCategoryHorizons    profileCategory = "horizons"
	profileCategoryRisk        profileCategory = "risk"
	profileCategoryReportStyle profileCategory = "report-style"
	profileCategoryGoals       profileCategory = "goals"
)

func marketOptions() []profileOption {
	return []profileOption{
		{Value: profileMarketUSStocksETFs, EN: "US stocks / ETFs", KO: "미국 주식/ETF"},
		{Value: profileMarketKoreanStocksETFs, EN: "Korean stocks / ETFs", KO: "한국 주식/ETF"},
		{Value: profileMarketGlobalExUS, EN: "Global ex-US", KO: "미국 외 글로벌"},
		{Value: profileMarketCrypto, EN: "Crypto", KO: "가상자산"},
		{Value: profileMarketBondsCash, EN: "Bonds / cash-like", KO: "채권/현금성"},
		{Value: profileValueNotSure, EN: "Not sure", KO: "잘 모르겠음"},
	}
}

func horizonOptions() []profileOption {
	return []profileOption{
		{Value: profileHorizonOneToThreeMonths, EN: "1-3 months", KO: "1-3개월"},
		{Value: profileHorizonSixToTwelve, EN: "6-12 months", KO: "6-12개월"},
		{Value: profileHorizonOneToThreeYears, EN: "1-3 years", KO: "1-3년"},
		{Value: profileHorizonFiveYearsPlus, EN: "5+ years", KO: "5년 이상"},
		{Value: profileValueNotSure, EN: "Not sure", KO: "잘 모르겠음"},
	}
}

func riskOptions() []profileOption {
	return []profileOption{
		{Value: profileRiskConservative, EN: "Conservative", KO: "보수형"},
		{Value: profileRiskBalanced, EN: "Balanced", KO: "균형형"},
		{Value: profileRiskAggressive, EN: "Aggressive", KO: "공격형"},
		{Value: profileValueNotSure, EN: "Not sure", KO: "잘 모르겠음"},
	}
}

func reportStyleOptions() []profileOption {
	return []profileOption{
		{Value: profileReportBasicFirst, EN: "Basic first, details on request", KO: "기본 리포트 먼저, 상세 분석은 요청 시"},
		{Value: profileReportDetailed, EN: "Detailed every time", KO: "항상 상세하게"},
		{Value: profileReportConcise, EN: "Very concise", KO: "아주 간결하게"},
	}
}

func goalOptions() []profileOption {
	return []profileOption{
		{Value: profileGoalDiscover, EN: "Discover what to invest in", KO: "무엇에 투자할지 탐색"},
		{Value: profileGoalTheme, EN: "Compare a theme or sector", KO: "테마나 섹터 비교"},
		{Value: profileGoalEntry, EN: "Plan entry timing", KO: "매수 진입 계획"},
		{Value: profileGoalHolding, EN: "Manage existing holdings", KO: "보유 종목 관리"},
		{Value: profileGoalPortfolio, EN: "Review portfolio", KO: "포트폴리오 점검"},
	}
}

func profileCategoryTitle(category profileCategory, lang language) string {
	if lang == languageKorean {
		switch category {
		case profileCategoryMarkets:
			return "관심 시장"
		case profileCategoryHorizons:
			return "투자 기간"
		case profileCategoryRisk:
			return "리스크 성향"
		case profileCategoryReportStyle:
			return "리포트 스타일"
		case profileCategoryGoals:
			return "자주 하는 판단"
		default:
			return "투자 프로필"
		}
	}
	switch category {
	case profileCategoryMarkets:
		return "Primary markets"
	case profileCategoryHorizons:
		return "Usual time horizons"
	case profileCategoryRisk:
		return "Risk preference"
	case profileCategoryReportStyle:
		return "Report style"
	case profileCategoryGoals:
		return "Common goals"
	default:
		return "Research Profile"
	}
}

func profileOptionsFor(category profileCategory) []profileOption {
	switch category {
	case profileCategoryMarkets:
		return marketOptions()
	case profileCategoryHorizons:
		return horizonOptions()
	case profileCategoryRisk:
		return riskOptions()
	case profileCategoryReportStyle:
		return reportStyleOptions()
	case profileCategoryGoals:
		return goalOptions()
	default:
		return nil
	}
}

func profileLabel(value string, options []profileOption, lang language) string {
	for _, option := range options {
		if option.Value == value {
			if lang == languageKorean {
				return option.KO
			}
			return option.EN
		}
	}
	return value
}

func profileMultiLabels(value profileMulti, options []profileOption, lang language) []string {
	labels := make([]string, 0, len(value.Selected))
	for _, selected := range value.Selected {
		labels = append(labels, profileLabel(selected, options, lang))
	}
	return labels
}

func profileHasValues(profile gachaProfile) bool {
	return len(profile.Markets.Selected) > 0 ||
		len(profile.Horizons.Selected) > 0 ||
		strings.TrimSpace(profile.Risk) != "" ||
		strings.TrimSpace(profile.ReportStyle) != "" ||
		len(profile.Goals.Selected) > 0
}

func profileIsZero(profile gachaProfile) bool {
	return !profileHasValues(profile) &&
		!profile.Onboarding.Completed &&
		!profile.Onboarding.Skipped
}

func shouldShowProfileOnboarding(profile gachaProfile) bool {
	return !profile.Onboarding.Completed && !profile.Onboarding.Skipped
}

func normalizeProfile(profile gachaProfile) gachaProfile {
	profile.Markets = normalizeProfileMulti(profile.Markets, marketOptions(), true)
	profile.Horizons = normalizeProfileMulti(profile.Horizons, horizonOptions(), true)
	profile.Goals = normalizeProfileMulti(profile.Goals, goalOptions(), false)
	if !validProfileOption(profile.Risk, riskOptions()) {
		profile.Risk = ""
	}
	if !validProfileOption(profile.ReportStyle, reportStyleOptions()) {
		profile.ReportStyle = ""
	}
	return profile
}

func normalizeProfileMulti(value profileMulti, options []profileOption, notSureExclusive bool) profileMulti {
	allowed := map[string]bool{}
	for _, option := range options {
		allowed[option.Value] = true
	}
	var selected []string
	seen := map[string]bool{}
	for _, item := range value.Selected {
		item = strings.TrimSpace(item)
		if !allowed[item] || seen[item] {
			continue
		}
		selected = append(selected, item)
		seen[item] = true
	}
	if notSureExclusive && seen[profileValueNotSure] {
		selected = []string{profileValueNotSure}
		seen = map[string]bool{profileValueNotSure: true}
	}
	value.Selected = selected
	if !seen[value.Default] || value.Default == profileValueNotSure {
		if len(selected) > 0 && selected[0] != profileValueNotSure {
			value.Default = selected[0]
		} else {
			value.Default = ""
		}
	}
	return value
}

func validProfileOption(value string, options []profileOption) bool {
	if strings.TrimSpace(value) == "" {
		return true
	}
	for _, option := range options {
		if option.Value == value {
			return true
		}
	}
	return false
}

func updateConfigProfile(profile gachaProfile) error {
	config, err := loadGachaConfig()
	if err != nil {
		return err
	}
	config.Profile = normalizeProfile(profile)
	return saveGachaConfig(config)
}

func resetConfigProfile() error {
	config, err := loadGachaConfig()
	if err != nil {
		return err
	}
	config.Profile = gachaProfile{}
	return saveGachaConfig(config)
}

func printProfile() error {
	return printProfileTo(defaultEnv().Stdout)
}

func printProfileTo(writer io.Writer) error {
	config, err := configWithDefaults()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, stripANSI(profileDetails(config.Profile, detectLanguage())))
	return err
}

func profileSummary(profile gachaProfile, lang language) string {
	profile = normalizeProfile(profile)
	if !profileHasValues(profile) {
		if lang == languageKorean {
			return "투자 프로필 없음"
		}
		return "No research profile"
	}
	var parts []string
	if labels := profileMultiLabels(profile.Markets, marketOptions(), lang); len(labels) > 0 {
		parts = append(parts, strings.Join(labels, ", "))
	}
	if profile.Horizons.Default != "" {
		parts = append(parts, profileLabel(profile.Horizons.Default, horizonOptions(), lang))
	} else if labels := profileMultiLabels(profile.Horizons, horizonOptions(), lang); len(labels) > 0 {
		parts = append(parts, strings.Join(labels, ", "))
	}
	if profile.Risk != "" {
		parts = append(parts, profileLabel(profile.Risk, riskOptions(), lang))
	}
	if len(parts) == 0 {
		if lang == languageKorean {
			return "투자 프로필 설정됨"
		}
		return "Research profile set"
	}
	return strings.Join(parts, " · ")
}

func profileTitleForLang(lang language) string {
	if lang == languageKorean {
		return "투자 프로필"
	}
	return "Research Profile"
}

func profileDetails(profile gachaProfile, lang language) string {
	profile = normalizeProfile(profile)
	title := "Research Profile"
	if lang == languageKorean {
		title = "투자 프로필"
	}
	lines := []string{titleStyle.Render(title), ""}
	lines = append(lines, profileDetailLine("Markets", "관심 시장", strings.Join(profileMultiLabels(profile.Markets, marketOptions(), lang), ", "), lang))
	lines = append(lines, profileDetailLine("Default market", "기본 시장", profileLabel(profile.Markets.Default, marketOptions(), lang), lang))
	lines = append(lines, profileDetailLine("Horizons", "투자 기간", strings.Join(profileMultiLabels(profile.Horizons, horizonOptions(), lang), ", "), lang))
	lines = append(lines, profileDetailLine("Default horizon", "기본 기간", profileLabel(profile.Horizons.Default, horizonOptions(), lang), lang))
	lines = append(lines, profileDetailLine("Risk", "리스크 성향", profileLabel(profile.Risk, riskOptions(), lang), lang))
	lines = append(lines, profileDetailLine("Report style", "리포트 스타일", profileLabel(profile.ReportStyle, reportStyleOptions(), lang), lang))
	lines = append(lines, profileDetailLine("Goals", "자주 하는 판단", strings.Join(profileMultiLabels(profile.Goals, goalOptions(), lang), ", "), lang))
	return strings.Join(lines, "\n")
}

func profileDetailLine(enLabel, koLabel, value string, lang language) string {
	if strings.TrimSpace(value) == "" {
		if lang == languageKorean {
			value = "미설정"
		} else {
			value = "Not set"
		}
	}
	label := enLabel
	if lang == languageKorean {
		label = koLabel
	}
	return fmt.Sprintf("%-16s %s", label, value)
}

func profilePromptBlock(config gachaConfig, lang language) string {
	profile := normalizeProfile(config.Profile)
	if config.Profile.Onboarding.Skipped || !profileHasValues(profile) {
		return ""
	}
	lines := []string{"User research profile:"}
	if labels := profileMultiLabels(profile.Markets, marketOptions(), lang); len(labels) > 0 {
		lines = append(lines, "- Markets of interest: "+strings.Join(labels, ", "))
	}
	if profile.Markets.Default != "" {
		lines = append(lines, "- Default market when unspecified: "+profileLabel(profile.Markets.Default, marketOptions(), lang))
	}
	if labels := profileMultiLabels(profile.Horizons, horizonOptions(), lang); len(labels) > 0 {
		lines = append(lines, "- Time horizons of interest: "+strings.Join(labels, ", "))
	}
	if profile.Horizons.Default != "" {
		lines = append(lines, "- Default horizon when unspecified: "+profileLabel(profile.Horizons.Default, horizonOptions(), lang))
	}
	if profile.Risk != "" {
		lines = append(lines, "- Risk preference: "+profileLabel(profile.Risk, riskOptions(), lang))
	}
	if labels := profileMultiLabels(profile.Goals, goalOptions(), lang); len(labels) > 0 {
		lines = append(lines, "- Common research goals: "+strings.Join(labels, ", "))
	}
	if profile.ReportStyle != "" {
		lines = append(lines, "- Default report style: "+profileLabel(profile.ReportStyle, reportStyleOptions(), lang))
	}
	lines = append(lines, "", "If the user's question specifies a different market, horizon, goal, risk preference, or report style, prefer the user's current question over the saved profile.")
	return strings.Join(lines, "\n")
}
