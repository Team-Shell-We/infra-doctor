package ui

import (
	"fmt"
	"strings"

	"golang.org/x/text/width"
)

const boxWidth = 60

const borderWidth = boxWidth + 2

func Header(title string) {
	fmt.Println("╔" + strings.Repeat("═", borderWidth) + "╗")

	fmt.Println("║ " + Center(title, boxWidth) + " ║")

	fmt.Println("╚" + strings.Repeat("═", borderWidth) + "╝")
}

func Footer() {
	fmt.Println("╚" + strings.Repeat("═", borderWidth) + "╝")
}

func Blank() {

	fmt.Println("║ " + PadRight("", boxWidth) + " ║")
}

func Line(text string) {
	fmt.Println("║ " + PadRight(text, boxWidth) + " ║")

}

func ProgressBar(percent, w int) string {

	if percent < 0 {
		percent = 0
	}

	if percent > 100 {
		percent = 100
	}

	filled := (percent * w) / 100

	return strings.Repeat("█", filled) + strings.Repeat("░", w-filled)
}

// DisplayWidth : text가 터미널에서 실제로 차지하는 컬럼 수를 계산
func DisplayWidth(text string) int {

	w := 0

	for _, r := range text {
		switch width.LookupRune(r).Kind() {
		case width.EastAsianWide, width.EastAsianFullwidth:
			w += 2
		default:
			w += 1
		}
	}

	return w
}

// Wrap : 텍스트를 화면 폭 w보다 넓지 않은 줄들로 나눔. 단어 자체가 w보다
// 넓으면(URL 등 공백 없는 긴 토큰) splitByWidth로 강제로 끊어서
// 어떤 줄도 w를 넘지 않도록 함
func Wrap(text string, w int) []string {

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var current string

	for _, word := range words {

		for DisplayWidth(word) > w {
			var chunk string
			chunk, word = splitByWidth(word, w)

			if current != "" {
				lines = append(lines, current)
				current = ""
			}

			lines = append(lines, chunk)
		}

		if current == "" {
			current = word
			continue
		}

		if DisplayWidth(current)+1+DisplayWidth(word) > w {
			lines = append(lines, current)
			current = word
			continue
		}

		current += " " + word
	}

	return append(lines, current)
}

// splitByWidth : s 앞부분에서 DisplayWidth가 w를 넘지 않는 최대 조각을 잘라
// (head, 나머지)로 반환. 룬 단위로 계산해 와이드 문자(한글 등)가 중간에
// 안 잘리게 함. w가 0 이하면 자를 수 없으므로 s 전체를 head로 반환(무한 루프 방지)
func splitByWidth(s string, w int) (head, rest string) {

	if w <= 0 {
		return s, ""
	}

	runes := []rune(s)
	width := 0

	for i, r := range runes {

		rw := DisplayWidth(string(r))

		if width+rw > w {
			return string(runes[:i]), string(runes[i:])
		}

		width += rw
	}

	return s, ""
}

// PadRight : text 뒤에 공백을 채워 DisplayWidth 기준으로 w까지 맞춤
func PadRight(text string, w int) string {

	padding := w - DisplayWidth(text)
	if padding <= 0 {
		return text
	}
	return text + strings.Repeat(" ", padding)
}

// Center : text를 w 안에서 가운데 정렬. DisplayWidth 기준으로 계산
func Center(text string, w int) string {
	padding := w - DisplayWidth(text)
	if padding <= 0 {
		return text
	}

	left := padding / 2
	right := padding - left

	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}
