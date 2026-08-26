package cmd

import (
	"github.com/Team-Shell-We/infra-doctor/internal/ai"
	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

// currentLang : 저장된 언어 설정을 반환한다(로그인 여부 무관, 미설정 시 영어) — 모든 cmd 커맨드가 출력 전 언어 결정에 쓴다
func currentLang() string {

	lang := ai.LoadOrDefault().Language
	if lang == "" {
		return i18n.English
	}

	return lang
}
