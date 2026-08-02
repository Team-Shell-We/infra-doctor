package cmd

import (
 	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func helpCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "help",
		Short: "Show help information",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			
			utils.PrintHelp()
		},
	}
}

func init() {
	rootCmd.AddCommand(helpCommand())
}