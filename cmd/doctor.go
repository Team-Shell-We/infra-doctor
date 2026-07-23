package cmd

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor [path]",
	Short: "Analyze project and provide recommendations",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		root := "."

		if len(args) == 1 {
			root = args[0]
		}

		info, err := analyzer.AnalyzeProject(root)
		if err != nil {
			fmt.Println(err)
			return
		}

		result := doctor.Analyze(info)

		fmt.Println("Infra Doctor Report")
		fmt.Println("-------------------")

		if len(result.Diagnoses) == 0 {
			fmt.Println("✓ No issues found.")
			return
		}

		for _, d := range result.Diagnoses {
			fmt.Printf("[%s] %s\n", d.Level, d.Title)
			fmt.Printf("  %s\n", d.Message)
			fmt.Printf("  Recommendation: %s\n\n", d.Fix)
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
