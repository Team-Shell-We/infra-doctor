package cmd

// doctor 명령 : 경로를 분석해 준비도 점수·체크리스트(Docker/Compose/헬스체크/리버스프록시/모니터링/로그 로테이션/DB 백업)·추천·다음 단계를 출력한다

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
	"github.com/Team-Shell-We/infra-doctor/internal/nextstep"
	"github.com/Team-Shell-We/infra-doctor/internal/project"
	"github.com/Team-Shell-We/infra-doctor/internal/ui"
	"github.com/spf13/cobra"
)

var (
	doctorJSON      bool
	doctorFailUnder int
)

// doctorCheckNameKeys : Checklist()가 반환하는 영어 체크명을 렌더링 시점에 i18n key로 매핑한다(internal/doctor는 그대로 둠)
var doctorCheckNameKeys = map[string]string{
	"Docker":         "doctor.check.docker",
	"Docker Compose": "doctor.check.dockerCompose",
	"Health Check":   "doctor.check.healthCheck",
	"Reverse Proxy":  "doctor.check.reverseProxy",
	"Monitoring":     "doctor.check.monitoring",
	"Log Rotation":   "doctor.check.logRotation",
	"DB Backup":      "doctor.check.dbBackup",
}

var doctorCmd = &cobra.Command{
	Use:   "doctor [path]",
	Short: "Analyze project and provide recommendations",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		lang := currentLang()

		root := "."

		if len(args) == 1 {
			root = args[0]
		}

		info, err := analyzer.AnalyzeProject(root)
		if err != nil {
			fmt.Println(err)
			return
		}

		result := doctor.Analyze(info)

		if doctorJSON {
			if err := writeDoctorJSON(os.Stdout, result); err != nil {
				fmt.Println(err)
				return
			}
		} else {
			printDoctorBox(lang, info, result)
		}

		if doctorShouldFail(result.Score, doctorFailUnder, cmd.Flags().Changed("fail-under")) {
			os.Exit(1)
		}
	},
}

func writeDoctorJSON(w io.Writer, result *doctor.Result) error {

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, string(data))
	return err
}

// doctorShouldFail : --fail-under 지정 시에만 게이트를 켠다. 예: 미지정→false, "--fail-under 80"+score 75→true. 안 그러면 기본값 0 때문에 항상 게이트가 걸린다
func doctorShouldFail(score, failUnder int, failUnderSet bool) bool {
	return failUnderSet && score < failUnder
}

func printDoctorBox(lang string, info *project.Info, result *doctor.Result) {

	ui.Header("🩺 " + i18n.Get(lang, "doctor.title"))
	ui.Blank()

	ui.Line(i18n.Get(lang, "doctor.readiness"))
	ui.Blank()
	ui.Line(fmt.Sprintf(" %s %d%%", ui.ProgressBar(result.Score, 30), result.Score))
	ui.Blank()

	ui.Line(i18n.Get(lang, "doctor.infraCheck"))
	ui.Blank()

	for _, check := range result.Checks {

		mark := "✗"
		if check.Passed {
			mark = "✓"
		}

		name := check.Name
		if key, ok := doctorCheckNameKeys[check.Name]; ok {
			name = i18n.Get(lang, key)
		}

		ui.Line(fmt.Sprintf(" %s %s", mark, name))
	}

	ui.Blank()

	ui.Line(i18n.Get(lang, "doctor.recommendation"))
	ui.Blank()

	if len(result.Diagnoses) == 0 {

		ui.Line(" ✓ " + i18n.Get(lang, "doctor.noIssues"))

	} else {

		for _, d := range result.Diagnoses {

			wrapped := ui.Wrap(strings.TrimSpace(d.Fix), 56)

			for i, line := range wrapped {

				if i == 0 {
					ui.Line(fmt.Sprintf(" • %s", line))
				} else {
					ui.Line(fmt.Sprintf("   %s", line))
				}
			}
		}
	}

	ui.Blank()

	ui.Line(i18n.Get(lang, "doctor.nextStep"))
	ui.Blank()

	for _, step := range nextstep.Suggest(info, false) {
		ui.Line(" " + step)
	}

	ui.Footer()
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "output results as JSON")
	doctorCmd.Flags().IntVar(&doctorFailUnder, "fail-under", 0, "exit non-zero if the score is below this threshold")
	rootCmd.AddCommand(doctorCmd)
}
