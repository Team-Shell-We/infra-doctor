package cmd

import (
	"github.com/Team-Shell-We/infra-doctor/internal/ai"
	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

// currentLang : 저장된 언어 설정을 읽는다(로그인 여부와 무관, 설정 안
// 돼있으면 영어). 모든 cmd/*.go가 출력 시작 전에 이걸로 언어를 결정한다.
func currentLang() string {

	lang := ai.LoadOrDefault().Language
	if lang == "" {
		return i18n.English
	}

	return lang
}
