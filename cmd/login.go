package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
	"github.com/Team-Shell-We/infra-doctor/internal/ai/openai"
	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
	"github.com/Team-Shell-We/infra-doctor/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to enable AI-powered commands",
	Args:  cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		runLogin(cmd.InOrStdin(), cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(in io.Reader, out io.Writer) {

	lang := currentLang()
	reader := bufio.NewReader(in)

	ui.Header("🔐 " + i18n.Get(lang, "login.title"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, i18n.Get(lang, "login.selectMethod"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "  "+i18n.Get(lang, "login.option1"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "  "+i18n.Get(lang, "login.option2"))
	fmt.Fprintln(out)
	fmt.Fprint(out, i18n.Get(lang, "login.choose"))

	switch readLine(reader) {

	case "1":
		loginWithOpenAI(reader, out, lang)

	case "2":
		fmt.Fprintln(out)
		fmt.Fprintln(out, i18n.Get(lang, "login.accountUnavailable"))

	default:
		fmt.Fprintln(out)
		fmt.Fprintln(out, i18n.Get(lang, "login.invalidChoice"))
	}
}

func loginWithOpenAI(reader *bufio.Reader, out io.Writer, lang string) {

	fmt.Fprintln(out)
	fmt.Fprintln(out, i18n.Get(lang, "login.apiKeyPrompt"))
	fmt.Fprint(out, "> ")

	apiKey := readAPIKey(reader)
	if apiKey == "" {
		fmt.Fprintln(out)
		fmt.Fprintln(out, i18n.Get(lang, "login.apiKeyEmpty"))
		return
	}

	client := openai.New(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.VerifyCredentials(ctx); err != nil {

		fmt.Fprintln(out)

		if errors.Is(err, ai.ErrInvalidAPIKey) {
			fmt.Fprintln(out, "❌ "+i18n.Get(lang, "login.invalidKey"))
			fmt.Fprintln(out)
			fmt.Fprintln(out, i18n.Get(lang, "login.checkKey"))
			return
		}

		fmt.Fprintf(out, i18n.Get(lang, "login.verifyFailed")+"\n", err)
		return
	}

	creds := ai.LoadOrDefault()
	creds.Provider = "openai"
	creds.APIKey = apiKey
	creds.Login = true

	if err := ai.Save(creds); err != nil {
		fmt.Fprintln(out)
		fmt.Fprintf(out, i18n.Get(lang, "login.saveFailed")+"\n", err)
		return
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "✅ "+i18n.Get(lang, "login.verified"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "✅ "+i18n.Get(lang, "login.completed"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, i18n.Get(lang, "login.provider"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "OpenAI")
}

// readAPIKey : stdin이 실제 대화형 터미널일 때만 입력을 마스킹
// 그 외엔 일반 버퍼 읽기로 대체(파이프 입력으로도 스크립팅/테스트 가능)
func readAPIKey(reader *bufio.Reader) string {

	if term.IsTerminal(int(os.Stdin.Fd())) {

		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		if err == nil {
			return strings.TrimSpace(string(bytePassword))
		}
	}

	return readLine(reader)
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
