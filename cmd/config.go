package cmd

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
	"github.com/spf13/cobra"
)

func configCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show or change CLI configuration",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {

			creds := ai.LoadOrDefault()

			if lang, _ := cmd.Flags().GetString("lang"); lang != "" {

				if !i18n.IsSupported(lang) {
					fmt.Printf("Unsupported language %q (supported: %v)\n", lang, i18n.Supported())
					return
				}

				creds.Language = lang

				if err := ai.Save(creds); err != nil {
					fmt.Println(err)
					return
				}
			}

			printConfig(creds)
		},
	}

	cmd.Flags().String("lang", "", "Set output language (en, ko)")

	return cmd
}

func printConfig(creds ai.Credentials) {

	lang := creds.Language
	if lang == "" {
		lang = i18n.English
	}

	provider := creds.Provider
	if provider == "" {
		provider = i18n.Get(lang, "config.unconfigured")
	}

	output := creds.OutputFormat
	if output == "" {
		output = "ASCII + Mermaid"
	}

	autoExport := i18n.Get(lang, "config.disabled")
	if creds.AutoExport {
		autoExport = i18n.Get(lang, "config.enabled")
	}

	fmt.Println(i18n.Get(lang, "config.llm"))
	fmt.Println(provider)
	fmt.Println()
	fmt.Println(i18n.Get(lang, "config.language"))
	fmt.Println(lang)
	fmt.Println()
	fmt.Println(i18n.Get(lang, "config.output"))
	fmt.Println(output)
	fmt.Println()
	fmt.Println(i18n.Get(lang, "config.autoExport"))
	fmt.Println(autoExport)
}

func init() {
	rootCmd.AddCommand(configCommand())
}
