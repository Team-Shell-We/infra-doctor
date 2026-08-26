package ai

// languageNames : i18n 언어 코드를 모델에게 지시할 사람이 읽는 이름으로
// 매핑. 고정 문구 번역(internal/i18n)과는 별개 메커니즘 — AI 응답은
// 정적 번역이 아니라 프롬프트 지시로 언어를 맞춘다.
var languageNames = map[string]string{
	"ko": "Korean",
	"en": "English",
}

// LanguageName : lang에 대응하는 사람이 읽는 언어 이름. 지원 안 하는
// 코드면 English로 폴백
func LanguageName(lang string) string {

	if name, ok := languageNames[lang]; ok {
		return name
	}

	return languageNames["en"]
}
