package utils

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

func PrintAbout(lang string) {
	fmt.Printf(`
        Infra Doctor

%s

%s

%s

github.com/Team-Shell-We/infra-doctor

%s

Apache 2.0
`,
		i18n.Get(lang, "about.tagline"),
		i18n.Get(lang, "about.builtFor"),
		i18n.Get(lang, "about.github"),
		i18n.Get(lang, "about.license"),
	)
}
