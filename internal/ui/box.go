package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const boxWidth = 60

func Header(title string) {
	fmt.Println("╔" + strings.Repeat("═", boxWidth) + "╗")
	fmt.Printf("║ %-57s ║\n", center(title, boxWidth))
	fmt.Println("╚" + strings.Repeat("═", boxWidth) + "╝")
}

func Footer() {
	fmt.Println("╚" + strings.Repeat("═", boxWidth) + "╝")
}

func Blank() {
	fmt.Printf("║ %-60s ║\n", "")
}

func Line(text string) {
	fmt.Printf("║ %-60s ║\n", text)
}

func ProgressBar(percent, width int) string {

	if percent < 0 {
		percent = 0
	}

	if percent > 100 {
		percent = 100
	}

	filled := (percent * width) / 100

	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// Wrap : 텍스트를 width보다 길지 않은 줄들로 나눔
// 단어 단위로 끊어서 고정폭 박스 안에 긴 텍스트도 출력할 수 있도록
func Wrap(text string, width int) []string {

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	current := words[0]

	for _, word := range words[1:] {

		if len(current)+1+len(word) > width {
			lines = append(lines, current)
			current = word
			continue
		}

		current += " " + word
	}

	return append(lines, current)
}

// center : text를 width 안에서 가운데 정렬. rune 개수 기준으로 계산 —
// len()(바이트 길이)을 쓰면 한글 등 멀티바이트 문자가 섞인 제목에서
// 패딩이 틀어진다. (완전한 터미널 컬럼 폭 정렬은 한글이 2칸을 차지하는
// East Asian Width까지 고려해야 해서 별도이며, 여기선 rune 개수까지만 맞춘다.)
func center(text string, width int) string {
	padding := width - utf8.RuneCountInString(text)
	if padding <= 0 {
		return text
	}

	left := padding / 2
	right := padding - left

	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}
