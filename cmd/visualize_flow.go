package cmd

import (
	"os"
	"path/filepath"

	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/visualize"
	"github.com/spf13/cobra"
)

var flowFormat string
var flowOutput string

var visualizeFlowCmd = &cobra.Command{
	Use:   "flow [path]",
	Short: "Visualize the build and deployment pipeline",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot := "."

		if len(args) == 1 {
			projectRoot = args[0]
		}

		absoluteProjectRoot, err := filepath.Abs(projectRoot)
		if err != nil {
			return err
		}

		// 애플리케이션은 build.gradle이 있는 gogildong에서 분석한다.
		info, err := analyzer.AnalyzeProject(absoluteProjectRoot)
		if err != nil {
			return err
		}

		// GitHub Actions는 현재 경로 또는 가까운 상위 경로에서 찾는다.
		workflowRoot, found := findWorkflowRoot(
			absoluteProjectRoot,
		)

		if found && workflowRoot != absoluteProjectRoot {
			github, err := analyzer.AnalyzeGitHub(workflowRoot)
			if err != nil {
				return err
			}

			info.Github = *github
		}

		diagram, err := visualize.BuildDeploymentFlow(
			absoluteProjectRoot,
			*info,
		)
		if err != nil {
			return err
		}

		return writeVisualization(
			diagram,
			flowFormat,
			flowOutput,
			cmd,
		)
	},
}

// findWorkflowRoot searches the current directory and its parents for
// .github/workflows.
func findWorkflowRoot(start string) (string, bool) {
	current := filepath.Clean(start)

	for {
		workflows := filepath.Join(
			current,
			".github",
			"workflows",
		)

		stat, err := os.Stat(workflows)
		if err == nil && stat.IsDir() {
			return current, true
		}

		parent := filepath.Dir(current)

		if parent == current {
			return "", false
		}

		current = parent
	}
}

func init() {
	flowFlags := visualizeFlowCmd.Flags()

	flowFlags.StringVar(
		&flowFormat,
		"format",
		"ascii",
		"ascii, mermaid, or markdown",
	)

	flowFlags.StringVar(
		&flowOutput,
		"output",
		"",
		"write output to a file",
	)

	visualizeCmd.AddCommand(visualizeFlowCmd)
}