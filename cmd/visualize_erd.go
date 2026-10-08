package cmd

import (
	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/visualize/erd"
	"github.com/spf13/cobra"
)

var erdFormat string
var erdOutput string

var visualizeERDCmd = &cobra.Command{
	Use:   "erd [path]",
	Short: "Visualize entity relationships (ERD)",
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

		diagram := erd.Build(*info)

		content, err := erd.Render(diagram, erd.Format(erdFormat))
		if err != nil {
			return err
		}

		return writeVisualization(content, erdOutput, cmd)
	},
}

func init() {
	erdFlags := visualizeERDCmd.Flags()
	erdFlags.StringVar(
		&erdFormat,
		"format",
		"ascii",
		"ascii, mermaid, or markdown",
	)
	erdFlags.StringVar(
		&erdOutput,
		"output",
		"",
		"write output to a file",
	)

	visualizeCmd.AddCommand(visualizeERDCmd)
}
