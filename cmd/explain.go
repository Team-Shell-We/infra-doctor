package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
	"github.com/Team-Shell-We/infra-doctor/internal/ai/explain"
	"github.com/Team-Shell-We/infra-doctor/internal/ai/openai"
	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
	"github.com/Team-Shell-We/infra-doctor/internal/ui"
	"github.com/spf13/cobra"
)

// explainStatusLabelKeys : BuildStatus()가 반환하는 영어 라벨을 렌더링 시점에 i18n key로 매핑한다
var explainStatusLabelKeys = map[string]string{
	"Dockerfile":                             "explain.status.dockerfile",
	"Docker Compose":                         "explain.status.dockerCompose",
	"Health Check":                           "explain.status.healthCheck",
	"docker-compose.yml":                     "explain.status.dockerComposeYml",
	"GitHub Actions workflow":                "explain.status.githubActionsWorkflow",
	"Kubernetes manifests":                   "explain.status.kubernetesManifests",
	"Nginx configuration":                    "explain.status.nginxConfig",
	"PostgreSQL":                             "explain.status.postgresql",
	"AWS SDK dependency":                     "explain.status.awsSdk",
	"Relational database (PostgreSQL/MySQL)": "explain.status.relationalDb",
	"Redis":                                  "explain.status.redis",
}

var explainCmd = &cobra.Command{
	Use:       "explain <topic> [path]",
	Short:     "Explain an infrastructure concept in the context of your project",
	Args:      explainArgs,
	ValidArgs: explain.Topics,

	Run: func(cmd *cobra.Command, args []string) {

		lang := currentLang()
		topic := args[0]

		root := "."
		if len(args) == 2 {
			root = args[1]
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
		status := explain.BuildStatus(topic, info)

		req, err := explain.BuildRequest(topic, summary, status, lang)
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

		result, err := explain.Parse(resp.Content)
		if err != nil {
			fmt.Println(err)
			return
		}

		renderExplainResult(lang, topic, result, status)
	},
}

func renderExplainResult(lang, topic string, result *explain.Result, status []explain.StatusItem) {

	name := explain.DisplayName(topic)

	ui.Header("💡 " + name + " " + i18n.Get(lang, "explain.suffix"))
	ui.Blank()

	// topic이 프로젝트에 없으면(예: k8s) 이후 내용이 가정 시나리오임을 코드가 먼저 알린다 — AI 응답 톤에만 의존하면 가정법이 항상 지켜진다는 보장이 없다
	if !explain.TopicPresent(status) {
		printWrapped(" ⚠ ", i18n.Get(lang, "explain.notAdopted"))
		ui.Blank()
	}

	// "현재 상태"(결정론적 ✓/✗)를 AI 서술보다 먼저 보여준다 — 순서가 반대면 미도입 기술(예: k8s)도 이미 쓰는 것처럼 오인할 수 있다
	ui.Line(i18n.Get(lang, "explain.currentStatus"))
	ui.Blank()

	for _, item := range status {

		mark := "✗"
		if item.Present {
			mark = "✓"
		}

		label := item.Label
		if key, ok := explainStatusLabelKeys[item.Label]; ok {
			label = i18n.Get(lang, key)
		}

		printWrapped(" "+mark+" ", label)
	}

	ui.Blank()
	ui.Line(i18n.Get(lang, "explain.currentProject"))
	ui.Blank()

	for _, item := range result.CurrentProject {
		printWrapped(" ✓ ", item)
	}

	ui.Blank()
	ui.Line(i18n.Get(lang, "explain.buildFlow"))
	ui.Blank()

	for _, step := range result.BuildFlow {
		printWrapped(" ", step)
	}

	ui.Blank()
	ui.Line(fmt.Sprintf(i18n.Get(lang, "explain.whyTopic"), name))
	ui.Blank()

	for _, reason := range result.WhyTopic {
		printWrapped(" • ", reason)
	}

	ui.Footer()
}

// printWrapped : 박스 너비에 맞게 텍스트를 줄바꿈해 출력
func printWrapped(prefix, text string) {

	width := ui.DisplayWidth(prefix)
	indent := strings.Repeat(" ", width)

	for i, line := range ui.Wrap(text, 60-width) {
		if i == 0 {
			ui.Line(prefix + line)
		} else {
			ui.Line(indent + line)
		}
	}
}

// explainArgs : 1번째 인자는 유효한 topic, 2번째(경로)는 생략 가능(기본값 현재 디렉터리, doctor/scan의 [path]와 동일). cobra.OnlyValidArgs는 경로도 topic으로 검사해 못 쓴다
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
