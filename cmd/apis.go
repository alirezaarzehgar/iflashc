/*
Copyright © 2025 Alireza Arzehgar <alirezaarzehgar82@gmail.com>
*/
package cmd

import (
	"net/http"

	"fyne.io/fyne/v2"
	"github.com/alirezaarzehgar/iflashc/internal/gui"
	"github.com/alirezaarzehgar/iflashc/internal/setproxy"
	"github.com/alirezaarzehgar/iflashc/internal/translate"
	"github.com/spf13/cobra"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

var (
	apisConfig struct {
		destLang     string
		groqApiKey   string
		groqLlmModel string
		groqSocks5   string
	}
)

func generalTranslate(translator translate.Translator) {
	c := clipboard.New(clipboard.ClipboardOptions{Primary: true})
	selectedText, err := c.PasteText()
	if err != nil {
		gui.ShowText(gui.TextBox{Title: "failed to paste", Text: err.Error()})
		return
	}

	if apisConfig.groqSocks5 != "" {
		client, err := setproxy.NewSocks5Client(apisConfig.groqSocks5, nil)
		if err == nil {
			http.DefaultClient = client
		}
	}

	response, err := translator.Translate(selectedText)
	if err != nil {
		gui.ShowText(gui.TextBox{Title: "failed to translate", Text: err.Error()})
		return
	}

	err = gui.ShowText(gui.TextBox{Title: selectedText, Text: response, HaveBtn: true})
	if err != nil {
		gui.ShowText(gui.TextBox{Title: "unable showing on text box", Text: err.Error()})
	}
}

var apisCmd = &cobra.Command{
	Use:   "apis",
	Short: "translate selected text based on customized apis",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: apisConfig.destLang}
		generalTranslate(translate.New(translate.TypeGoogle, cfg))
	},
}

var groqCmd = &cobra.Command{
	Use:   "groq",
	Short: "translate using groq",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: apisConfig.destLang, ApiKey: apisConfig.groqApiKey, LLMmodel: apisConfig.groqLlmModel}
		generalTranslate(translate.New(translate.TypeGroq, cfg))
	},
}

var googleCmd = &cobra.Command{
	Use:   "google",
	Short: "translate using google translate",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: apisConfig.destLang}
		generalTranslate(translate.New(translate.TypeGoogle, cfg))
	},
}

var dictapi = &cobra.Command{
	Use:   "dictapi",
	Short: "translate using dictionaryapi.dev",
	Run: func(cmd *cobra.Command, args []string) {
		generalTranslate(translate.New(translate.TypeDictionaryApi, translate.TranslatorConfig{}))
	},
}

func generalAnalyzer(translator translate.Translator) {
	c := clipboard.New(clipboard.ClipboardOptions{Primary: true})
	selectedText, err := c.PasteText()
	if err != nil {
		gui.ShowText(gui.TextBox{Title: "failed to paste", Text: err.Error()})
		return
	}

	if apisConfig.groqSocks5 != "" {
		client, err := setproxy.NewSocks5Client(apisConfig.groqSocks5, nil)
		if err == nil {
			http.DefaultClient = client
		}
	}

	response, err := translator.Translate(selectedText)
	if err != nil {
		gui.ShowText(gui.TextBox{Title: "failed to translate", Text: err.Error()})
		return
	}

	gui.DefaultWindowSize = fyne.NewSize(800, 0)
	err = gui.ShowText(gui.TextBox{Title: "Analyzer", Text: response})
	if err != nil {
		gui.ShowText(gui.TextBox{Title: "unable showing on text box", Text: err.Error()})
	}
}

var groqAnalyzeCmd = &cobra.Command{
	Use:   "groq-analyze",
	Short: "analyze text using groq",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{ApiKey: apisConfig.groqApiKey, LLMmodel: apisConfig.groqLlmModel}
		generalAnalyzer(translate.New(translate.TypeGroqAlayzer, cfg))
	},
}

func init() {
	rootCmd.AddCommand(apisCmd)
	apisCmd.AddCommand(groqCmd, groqAnalyzeCmd, googleCmd, dictapi)

	apisCmd.PersistentFlags().StringVar(&apisConfig.groqLlmModel, "socks5", "", "Socks5 proxy for all requests")
	apisCmd.PersistentFlags().StringVar(&apisConfig.destLang, "dest-lang", "fa", "Destination language")

	groqCmd.PersistentFlags().StringVar(&apisConfig.groqApiKey, "api-key", "", "API Key for groq")
	groqCmd.PersistentFlags().StringVar(&apisConfig.groqLlmModel, "llm-model", "", "LLM Model name for groq")
	groqCmd.MarkPersistentFlagRequired("api-key")
	groqCmd.MarkPersistentFlagRequired("llm-model")

	groqAnalyzeCmd.PersistentFlags().StringVar(&apisConfig.groqApiKey, "api-key", "", "API Key for groq")
	groqAnalyzeCmd.PersistentFlags().StringVar(&apisConfig.groqLlmModel, "llm-model", "", "LLM Model name for groq")
	groqAnalyzeCmd.MarkPersistentFlagRequired("api-key")
	groqAnalyzeCmd.MarkPersistentFlagRequired("llm-model")
}
