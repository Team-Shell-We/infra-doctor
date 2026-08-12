// Package i18n : 출력 문구 다국어 지원의 최소 메커니즘. 지금은 config
// 명령어 자체의 문구만 여기 등록돼 있고, scan/login/doctor 등 나머지
// 명령어의 출력 문구(~120개)는 아직 옮기지 않음 — 별도 작업으로 남겨둠.
package i18n

const (
	Korean  = "ko"
	English = "en"
)

var messages = map[string]map[string]string{
	English: {
		"config.llm":          "LLM",
		"config.language":     "Language",
		"config.output":       "Output",
		"config.autoExport":   "Auto Export",
		"config.enabled":      "Enabled",
		"config.disabled":     "Disabled",
		"config.unconfigured": "Not configured",
	},
	Korean: {
		"config.llm":          "LLM",
		"config.language":     "언어",
		"config.output":       "출력 형식",
		"config.autoExport":   "자동 내보내기",
		"config.enabled":      "사용",
		"config.disabled":     "미사용",
		"config.unconfigured": "설정 안 됨",
	},
}

// Get : lang에 해당하는 key의 문구를 반환. lang이 지원 안 되거나 key가
// 없으면 영어로, 영어에도 없으면 key 자체를 그대로 반환.
func Get(lang, key string) string {

	if translated, ok := messages[lang]; ok {
		if text, ok := translated[key]; ok {
			return text
		}
	}

	if text, ok := messages[English][key]; ok {
		return text
	}

	return key
}

// Supported : 현재 지원하는 언어 코드 목록.
func Supported() []string {
	return []string{English, Korean}
}

// IsSupported : lang이 지원하는 언어 코드인지 확인.
func IsSupported(lang string) bool {

	for _, supported := range Supported() {
		if lang == supported {
			return true
		}
	}

	return false
}
