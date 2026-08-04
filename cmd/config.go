package cmd

import (
 	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func configCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "config",
		Short: "한영 변환 가능하게 하는 유틸",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			
			utils.PrintConfig()
		},
	}
}

func init() {
	rootCmd.AddCommand(configCommand())
}