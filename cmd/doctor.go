package cmd

import (
	"fmt"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/Team-Shell-We/infra-doctor/internal/ui"
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

		ui.Header("🩺 Infrastructure Doctor")
		ui.Blank()

		ui.Line("Deployment Readiness")
		ui.Blank()
		ui.Line(fmt.Sprintf(" %s %d%%", ui.ProgressBar(result.Score, 30), result.Score))
		ui.Blank()

		ui.Line("Infrastructure Check")
		ui.Blank()

		for _, check := range result.Checks {

			mark := "✗"
			if check.Passed {
				mark = "✓"
			}

			ui.Line(fmt.Sprintf(" %s %s", mark, check.Name))
		}

		ui.Blank()

		ui.Line("Recommendation")
		ui.Blank()

		if len(result.Diagnoses) == 0 {

			ui.Line(" ✓ No issues found.")

		} else {

			for _, d := range result.Diagnoses {

				wrapped := ui.Wrap(strings.TrimSpace(d.Fix), 56)

				for i, line := range wrapped {

					if i == 0 {
						ui.Line(fmt.Sprintf(" • %s", line))
					} else {
						ui.Line(fmt.Sprintf("   %s", line))
					}
				}
			}
		}

		ui.Footer()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
