/*
Copyright © 2026 Team-Shell-We
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd : 인자 없이 실행됐을 때의 최상위 명령어.
var rootCmd = &cobra.Command{
	Use:   "infra-doctor",
	Short: "AI-powered infrastructure analysis CLI for Spring Boot projects",
	Long: `Infra Doctor analyzes a Spring Boot project's current infrastructure,
diagnoses deployment readiness, and explains and recommends deployment
strategies in the context of your actual project.`,
}

// Execute : 모든 서브커맨드를 rootCmd에 등록하고 실행. main.main()에서 한 번만 호출됨.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
