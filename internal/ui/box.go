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

// Wrap : 텍스트를 화면 폭 w보다 넓지 않은 줄들로 나눔
func Wrap(text string, w int) []string {

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	current := words[0]

	for _, word := range words[1:] {

		if DisplayWidth(current)+1+DisplayWidth(word) > w {
			lines = append(lines, current)
			current = word
			continue
		}

		current += " " + word
	}

	return append(lines, current)
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
