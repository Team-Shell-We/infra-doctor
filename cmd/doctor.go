package cmd

/*
doctor 명령을 정의
사용자로부터 경로인자를 받고, analyzer를 호출해 프로젝트 분석 결과를 받아 터미널에 요약을 출력
출력내용: 프레임워크, 데이터베이스, docker, docker compose 파일, ci/cd, profiles 등
*/
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

		info, err := analyzer.AnalyzeProject(root)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("╔════════════════════════════════════════════════════════════╗")
		fmt.Println("║                   🔍 Project Analysis                     ║")
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
		// Docker
		// ---------------------------------------------------------------------

		fmt.Println()
		fmt.Println(" Docker")

		for _, docker := range info.Infrastructure.Docker.Dockerfiles {
			fmt.Printf("   ✓ %s\n", docker.File)
		}

		for _, compose := range info.Infrastructure.Docker.Compose {
			fmt.Printf("   ✓ %s\n", compose.File)
		}

		// ---------------------------------------------------------------------
		// GitHub
		// ---------------------------------------------------------------------

		fmt.Println()
		fmt.Println(" CI/CD")

		for _, workflow := range info.Github.Workflows {
			fmt.Printf("   ✓ %s\n", workflow.File)
		}

		// ---------------------------------------------------------------------
		// Profiles
		// ---------------------------------------------------------------------

		fmt.Println()
		fmt.Println(" Profiles")

		for _, profile := range info.Profiles {
			fmt.Printf("   ✓ %-8s %s\n", profile.Name, profile.File)
		}

		fmt.Println()
		fmt.Println("╚════════════════════════════════════════════════════════════╝")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}