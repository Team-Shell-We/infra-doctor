package cmd

import (
	"fmt"
	"os"

	"github.com/Team-Shell-We/infra-doctor/internal/visualize"
	"github.com/spf13/cobra"
)

var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize the analyzed project infrastructure",
}

func writeVisualization(
	diagram visualize.Diagram,
	format string,
	output string,
	cmd *cobra.Command,
) error {
	content, err := visualize.Render(
		diagram,
		visualize.Format(format),
	)
	if err != nil {
		return err
	}

	if output != "" {
		return os.WriteFile(output, []byte(content), 0o644)
	}

	_, err = fmt.Fprint(cmd.OutOrStdout(), content)
	return err
}

func init() {
	rootCmd.AddCommand(visualizeCmd)
}
