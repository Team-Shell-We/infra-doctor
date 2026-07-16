package cmd

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor [path]",
	Short: "Analyze project deployment readiness",
	Long:  "Analyze the current project and check deployment readiness.",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		root := "."

		if len(args) == 1 {
			root = args[0]
		}

		fmt.Printf("Analyzing project: %s\n\n", root)

		_, err := analyzer.AnalyzeProject(root)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("\nDone.")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}