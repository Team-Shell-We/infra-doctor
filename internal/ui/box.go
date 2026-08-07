package ui

import (
	"fmt"
	"strings"
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

func center(text string, width int) string {
	padding := width - len(text)
	if padding <= 0 {
		return text
	}

	left := padding / 2
	right := padding - left

	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}
