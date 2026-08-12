package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
	"github.com/Team-Shell-We/infra-doctor/internal/ai/explain"
	"github.com/Team-Shell-We/infra-doctor/internal/ai/openai"
	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/ui"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:       "explain <topic> [path]",
	Short:     "Explain an infrastructure concept in the context of your project",
	Args:      explainArgs,
	ValidArgs: explain.Topics,

	Run: func(cmd *cobra.Command, args []string) {

		topic := args[0]

		root := "."
		if len(args) == 2 {
			root = args[1]
		}

		creds, err := ai.Load()
		if err != nil {

			if errors.Is(err, ai.ErrNotLoggedIn) {
				fmt.Println("You're not logged in. Run 'infra-doctor login' to set up your OpenAI API Key first.")
				return
			}

			fmt.Println(err)
			return
		}

		info, err := analyzer.AnalyzeProject(root)
		if err != nil {
			fmt.Println(err)
			return
		}

		summary := ai.BuildSummary(info)
		status := explain.BuildStatus(topic, info)

		req, err := explain.BuildRequest(topic, summary, status)
		if err != nil {
			fmt.Println(err)
			return
		}

		client := openai.New(creds.APIKey)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		resp, err := client.Complete(ctx, req)
		if err != nil {
			fmt.Printf("Failed to reach OpenAI: %v\n", err)
			return
		}

		result, err := explain.Parse(resp.Content)
		if err != nil {
			fmt.Println(err)
			return
		}

		renderExplainResult(topic, result, status)
	},
}

func renderExplainResult(topic string, result *explain.Result, status []explain.StatusItem) {

	name := explain.DisplayName(topic)

	ui.Header("💡 " + name + " Explained")
	ui.Blank()

	ui.Line("Current Project")
	ui.Blank()

	for _, item := range result.CurrentProject {
		printWrapped(" ✓ ", item)
	}

	ui.Blank()
	ui.Line("Build Flow")
	ui.Blank()

	for _, step := range result.BuildFlow {
		printWrapped(" ", step)
	}

	ui.Blank()
	ui.Line(fmt.Sprintf("Why %s?", name))
	ui.Blank()

	for _, reason := range result.WhyTopic {
		printWrapped(" • ", reason)
	}

	ui.Blank()
	ui.Line("Current Status")
	ui.Blank()

	for _, item := range status {

		mark := "✗"
		if item.Present {
			mark = "✓"
		}

		printWrapped(" "+mark+" ", item.Label)
	}

	ui.Footer()
}

// printWrapped : 박스 너비에 맞게 텍스트를 줄바꿈해 출력
func printWrapped(prefix, text string) {

	width := utf8.RuneCountInString(prefix)
	indent := strings.Repeat(" ", width)

	for i, line := range ui.Wrap(text, 60-width) {
		if i == 0 {
			ui.Line(prefix + line)
		} else {
			ui.Line(indent + line)
		}
	}
}

// explainArgs : 첫 번째 인자는 유효한 topic이어야 하고, 두 번째 인자(경로)는
// 선택 — 생략하면 현재 디렉터리(doctor/scan의 [path]와 동일 패턴).
// cobra.OnlyValidArgs는 경로까지 topic으로 검사해버려서 못 씀
func explainArgs(cmd *cobra.Command, args []string) error {

	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("accepts between 1 and 2 arg(s), received %d", len(args))
	}

	for _, topic := range explain.Topics {
		if args[0] == topic {
			return nil
		}
	}

	return fmt.Errorf("invalid topic %q for %q", args[0], cmd.CommandPath())
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
