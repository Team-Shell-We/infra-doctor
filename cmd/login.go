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

	reader := bufio.NewReader(in)

	ui.Header("🔐 Login")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Select Login Method")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "  1. OpenAI API Key")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "  2. Infra Doctor Account")
	fmt.Fprintln(out)
	fmt.Fprint(out, "Choose (1-2): ")

	switch readLine(reader) {

	case "1":
		loginWithOpenAI(reader, out)

	case "2":
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Infra Doctor Account login is not available yet. Please use an OpenAI API Key for now.")

	default:
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Invalid choice. Run 'infra-doctor login' again and choose 1 or 2.")
	}
}

func loginWithOpenAI(reader *bufio.Reader, out io.Writer) {

	fmt.Fprintln(out)
	fmt.Fprintln(out, "OpenAI API Key")
	fmt.Fprint(out, "> ")

	apiKey := readAPIKey(reader)
	if apiKey == "" {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "API Key cannot be empty.")
		return
	}

	client := openai.New(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.VerifyCredentials(ctx); err != nil {

		fmt.Fprintln(out)

		if errors.Is(err, ai.ErrInvalidAPIKey) {
			fmt.Fprintln(out, "❌ Invalid OpenAI API Key.")
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Please check your API Key.")
			return
		}

		fmt.Fprintf(out, "Failed to verify API Key: %v\n", err)
		return
	}

	creds := ai.Credentials{Provider: "openai", APIKey: apiKey, Login: true}

	if err := ai.Save(creds); err != nil {
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Failed to save credentials: %v\n", err)
		return
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "✅ OpenAI API Key verified.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "✅ Login completed.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Provider")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "OpenAI")
}

// readAPIKey masks input when stdin is a real interactive terminal
// (matching the spec's `sk-********` masked example); otherwise it falls
// back to plain buffered reading so login stays scriptable/testable via
// piped input.
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
