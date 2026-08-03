package cmd

import (
 	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func aboutCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "about",
		Short: "Show about CLI Information",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			
			utils.PrintAbout()
		},
	}
}

func init() {
	rootCmd.AddCommand(aboutCommand())
}