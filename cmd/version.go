package cmd

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of infra-doctor",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(utils.VersionString())
		},
	}
}


func init() {
	rootCmd.AddCommand(versionCommand())
}