package cmd

import (
	"fmt"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan project structure",
	Long:  "Analyze the project structure and detect technologies in use.",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		root := "."

		if len(args) == 1 {
			root = args[0]
		}

		info, err := analyzer.AnalyzeProject(root)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("╔════════════════════════════════════════════════════════════╗")
		fmt.Println("║                     🔍 Project Scan                        ║")
		fmt.Println("╠════════════════════════════════════════════════════════════╣")

		// ---------------------------------------------------------------------
		// Framework
		// ---------------------------------------------------------------------

		fmt.Println()
		fmt.Println(" Framework")

		if info.Framework.SpringBoot.Version != "" {
			fmt.Printf("   ✓ Spring Boot %s\n", info.Framework.SpringBoot.Version)
		}

		fmt.Printf("   ✓ %s\n", info.Framework.BuildTool.Type)

		if info.Framework.Java.Version != "" {
			fmt.Printf("   ✓ Java %s\n", info.Framework.Java.Version)
		}

		// ---------------------------------------------------------------------
		// Dependencies
		// ---------------------------------------------------------------------

		fmt.Println()
		fmt.Println(" Dependencies")

		if info.Dependencies.Security.Enabled {
			fmt.Println("   ✓ Spring Security")
		}

		if info.Dependencies.JPA.Enabled {
			fmt.Println("   ✓ Spring Data JPA")
		}

		if info.Dependencies.Kafka.Enabled {
			fmt.Println("   ✓ Kafka")
		}

		if info.Dependencies.AWS.Enabled {
			fmt.Println("   ✓ AWS SDK")
		}

		if info.Dependencies.Lombok.Enabled {
			fmt.Println("   ✓ Lombok")
		}

		if info.Dependencies.Actuator.Enabled {
			fmt.Println("   ✓ Spring Boot Actuator")
		}

		if info.Dependencies.OpenAPI.Enabled {
			fmt.Println("   ✓ SpringDoc OpenAPI")
		}

		// ---------------------------------------------------------------------
		// Database
		// ---------------------------------------------------------------------

		fmt.Println()
		fmt.Println(" Database")

		if info.Database.Primary.Type != "" &&
			info.Database.Primary.Type != "Unknown" {
			fmt.Printf("   ✓ %s\n", info.Database.Primary.Type)
		}

		if info.Database.Redis != nil {
			fmt.Println("   ✓ Redis")
		}

		// ---------------------------------------------------------------------
		// Infrastructure
		// ---------------------------------------------------------------------

		fmt.Println()
		fmt.Println(" Infrastructure")

		// Docker
		if info.Infrastructure.Docker.Enabled {

			fmt.Println("   ✓ Docker")

			for _, docker := range info.Infrastructure.Docker.Dockerfiles {
				fmt.Printf("      └─ %s\n", docker.File)
			}

		} else {
			fmt.Println("   ✗ Docker")
		}

		// Docker Compose
		if len(info.Infrastructure.Docker.Compose) > 0 {

			fmt.Println("   ✓ Docker Compose")

			for _, compose := range info.Infrastructure.Docker.Compose {
				fmt.Printf("      └─ %s\n", compose.File)
			}

		} else {
			fmt.Println("   ✗ Docker Compose")
		}

		// Kubernetes
		if info.Infrastructure.Kubernetes.Enabled {

			fmt.Println("   ✓ Kubernetes")

			for _, file := range info.Infrastructure.Kubernetes.Files {
				fmt.Printf("      └─ %s\n", file.File)
			}

		} else {

			fmt.Println("   ✗ Kubernetes")
		}

		// ---------------------------------------------------------------------
		// CI/CD
		// ---------------------------------------------------------------------

		fmt.Println()
		fmt.Println(" CI/CD")

		if len(info.Github.Workflows) > 0 {

			fmt.Println("   ✓ GitHub Actions")

			for _, workflow := range info.Github.Workflows {

				fmt.Printf("      └─ %s\n", workflow.Name)

				for _, trigger := range workflow.Triggers {

					fmt.Printf("         Trigger : %s\n", trigger.Event)

					if len(trigger.Branches) > 0 {
						fmt.Printf("         Branch  : %s\n",
							strings.Join(trigger.Branches, ", "))
					}
				}

				if len(workflow.Jobs) > 0 {

					var jobs []string

					for _, job := range workflow.Jobs {
						jobs = append(jobs, job.Name)
					}

					fmt.Printf("         Jobs    : %s\n",
						strings.Join(jobs, ", "))
				}
			}

		} else {

			fmt.Println("   ✗ GitHub Actions")
		}

		// ---------------------------------------------------------------------
		// Profiles
		// ---------------------------------------------------------------------

		// ---------------------------------------------------------------------
		// Profiles
		// ---------------------------------------------------------------------

		fmt.Println()
		fmt.Println(" Profiles")

		if len(info.Profiles) > 0 {

			for _, profile := range info.Profiles {

				fmt.Printf("   ✓ %s\n", profile.Name)
				fmt.Printf("      └─ %s\n", profile.File)
			}

		} else {

			fmt.Println("   ✗ Profiles")
		}

		fmt.Println()
		fmt.Println("╚════════════════════════════════════════════════════════════╝")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
