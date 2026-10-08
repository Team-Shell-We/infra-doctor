package cmd

import (
	"context"
	"io"

	exportapp "github.com/Team-Shell-We/infra-doctor/internal/export"
	"github.com/spf13/cobra"
)

type ExportRunner interface {
	Run(context.Context, exportapp.Request, io.Reader, io.Writer) error
}

func exportCommand(runner ExportRunner) *cobra.Command {
	request := exportapp.Request{}
	command := &cobra.Command{
		Use: "export [path]", Short: "Export analysis and generated infrastructure files",
		Args: cobra.MaximumNArgs(1), SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) == 1 {
				request.Root = args[0]
			}
			request.Lang = currentLang()
			return runner.Run(command.Context(), request, command.InOrStdin(), command.OutOrStdout())
		},
	}
	command.Flags().BoolVarP(&request.Force, "force", "f", false, "overwrite existing export files")
	command.Flags().BoolVar(&request.DryRun, "dry-run", false, "show export files without writing")
	command.Flags().BoolVar(&request.Select, "select", false, "choose which export items to include interactively")
	return command
}

func init() {
	rootCmd.AddCommand(exportCommand(exportapp.NewApplication()))
}
