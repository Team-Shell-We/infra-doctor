package cmd

/*
doctor 명령을 정의
사용자로부터 경로인자를 받고, analyzer를 호출해 프로젝트 분석 결과를 받아 터미널에 요약을 출력
출력내용: 프레임워크, 데이터베이스, docker, docker compose 파일, ci/cd, profiles 등
*/
import (
	"fmt"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
	"github.com/Team-Shell-We/infra-doctor/internal/ui"
	"github.com/spf13/cobra"
)

// doctorCheckNameKeys : doctor.Checklist()가 반환하는 고정된 체크명(영어)을
// i18n key로 매핑. internal/doctor 패키지 자체는 건드리지 않고 렌더링
// 시점에만 매핑한다.
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

		ui.Footer()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
