package utils

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

func PrintDonateInfo(lang string) {
	fmt.Printf(`
        ❤️ %s

		GitHub Sponsors

		https://github.com/Team-Shell-We/infra-doctor
`, i18n.Get(lang, "donate.thanks"))
}
