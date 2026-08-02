package utils

import (
    "fmt"

    "github.com/spf13/cobra"
)

func helpCommand(rootCmd *cobra.Command) *cobra.Command {
    return &cobra.Command{
        Use:   "help [command]",
        Short: "Show help information",
        Long:  "Show help information for commands",
        Args:  cobra.MaximumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            if len(args) == 0 {
                PrintMainHelp()
                return
            }

            target, _, err := rootCmd.Find(args)
            if err == nil && target != nil {
                _ = target.Help()
                return
            }

            PrintMainHelp()
        },
    }
}

func PrintMainHelp() {
    fmt.Println(`Infra Doctor v1.0.0

Usage
  infra-doctor <command> [arguments]

Core Commands

  doctor
      Analyze project readiness

  visualize architecture
      Show infrastructure architecture

  visualize flow
      Show deployment flow

  explain <topic>
      Explain infrastructure concepts

  recommend
      Recommend deployment strategy

  generate <target>
      Generate configuration files

  export
      Export project report

Run

  infra-doctor help <command>

for detailed information.`)
}
