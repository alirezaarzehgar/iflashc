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

var grokApiKey, grokLlmModel string

var grokCmd = &cobra.Command{
	Use:   "grok",
	Short: "translate using grok",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: "fa", ApiKey: grokApiKey, LLMmodel: grokLlmModel}
		generalTranslate(translate.New("grok", cfg))
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
	translateCmd.AddCommand(grokCmd, googleCmd)

	grokCmd.PersistentFlags().StringVar(&grokApiKey, "api-key", "", "API Key for Grok")
	grokCmd.PersistentFlags().StringVar(&grokLlmModel, "llm-model", "", "LLM Model name for Grok")
	grokCmd.MarkPersistentFlagRequired("api-key")
	grokCmd.MarkPersistentFlagRequired("llm-model")
}
