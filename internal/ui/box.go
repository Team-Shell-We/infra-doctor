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

func center(text string, width int) string {
	padding := width - len(text)
	if padding <= 0 {
		return text
	}

	left := padding / 2
	right := padding - left

	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}
