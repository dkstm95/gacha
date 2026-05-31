package app

import (
	"testing"
)

func TestDetectLanguageFromLocale(t *testing.T) {
	lang := detectLanguageFromEnv(func(key string) string {
		if key == "LANG" {
			return "ko_KR.UTF-8"
		}
		return ""
	})
	if lang != languageKorean {
		t.Fatalf("unexpected language: %s", lang)
	}
}

func TestResponseLanguageFallsBackToDetectedLocale(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("LANG", "en_US.UTF-8")
	if got := responseLanguage("Should I buy NVDA now?"); got != languageEnglish {
		t.Fatalf("unexpected language: %s", got)
	}
	if got := responseLanguage("NVDA 지금 사도 될까?"); got != languageKorean {
		t.Fatalf("unexpected language: %s", got)
	}
}

func TestDetectLanguageUsesConfigOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("GACHA_LANG", "")
	t.Setenv("LANG", "en_US.UTF-8")
	if err := saveGachaConfig(gachaConfig{Language: languageSettingKorean}); err != nil {
		t.Fatal(err)
	}
	if got := detectLanguage(); got != languageKorean {
		t.Fatalf("unexpected language: %s", got)
	}
}

func TestGachaLangEnvOverridesConfigLanguage(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("GACHA_LANG", "en")
	if err := saveGachaConfig(gachaConfig{Language: languageSettingKorean}); err != nil {
		t.Fatal(err)
	}
	if got := detectLanguage(); got != languageEnglish {
		t.Fatalf("unexpected language: %s", got)
	}
}
