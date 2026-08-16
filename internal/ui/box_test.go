package ui

import (
	"strings"
	"testing"
)

// 회귀 테스트: 한글은 터미널에서 2칸(East Asian Wide)을 차지하므로
// DisplayWidth 기준으로 패딩해야 실제 터미널 폭에 정확히 맞음
func TestCenterMultibyteText(t *testing.T) {

	got := Center("한글", 10)

	if DisplayWidth(got) != 10 {
		t.Errorf("center(\"한글\", 10) display width = %d, want 10 (got %q)", DisplayWidth(got), got)
	}
}

func TestCenterASCII(t *testing.T) {

	got := Center("hi", 6)

	if DisplayWidth(got) != 6 {
		t.Errorf("center(\"hi\", 6) display width = %d, want 6 (got %q)", DisplayWidth(got), got)
	}
}


func TestExportedLayoutHelpers(t *testing.T) {
	if got := DisplayWidth(Center("한글", 10)); got != 10 {
		t.Fatalf("Center display width = %d, want 10", got)
	}
	if got := DisplayWidth(PadRight("한글", 10)); got != 10 {
		t.Fatalf("PadRight display width = %d, want 10", got)
	}
}


func TestCenterTextLongerThanWidth(t *testing.T) {

	if got := Center("very long title", 5); got != "very long title" {
		t.Errorf("center should return text unchanged when it's already >= width, got %q", got)
	}
}

func TestDisplayWidthHangulCountsAsTwoColumns(t *testing.T) {

	if got := DisplayWidth("한"); got != 2 {
		t.Errorf("DisplayWidth(\"한\") = %d, want 2", got)
	}

	if got := DisplayWidth("한글"); got != 4 {
		t.Errorf("DisplayWidth(\"한글\") = %d, want 4", got)
	}
}

func TestDisplayWidthASCIICountsAsOneColumnEach(t *testing.T) {

	if got := DisplayWidth("hello"); got != 5 {
		t.Errorf("DisplayWidth(\"hello\") = %d, want 5", got)
	}
}

// 회귀 테스트
func TestWrapUsesDisplayWidthNotByteLength(t *testing.T) {

	lines := Wrap("한글 한글 한글", 9)

	for _, line := range lines {
		if DisplayWidth(line) > 9 {
			t.Errorf("Wrap line %q has display width %d, want <= 9", line, DisplayWidth(line))
		}
	}

	if len(lines) != 2 {
		t.Errorf("Wrap(\"한글 한글 한글\", 9) produced %d lines, want 2 (got %q)", len(lines), lines)
	}
}

// 회귀 테스트: Line/Blank/Header가 fmt의 %-Ns(rune 개수 기준) 대신
// DisplayWidth 기준으로 패딩해야, 한글이 섞인 행도 ASCII 행과 동일한 터미널 폭으로 정렬됨
func TestLineAndHeaderRowsMatchWidthForKoreanAndASCII(t *testing.T) {

	asciiRow := "║ " + PadRight("hello world", boxWidth) + " ║"
	koreanRow := "║ " + PadRight("한글도 잘 맞아야 함", boxWidth) + " ║"

	if DisplayWidth(asciiRow) != DisplayWidth(koreanRow) {
		t.Errorf("row widths differ: ascii=%d korean=%d", DisplayWidth(asciiRow), DisplayWidth(koreanRow))
	}

	asciiHeader := "║ " + Center("Infra Doctor", boxWidth) + " ║"
	koreanHeader := "║ " + Center("한글 제목 테스트", boxWidth) + " ║"

	if DisplayWidth(asciiHeader) != DisplayWidth(koreanHeader) {
		t.Errorf("header widths differ: ascii=%d korean=%d", DisplayWidth(asciiHeader), DisplayWidth(koreanHeader))
	}
}

// 회귀 테스트: 테두리 행("╔"+"═"*n+"╗")과 본문 행("║ "+content+" ║")의
// 총 폭이 정확히 같아야 오른쪽 테두리가 세로로 정렬됨
func TestBorderRowMatchesContentRowWidth(t *testing.T) {

	border := "╔" + strings.Repeat("═", borderWidth) + "╗"
	content := "║ " + Center("🔍 프로젝트 스캔", boxWidth) + " ║"

	if DisplayWidth(border) != DisplayWidth(content) {
		t.Errorf("border width = %d, content row width = %d, want equal", DisplayWidth(border), DisplayWidth(content))
	}
}
