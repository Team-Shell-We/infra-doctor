package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
	"github.com/Team-Shell-We/infra-doctor/internal/ai/openai"
	"github.com/Team-Shell-We/infra-doctor/internal/ai/recommend"
	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
	"github.com/Team-Shell-We/infra-doctor/internal/ui"
	"github.com/spf13/cobra"
)

var recommendCmd = &cobra.Command{
	Use:   "recommend [path]",
	Short: "Recommend a deployment strategy for your project's scale",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		lang := currentLang()

		root := "."
		if len(args) == 1 {
			root = args[0]
		}

		creds, err := ai.Load()
		if err != nil {

			if errors.Is(err, ai.ErrNotLoggedIn) {
				fmt.Println(i18n.Get(lang, "common.notLoggedIn"))
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
		decision := recommend.Decide(info)

		req, err := recommend.BuildRequest(summary, decision, lang)
		if err != nil {
			fmt.Println(err)
			return
		}

		client := openai.New(creds.APIKey)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		resp, err := client.Complete(ctx, req)
		if err != nil {
			fmt.Printf(i18n.Get(lang, "common.openaiFailed")+"\n", err)
			return
		}

		result, err := recommend.Parse(resp.Content)
		if err != nil {
			fmt.Println(err)
			return
		}

		renderRecommendResult(lang, summary, decision, result, recommend.NextSteps(info, decision))
	},
}

func renderRecommendResult(lang string, summary ai.Summary, decision recommend.Decision, result *recommend.Result, nextSteps []string) {

	ui.Header("🚀 " + i18n.Get(lang, "recommend.title"))
	ui.Blank()

	// Current Stack, Recommended, Kubernetes, Next Step은 전부 스캔 결과에서
	// 바로 나오는 결정론적 값 — AI를 거치지 않는다.

	ui.Line(i18n.Get(lang, "recommend.currentStack"))
	ui.Blank()

	if summary.Framework != "" {
		ui.Line(" ✓ " + summary.Framework)
	}

	for _, db := range summary.Database {
		ui.Line(" ✓ " + db)
	}

	ui.Blank()
	ui.Line(i18n.Get(lang, "recommend.recommended"))
	ui.Blank()
	ui.Line(" ⭐ " + decision.Recommended)

	ui.Blank()
	ui.Line("Kubernetes")
	ui.Blank()

	if decision.KubernetesFit {
		ui.Line(" ✓ " + i18n.Get(lang, "recommend.fit.yes"))
	} else {
		ui.Line(" ✗ " + i18n.Get(lang, "recommend.fit.no"))
	}

	ui.Blank()
	ui.Line(i18n.Get(lang, "recommend.reason"))
	ui.Blank()

	for _, reason := range result.Reasons {
		printWrapped(" • ", reason)
	}

	ui.Blank()
	ui.Line(i18n.Get(lang, "recommend.nextStep"))
	ui.Blank()

	for _, step := range nextSteps {
		ui.Line(" " + step)
	}

	ui.Footer()
}

func init() {
	rootCmd.AddCommand(recommendCmd)
}
