package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type profileFlowMode string

const (
	profileFlowOnboarding profileFlowMode = "onboarding"
	profileFlowMenu       profileFlowMode = "menu"
	profileFlowEdit       profileFlowMode = "edit"
)

type profileFlow struct {
	Mode     profileFlowMode
	Draft    gachaProfile
	Step     int
	Category profileCategory
	Cursor   int
	Message  string
}

var profileOnboardingSteps = []profileCategory{
	profileCategoryMarkets,
	profileCategoryHorizons,
	profileCategoryRisk,
	profileCategoryReportStyle,
	profileCategoryGoals,
}

var profileMenuCategories = []profileCategory{
	profileCategoryMarkets,
	profileCategoryHorizons,
	profileCategoryRisk,
	profileCategoryReportStyle,
	profileCategoryGoals,
}

func newProfileOnboarding(profile gachaProfile) *profileFlow {
	draft := normalizeProfile(profile)
	return &profileFlow{
		Mode:     profileFlowOnboarding,
		Draft:    draft,
		Category: profileOnboardingSteps[0],
		Cursor:   profileInitialCursor(profileOnboardingSteps[0], draft),
	}
}

func newProfileMenu(profile gachaProfile) *profileFlow {
	return &profileFlow{
		Mode:  profileFlowMenu,
		Draft: normalizeProfile(profile),
	}
}

func (f *profileFlow) render(lang language, width int) string {
	switch f.Mode {
	case profileFlowMenu:
		return f.renderMenu(lang, width)
	default:
		return f.renderEdit(lang, width)
	}
}

func (m tuiModel) handleProfileKey(key string) (tea.Model, tea.Cmd) {
	if m.profile == nil {
		return m, nil
	}
	switch key {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.profile.Mode == profileFlowOnboarding {
			return m, tea.Quit
		}
		m.profile = nil
		m.status = m.text.Ready
		m.mode = m.text.Auto
		m.view.SetContent(welcomeContent(m.version, m.text, m.view.Width, m.view.Height))
		m.view.GotoTop()
		return m, nil
	case "up", "k":
		m.profile.move(-1)
	case "down", "j":
		m.profile.move(1)
	case " ":
		m.profile.toggleCurrent()
	case "d":
		m.profile.setDefaultCurrent()
	case "s":
		if m.profile.Mode == profileFlowOnboarding {
			if m.profile.skipCurrent() {
				if profileHasValues(m.profile.Draft) {
					m.profile.Draft.Onboarding.Completed = true
					m.profile.Draft.Onboarding.Skipped = false
				} else {
					m.profile.Draft.Onboarding.Completed = false
					m.profile.Draft.Onboarding.Skipped = true
				}
				if err := updateConfigProfile(m.profile.Draft); err != nil {
					return m.showError(err)
				}
				m.profile = nil
				m.status = m.text.Ready
				m.mode = m.text.Auto
				m.view.SetContent(welcomeContent(m.version, m.text, m.view.Width, m.view.Height))
				m.view.GotoTop()
				return m, nil
			}
			break
		}
		m.profile.skipCurrent()
	case "enter":
		done, _ := m.profile.saveCurrent(m.lang)
		if done {
			if profileHasValues(m.profile.Draft) {
				m.profile.Draft.Onboarding.Completed = true
				m.profile.Draft.Onboarding.Skipped = false
			}
			if err := updateConfigProfile(m.profile.Draft); err != nil {
				return m.showError(err)
			}
			if m.profile.Mode == profileFlowOnboarding {
				summary := profileSavedContent(m.profile.Draft, m.lang)
				m.profile = nil
				m.status = m.text.SettingsSaved
				m.mode = m.text.Auto
				m.view.SetContent(summary)
				m.view.GotoTop()
				return m, nil
			}
		}
	default:
		return m, nil
	}
	if m.profile != nil {
		m.view.SetContent(m.profile.render(m.lang, m.view.Width))
		m.view.GotoTop()
	}
	return m, nil
}

func profileSavedContent(profile gachaProfile, lang language) string {
	if lang == languageKorean {
		return strings.Join([]string{
			titleStyle.Render("프로필을 저장했습니다."),
			"",
			profileDetails(profile, lang),
			"",
			"투자 질문을 입력하세요.",
		}, "\n")
	}
	return strings.Join([]string{
		titleStyle.Render("Profile saved."),
		"",
		profileDetails(profile, lang),
		"",
		"Ask an investment question.",
	}, "\n")
}

func (f *profileFlow) renderMenu(lang language, width int) string {
	lines := []string{profileDetails(f.Draft, lang), ""}
	items := profileMenuLabels(lang)
	for i, item := range items {
		marker := " "
		if i == f.Cursor {
			marker = "›"
		}
		lines = append(lines, bulletStyle.Render(marker)+" "+item)
	}
	if strings.TrimSpace(f.Message) != "" {
		lines = append(lines, "", warningStyle.Render(wrapLine(f.Message, max(24, width-4))))
	}
	lines = append(lines, "", mutedStyle.Render(profileMenuHint(lang)))
	return strings.Join(lines, "\n")
}

func (f *profileFlow) renderEdit(lang language, width int) string {
	lines := []string{}
	if f.Mode == profileFlowOnboarding {
		lines = append(lines, brandLine(lang), "")
		lines = append(lines, wrapParagraphs(profileOnboardingIntro(lang), max(24, width-4)), "")
	}
	lines = append(lines, titleStyle.Render(profileCategoryTitle(f.Category, lang)))
	if profileCategoryIsMulti(f.Category) {
		lines = append(lines, mutedStyle.Render(profileMultiHint(lang)))
	} else {
		lines = append(lines, mutedStyle.Render(profileSingleHint(lang)))
	}
	lines = append(lines, "")
	options := profileOptionsFor(f.Category)
	for i, option := range options {
		marker := " "
		if i == f.Cursor {
			marker = "›"
		}
		lines = append(lines, bulletStyle.Render(marker)+" "+f.renderOption(option, lang))
	}
	if strings.TrimSpace(f.Message) != "" {
		lines = append(lines, "", warningStyle.Render(wrapLine(f.Message, max(24, width-4))))
	}
	return strings.Join(lines, "\n")
}

func (f *profileFlow) renderOption(option profileOption, lang language) string {
	label := profileLabel(option.Value, profileOptionsFor(f.Category), lang)
	if profileCategoryIsMulti(f.Category) {
		check := "[ ]"
		if profileMultiContains(f.multi(), option.Value) {
			check = "[x]"
		}
		suffix := ""
		if f.multi().Default == option.Value {
			suffix = mutedStyle.Render(" default")
		}
		return fmt.Sprintf("%s %s%s", check, label, suffix)
	}
	radio := "( )"
	if f.single() == option.Value {
		radio = "(*)"
	}
	return fmt.Sprintf("%s %s", radio, label)
}

func profileMenuLabels(lang language) []string {
	if lang == languageKorean {
		return []string{
			"관심 시장 수정",
			"투자 기간 수정",
			"리스크 성향 수정",
			"리포트 스타일 수정",
			"자주 하는 판단 수정",
			"프로필 초기화",
		}
	}
	return []string{
		"Edit markets",
		"Edit horizons",
		"Edit risk preference",
		"Edit report style",
		"Edit common goals",
		"Reset profile",
	}
}

func profileMenuHint(lang language) string {
	if lang == languageKorean {
		return "↑/↓ 이동   enter 수정   esc 돌아가기"
	}
	return "↑/↓ choose   enter edit   esc back"
}

func profileMultiHint(lang language) string {
	if lang == languageKorean {
		return "space 선택   d 기본값   enter 계속   s 건너뛰기"
	}
	return "space toggle   d default   enter continue   s skip"
}

func profileSingleHint(lang language) string {
	if lang == languageKorean {
		return "↑/↓ 이동   enter 계속   s 건너뛰기"
	}
	return "↑/↓ choose   enter continue   s skip"
}

func brandLine(lang language) string {
	if lang == languageKorean {
		return titleStyle.Render("GACHA") + "\n" + mutedStyle.Render("리서치로 확률을 높입니다.")
	}
	return titleStyle.Render("GACHA") + "\n" + mutedStyle.Render("Better odds through research.")
}

func profileOnboardingIntro(lang language) string {
	if lang == languageKorean {
		return "투자 결과를 완벽하게 예측할 수는 없습니다. Gacha는 꼼꼼한 리서치로 더 나은 판단 기준을 만들도록 돕습니다.\n\n투자 프로필을 설정하겠습니다. 언제든 /profile에서 바꿀 수 있습니다."
	}
	return "Investing is uncertain. Gacha helps improve the odds with disciplined research.\n\nLet's set your research profile. You can change this anytime with /profile."
}

func profileCategoryIsMulti(category profileCategory) bool {
	return category == profileCategoryMarkets || category == profileCategoryHorizons || category == profileCategoryGoals
}

func profileMultiContains(value profileMulti, item string) bool {
	for _, selected := range value.Selected {
		if selected == item {
			return true
		}
	}
	return false
}

func (f *profileFlow) multi() profileMulti {
	switch f.Category {
	case profileCategoryMarkets:
		return f.Draft.Markets
	case profileCategoryHorizons:
		return f.Draft.Horizons
	case profileCategoryGoals:
		return f.Draft.Goals
	default:
		return profileMulti{}
	}
}

func (f *profileFlow) setMulti(value profileMulti) {
	switch f.Category {
	case profileCategoryMarkets:
		f.Draft.Markets = value
	case profileCategoryHorizons:
		f.Draft.Horizons = value
	case profileCategoryGoals:
		f.Draft.Goals = value
	}
	f.Draft = normalizeProfile(f.Draft)
}

func (f *profileFlow) single() string {
	switch f.Category {
	case profileCategoryRisk:
		return f.Draft.Risk
	case profileCategoryReportStyle:
		return f.Draft.ReportStyle
	default:
		return ""
	}
}

func (f *profileFlow) setSingle(value string) {
	switch f.Category {
	case profileCategoryRisk:
		f.Draft.Risk = value
	case profileCategoryReportStyle:
		f.Draft.ReportStyle = value
	}
	f.Draft = normalizeProfile(f.Draft)
}

func (f *profileFlow) move(delta int) {
	maxItems := len(profileOptionsFor(f.Category))
	if f.Mode == profileFlowMenu {
		maxItems = len(profileMenuLabels(languageEnglish))
	}
	if maxItems == 0 {
		return
	}
	next := (f.Cursor + delta) % maxItems
	if next < 0 {
		next += maxItems
	}
	f.Cursor = next
}

func (f *profileFlow) toggleCurrent() {
	if !profileCategoryIsMulti(f.Category) {
		return
	}
	options := profileOptionsFor(f.Category)
	if f.Cursor < 0 || f.Cursor >= len(options) {
		return
	}
	option := options[f.Cursor]
	value := f.multi()
	if option.Value == profileValueNotSure {
		if profileMultiContains(value, profileValueNotSure) {
			value = profileMulti{}
		} else {
			value = profileMulti{Selected: []string{profileValueNotSure}}
		}
		f.setMulti(value)
		return
	}
	selected := make([]string, 0, len(value.Selected)+1)
	removed := false
	for _, item := range value.Selected {
		if item == profileValueNotSure {
			continue
		}
		if item == option.Value {
			removed = true
			continue
		}
		selected = append(selected, item)
	}
	if !removed {
		selected = append(selected, option.Value)
	}
	value.Selected = selected
	if value.Default == option.Value && removed {
		value.Default = ""
	}
	f.setMulti(value)
}

func (f *profileFlow) setDefaultCurrent() {
	if !profileCategoryIsMulti(f.Category) {
		return
	}
	options := profileOptionsFor(f.Category)
	if f.Cursor < 0 || f.Cursor >= len(options) {
		return
	}
	option := options[f.Cursor]
	if option.Value == profileValueNotSure {
		return
	}
	value := f.multi()
	if !profileMultiContains(value, option.Value) {
		value.Selected = append(value.Selected, option.Value)
	}
	value.Default = option.Value
	f.setMulti(value)
}

func (f *profileFlow) skipCurrent() bool {
	if f.Mode == profileFlowOnboarding {
		return f.nextStep()
	}
	f.Mode = profileFlowMenu
	f.Cursor = 0
	f.Message = ""
	return false
}

func (f *profileFlow) saveCurrent(lang language) (bool, string) {
	f.Message = ""
	if f.Mode == profileFlowMenu {
		if f.Cursor == len(profileMenuLabels(lang))-1 {
			f.Draft = gachaProfile{}
			return true, ""
		}
		f.Category = profileMenuCategories[f.Cursor]
		f.Mode = profileFlowEdit
		f.Cursor = profileInitialCursor(f.Category, f.Draft)
		return false, ""
	}
	if profileCategoryIsMulti(f.Category) {
		value := f.multi()
		if len(value.Selected) == 0 {
			if lang == languageKorean {
				f.Message = "하나 이상 선택하거나 s로 건너뛰세요."
			} else {
				f.Message = "Select at least one option, or press s to skip."
			}
			return false, ""
		}
		if value.Default == "" && !(len(value.Selected) == 1 && value.Selected[0] == profileValueNotSure) {
			value.Default = value.Selected[0]
			f.setMulti(value)
		}
	} else {
		options := profileOptionsFor(f.Category)
		if f.Cursor >= 0 && f.Cursor < len(options) {
			f.setSingle(options[f.Cursor].Value)
		}
	}
	if f.Mode == profileFlowOnboarding {
		if f.nextStep() {
			f.Draft.Onboarding.Completed = true
			f.Draft.Onboarding.Skipped = false
			return true, ""
		}
		return false, ""
	}
	f.Mode = profileFlowMenu
	f.Cursor = 0
	return true, ""
}

func (f *profileFlow) nextStep() bool {
	f.Step++
	if f.Step >= len(profileOnboardingSteps) {
		return true
	}
	f.Category = profileOnboardingSteps[f.Step]
	f.Cursor = profileInitialCursor(f.Category, f.Draft)
	f.Message = ""
	return false
}

func profileInitialCursor(category profileCategory, profile gachaProfile) int {
	options := profileOptionsFor(category)
	value := ""
	switch category {
	case profileCategoryMarkets:
		value = profile.Markets.Default
		if value == "" && len(profile.Markets.Selected) > 0 {
			value = profile.Markets.Selected[0]
		}
	case profileCategoryHorizons:
		value = profile.Horizons.Default
		if value == "" && len(profile.Horizons.Selected) > 0 {
			value = profile.Horizons.Selected[0]
		}
	case profileCategoryGoals:
		value = profile.Goals.Default
		if value == "" && len(profile.Goals.Selected) > 0 {
			value = profile.Goals.Selected[0]
		}
	case profileCategoryRisk:
		value = profile.Risk
		if value == "" {
			value = profileRiskBalanced
		}
	case profileCategoryReportStyle:
		value = profile.ReportStyle
		if value == "" {
			value = profileReportBasicFirst
		}
	}
	for i, option := range options {
		if option.Value == value {
			return i
		}
	}
	return 0
}
