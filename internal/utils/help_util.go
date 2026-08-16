package utils

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

func PrintHelp(lang string) {
	fmt.Printf(`
        Infra Doctor %s

        %s
        infra-doctor <command> [arguments]

        %s

        doctor
            %s

        visualize architecture
            %s

        visualize flow
            %s

        explain <topic>
            %s

        recommend
            %s

        generate <target>
            %s

        export
            %s

        %s

        infra-doctor help <command>

        %s
`,
		LocalizeVersion(lang, Version()),
		i18n.Get(lang, "help.usage"),
		i18n.Get(lang, "help.coreCommands"),
		i18n.Get(lang, "help.doctorDesc"),
		i18n.Get(lang, "help.visualizeArchDesc"),
		i18n.Get(lang, "help.visualizeFlowDesc"),
		i18n.Get(lang, "help.explainDesc"),
		i18n.Get(lang, "help.recommendDesc"),
		i18n.Get(lang, "help.generateDesc"),
		i18n.Get(lang, "help.exportDesc"),
		i18n.Get(lang, "help.run"),
		i18n.Get(lang, "help.forDetails"),
	)
}
