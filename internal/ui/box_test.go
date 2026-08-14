package ui

import (
	"testing"
	"unicode/utf8"
)

// 회귀 테스트: len()(바이트 길이) 대신 rune 개수로 패딩을 계산해야
// 한글 등 멀티바이트 문자가 섞인 제목도 올바르게 정렬된다.
func TestCenterMultibyteText(t *testing.T) {

	got := center("한글", 10)

	if utf8.RuneCountInString(got) != 10 {
		t.Errorf("center(\"한글\", 10) rune count = %d, want 10 (got %q)", utf8.RuneCountInString(got), got)
	}
}

func TestCenterASCII(t *testing.T) {

	got := center("hi", 6)

	if utf8.RuneCountInString(got) != 6 {
		t.Errorf("center(\"hi\", 6) rune count = %d, want 6 (got %q)", utf8.RuneCountInString(got), got)
	}
}

func TestCenterTextLongerThanWidth(t *testing.T) {

	if got := center("very long title", 5); got != "very long title" {
		t.Errorf("center should return text unchanged when it's already >= width, got %q", got)
	}
}
