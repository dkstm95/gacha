package app

import (
	"fmt"
	"sort"
	"strings"
)

func doctor() error {
	lang := detectLanguage()
	if lang == languageKorean {
		return doctorKorean()
	}
	return doctorEnglish()
}

func doctorEnglish() error {
	runtimeStatus := runtimeDoctorStatus()
	fmt.Println("Gacha Doctor")
	fmt.Println()
	fmt.Printf("Overall:      %s\n", overallDoctorStatus(runtimeStatus))
	fmt.Printf("AI runtime:   %s\n", runtimeStatus)
	fmt.Printf("Provider:     %s\n", providerDoctorStatus())
	resolution := resolveOpenCodeModel()
	fmt.Printf("Model:        %s\n", modelDoctorSummary(resolution))
	fmt.Printf("Reports:      %s\n", reportsDir())
	fmt.Println()
	fmt.Println("Details")
	fmt.Printf("Command:      %s\n", openCodeCommand)
	if resolved, ok := resolveCommand(openCodeCommand); ok && resolved != openCodeCommand {
		fmt.Printf("Resolved:     %s\n", resolved)
	}
	fmt.Printf("Auth file:    %s\n", openCodeAuthPath())
	if source := resolution.Source; source != "" {
		fmt.Printf("Model source: %s\n", source)
	}
	if runtimeStatus != "ready" {
		fmt.Println()
		fmt.Println("Next step:    run `gch setup`")
	}
	fmt.Println()
	fmt.Println("Fallback:     ready; Gacha can print a pasteable prompt if the runtime is unavailable.")
	return nil
}

func doctorKorean() error {
	runtimeStatus := runtimeDoctorStatus()
	fmt.Println("Gacha 점검")
	fmt.Println()
	fmt.Printf("전체 상태:    %s\n", koreanOverallDoctorStatus(runtimeStatus))
	fmt.Printf("AI 런타임:    %s\n", koreanRuntimeDoctorStatus(runtimeStatus))
	fmt.Printf("Provider:     %s\n", providerDoctorStatus())
	resolution := resolveOpenCodeModel()
	fmt.Printf("모델:         %s\n", modelDoctorSummary(resolution))
	fmt.Printf("리포트 저장:  %s\n", reportsDir())
	fmt.Println()
	fmt.Println("세부 정보")
	fmt.Printf("명령어:       %s\n", openCodeCommand)
	if resolved, ok := resolveCommand(openCodeCommand); ok && resolved != openCodeCommand {
		fmt.Printf("위치:         %s\n", resolved)
	}
	fmt.Printf("인증 파일:    %s\n", openCodeAuthPath())
	if source := resolution.Source; source != "" {
		fmt.Printf("모델 출처:    %s\n", source)
	}
	if runtimeStatus != "ready" {
		fmt.Println()
		fmt.Println("다음 단계:    `gch setup` 실행")
	}
	fmt.Println()
	fmt.Println("대체 경로:    준비됨; 런타임을 사용할 수 없으면 붙여넣기용 프롬프트를 출력합니다.")
	return nil
}

func runtimeDoctorStatus() string {
	if !hasRunnableCommand(openCodeCommand) {
		return "missing"
	}
	if !hasOpenCodeAuth() {
		return "login required"
	}
	return "ready"
}

func overallDoctorStatus(runtimeStatus string) string {
	if runtimeStatus == "ready" {
		return "ready"
	}
	return "setup needed"
}

func koreanOverallDoctorStatus(runtimeStatus string) string {
	if runtimeStatus == "ready" {
		return "준비됨"
	}
	return "설정 필요"
}

func koreanRuntimeDoctorStatus(runtimeStatus string) string {
	switch runtimeStatus {
	case "ready":
		return "준비됨"
	case "login required":
		return "로그인 필요"
	default:
		return "없음"
	}
}

func providerDoctorStatus() string {
	providers, err := openCodeAuthProviders()
	if err != nil || len(providers) == 0 {
		return "not connected"
	}
	names := make([]string, 0, len(providers))
	for name, credential := range providers {
		value := name
		if credential.Type != "" {
			value += " " + credential.Type
		}
		names = append(names, value)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func modelDoctorSummary(resolution modelResolution) string {
	if strings.HasPrefix(resolution.Reason, "auto: could not read provider model list") {
		return "OpenCode default (auto model discovery unavailable)"
	}
	return modelDescription(resolution)
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
