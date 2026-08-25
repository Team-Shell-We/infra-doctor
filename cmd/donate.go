package cmd

import (
	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func donateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "donate",
		Short: "Show donate information",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {

			utils.PrintDonateInfo(currentLang())
		},
	}
}

func init() {
	rootCmd.AddCommand(donateCommand())
}
