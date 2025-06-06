/*
Copyright © 2025 Alireza Arzehgar <alirezaarzehgar82@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
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
	apiConfig struct {
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

	if apiConfig.groqSocks5 != "" {
		client, err := setproxy.NewSocks5Client(apiConfig.groqSocks5, nil)
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

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "translate selected text based on customized apis",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: apiConfig.destLang}
		generalTranslate(translate.New(translate.TypeGoogle, cfg))
	},
}

var groqCmd = &cobra.Command{
	Use:   "groq",
	Short: "translate using groq",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: apiConfig.destLang, ApiKey: apiConfig.groqApiKey, LLMmodel: apiConfig.groqLlmModel}
		generalTranslate(translate.New(translate.TypeGroq, cfg))
	},
}

var googleCmd = &cobra.Command{
	Use:   "google",
	Short: "translate using google translate",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := translate.TranslatorConfig{To: apiConfig.destLang}
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

	if apiConfig.groqSocks5 != "" {
		client, err := setproxy.NewSocks5Client(apiConfig.groqSocks5, nil)
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
		cfg := translate.TranslatorConfig{ApiKey: apiConfig.groqApiKey, LLMmodel: apiConfig.groqLlmModel}
		generalAnalyzer(translate.New(translate.TypeGroqAlayzer, cfg))
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.AddCommand(groqCmd, groqAnalyzeCmd, googleCmd, dictapi)

	apiCmd.PersistentFlags().StringVar(&apiConfig.groqLlmModel, "socks5", "", "Socks5 proxy for all requests")
	apiCmd.PersistentFlags().StringVar(&apiConfig.destLang, "dest-lang", "fa", "Destination language")

	groqCmd.PersistentFlags().StringVar(&apiConfig.groqApiKey, "api-key", "", "API Key for groq")
	groqCmd.PersistentFlags().StringVar(&apiConfig.groqLlmModel, "llm-model", "", "LLM Model name for groq")
	groqCmd.MarkPersistentFlagRequired("api-key")
	groqCmd.MarkPersistentFlagRequired("llm-model")

	groqAnalyzeCmd.PersistentFlags().StringVar(&apiConfig.groqApiKey, "api-key", "", "API Key for groq")
	groqAnalyzeCmd.PersistentFlags().StringVar(&apiConfig.groqLlmModel, "llm-model", "", "LLM Model name for groq")
	groqAnalyzeCmd.MarkPersistentFlagRequired("api-key")
	groqAnalyzeCmd.MarkPersistentFlagRequired("llm-model")
}
