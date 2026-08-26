/*
Copyright © 2026 Team-Shell-We
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd : 모든 서브커맨드가 등록되는 최상위 명령
var rootCmd = &cobra.Command{
	Use:   "infra-doctor",
	Short: "AI-powered infrastructure analysis CLI for Spring Boot projects",
	Long: `Infra Doctor analyzes a Spring Boot project's current infrastructure,
diagnoses deployment readiness, and explains and recommends deployment
strategies in the context of your actual project.`,
}

// Execute : rootCmd를 실행한다(서브커맨드는 각 cmd 파일의 init()에서 이미 등록됨) — main.main()에서 한 번만 호출
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
