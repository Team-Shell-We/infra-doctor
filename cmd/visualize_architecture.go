package cmd

import (
	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/visualize"
	"github.com/spf13/cobra"
)

var architectureFormat string
var architectureOutput string

var visualizeArchitectureCmd = &cobra.Command{
	Use:   "architecture [path]",
	Short: "Visualize the runtime architecture",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) == 1 {
			root = args[0]
		}

		info, err := analyzer.AnalyzeProject(root)
		if err != nil {
			return err
		}

		diagram := visualize.Build(*info)

		return writeVisualization(
			diagram,
			architectureFormat,
			architectureOutput,
			cmd,
		)
	},
}

func init() {
	architectureFlags := visualizeArchitectureCmd.Flags()
	architectureFlags.StringVar(
		&architectureFormat,
		"format",
		"ascii",
		"ascii, mermaid, or markdown",
	)
	architectureFlags.StringVar(
		&architectureOutput,
		"output",
		"",
		"write output to a file",
	)

	visualizeCmd.AddCommand(visualizeArchitectureCmd)
}
