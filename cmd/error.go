package cmd

import (
 	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func errorCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "help",
		Short: "if error cammand, recommend infra-doctor help",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			
			utils.PrintError()
		},
	}
}

func init() {
	rootCmd.AddCommand(errorCommand())
}