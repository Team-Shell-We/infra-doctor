package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize the analyzed project infrastructure",
}

// writeVisualization : 이미 렌더링된 다이어그램 문자열을 파일(output이 있으면) 또는
// stdout에 쓴다. architecture/flow/erd 각각 자신의 Diagram 타입으로 먼저 Render를
// 호출한 뒤 결과 문자열만 이 함수에 넘긴다 — 그래야 서로 다른 visualize.Diagram/
// erd.Diagram 타입에 종속되지 않는다.
func writeVisualization(
	content string,
	output string,
	cmd *cobra.Command,
) error {
	if output != "" {
		return os.WriteFile(output, []byte(content), 0o644)
	}

	_, err := fmt.Fprint(cmd.OutOrStdout(), content)
	return err
}

func init() {
	rootCmd.AddCommand(visualizeCmd)
}
