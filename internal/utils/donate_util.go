package utils

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

func PrintDonateInfo(lang string) {
	fmt.Printf(`
❤️ %s

Ko-fi
https://ko-fi.com/shellwe
`, i18n.Get(lang, "donate.thanks"))
}
