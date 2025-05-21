/*
Copyright © 2025 Alireza Arzehgar <alirezaarzehgar82@gmail.com>
*/
package cmd

import (
	"github.com/alirezaarzehgar/iflashc/internal/gui"
	"github.com/alirezaarzehgar/iflashc/internal/translate"
	"github.com/spf13/cobra"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

var transType string

func generalTranslate(translator translate.Translator) {
	c := clipboard.New(clipboard.ClipboardOptions{Primary: true})
	selectedText, err := c.PasteText()
	if err != nil {
		gui.ShowWord("ERROR", "**unable copying the text**: "+err.Error())
		return
	}

	response, err := translator.Translate(selectedText)
	if err != nil {
		gui.ShowWord("ERROR", err.Error())
		return
	}

	err = gui.ShowWord(selectedText, response)
	if err != nil {
		gui.ShowWord("ERROR", err.Error())
		return
	}
}

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "translate selected text",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: "fa"}
		generalTranslate(translate.New("google", cfg))
	},
}

var groqApiKey, groqLlmModel string

var groqCmd = &cobra.Command{
	Use:   "groq",
	Short: "translate using groq",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: "fa", ApiKey: groqApiKey, LLMmodel: groqLlmModel}
		generalTranslate(translate.New("groq", cfg))
	},
}

var googleCmd = &cobra.Command{
	Use:   "google",
	Short: "translate using google translate",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: "fa"}
		generalTranslate(translate.New("google", cfg))
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.PersistentFlags().StringVar(&transType, "tt", "google", "Set translation type")
	translateCmd.AddCommand(groqCmd, googleCmd)

	groqCmd.PersistentFlags().StringVar(&groqApiKey, "api-key", "", "API Key for groq")
	groqCmd.PersistentFlags().StringVar(&groqLlmModel, "llm-model", "", "LLM Model name for groq")
	groqCmd.MarkPersistentFlagRequired("api-key")
	groqCmd.MarkPersistentFlagRequired("llm-model")
}
