package cmd

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func helpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "help [command]",
		Short: "Show help information",
		Args:  cobra.ArbitraryArgs,

		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				utils.PrintHelp(currentLang())
				return
			}

			target, _, err := cmd.Root().Find(args)
			if target == nil || err != nil {
				fmt.Printf("Unknown help topic %#q\n", args)
				_ = cmd.Root().Usage()
				return
			}

			target.InitDefaultHelpFlag()
			_ = target.Help()
		},
	}
}

func init() {
	// AddCommand 대신 SetHelpCommand를 쓴다 — "help"는 Cobra 내장 명령 이름이라 AddCommand로 등록하면 기본 help와 겹쳐 --help에 두 번 나온다
	rootCmd.SetHelpCommand(helpCommand())
}
