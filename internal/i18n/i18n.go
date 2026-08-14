// Package i18n : CLI 출력 문구 다국어 지원. 고정 문구(헤더/라벨/에러
// 메시지)는 이 패키지의 key-lookup으로 처리한다. explain/recommend의
// AI 생성 문장은 정적 번역 대상이 아니라 프롬프트에 언어 지시를 넣는
// 별도 메커니즘을 쓴다(internal/ai/explain, internal/ai/recommend의
// prompt.go 참고). 실제 문구 데이터는 messages.go에 있다.
package i18n

const (
	Korean  = "ko"
	English = "en"
)

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
