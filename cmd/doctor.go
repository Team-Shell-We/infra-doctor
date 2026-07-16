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

		info, err := analyzer.AnalyzeProject(root)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("Framework")
		fmt.Printf("  Build Tool     : %s\n", info.Framework.BuildTool)
		fmt.Printf("  Spring Boot    : %t\n", info.Framework.SpringBoot)
		fmt.Printf("  Spring Version : %s\n", info.Framework.SpringVersion)
		fmt.Printf("  Java Version   : %s\n", info.Framework.JavaVersion)

		fmt.Println()

		fmt.Println("Database")
		fmt.Printf("  Primary DB   : %s\n", info.Database.Type)

		if info.Database.Redis {
			fmt.Println("  Redis        : Enabled")
		} else {
			fmt.Println("  Redis        : Disabled")
		}
		
		fmt.Println()


		fmt.Println("Profiles")

		if len(info.Profiles) == 0 {
			fmt.Println("  None")
		} else {
			for _, profile := range info.Profiles {
				fmt.Printf("  ✓ %-8s %s\n", profile.Name, profile.File)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}	