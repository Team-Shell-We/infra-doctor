package cmd

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "check the current version and the latest version",
	Args:  cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		updateInfo, err := utils.GetUpdateInfo()
		if err != nil {
			return fmt.Errorf("failed to check update: %w", err)
		}

		fmt.Fprintf(
			cmd.OutOrStdout(),
			"Current\n\n%s\n\nLatest\n\n%s\n",
			updateInfo.CurrentVersion,
			updateInfo.LatestVersion,
		)

		if updateInfo.CurrentVersion == "dev" {
			fmt.Fprintln(
				cmd.OutOrStdout(),
				"\nUpdate status cannot be checked in development mode.",
			)
			return nil
		}

		if updateInfo.UpdateAvailable {
			fmt.Fprintln(
				cmd.OutOrStdout(),
				"\nRun\n\nbrew upgrade infra-doctor",
			)
		} else {
			fmt.Fprintln(
				cmd.OutOrStdout(),
				"\nInfra Doctor is already up to date.",
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}