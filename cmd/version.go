package cmd

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of infra-doctor",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			lang := currentLang()

			fmt.Printf(
				"Infra Doctor\n\n%s %s\n\n%s %s\n",
				i18n.Get(lang, "version.label"), utils.LocalizeVersion(lang, utils.Version()),
				i18n.Get(lang, "version.goLabel"), utils.LocalizeVersion(lang, utils.GoVersion()),
			)
		},
	}
}

func init() {
	rootCmd.AddCommand(versionCommand())
}
