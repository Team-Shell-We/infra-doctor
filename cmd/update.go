package cmd

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "check the current version and the latest version",
	Args:  cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {

		lang := currentLang()

		updateInfo, err := utils.GetUpdateInfo()
		if err != nil {
			return fmt.Errorf("failed to check update: %w", err)
		}

		fmt.Fprintf(
			cmd.OutOrStdout(),
			"%s\n\n%s\n\n%s\n\n%s\n",
			i18n.Get(lang, "update.current"),
			updateInfo.CurrentVersion,
			i18n.Get(lang, "update.latest"),
			updateInfo.LatestVersion,
		)

		if updateInfo.CurrentVersion == "dev" {
			fmt.Fprintln(
				cmd.OutOrStdout(),
				"\n"+i18n.Get(lang, "update.devMode"),
			)
			return nil
		}

		if updateInfo.UpdateAvailable {
			fmt.Fprintln(
				cmd.OutOrStdout(),
				"\n"+i18n.Get(lang, "update.run")+"\n\nbrew upgrade infra-doctor",
			)
		} else {
			fmt.Fprintln(
				cmd.OutOrStdout(),
				"\n"+i18n.Get(lang, "update.upToDate"),
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
