package i18n

import "testing"

func TestGetKnownKey(t *testing.T) {

	if got := Get(Korean, "config.language"); got != "언어" {
		t.Errorf("Get(ko, config.language) = %q, want 언어", got)
	}

	if got := Get(English, "config.language"); got != "Language" {
		t.Errorf("Get(en, config.language) = %q, want Language", got)
	}
}

func TestGetFallsBackToEnglish(t *testing.T) {

	if got := Get("fr", "config.language"); got != "Language" {
		t.Errorf("Get(fr, config.language) should fall back to English, got %q", got)
	}
}

func TestGetUnknownKeyReturnsKeyItself(t *testing.T) {

	if got := Get(Korean, "no.such.key"); got != "no.such.key" {
		t.Errorf("Get for an unknown key should return the key itself, got %q", got)
	}
}

func TestIsSupported(t *testing.T) {

	if !IsSupported("ko") || !IsSupported("en") {
		t.Error("ko and en should both be supported")
	}

	if IsSupported("fr") {
		t.Error("fr should not be supported")
	}
}

// 두 언어 맵에 정의된 key 집합이 정확히 같은지 확인 — 한쪽에만 있는 key가
// 생기면 Get이 조용히 영어로 새는데, 그걸 미리 잡아낸다.
func TestLanguagesHaveMatchingKeys(t *testing.T) {

	enKeys := keySet(messages[English])
	koKeys := keySet(messages[Korean])

	for key := range enKeys {
		if !koKeys[key] {
			t.Errorf("key %q exists in English but not Korean", key)
		}
	}

	for key := range koKeys {
		if !enKeys[key] {
			t.Errorf("key %q exists in Korean but not English", key)
		}
	}
}

func keySet(m map[string]string) map[string]bool {

	set := make(map[string]bool, len(m))

	for key := range m {
		set[key] = true
	}

	return set
}
