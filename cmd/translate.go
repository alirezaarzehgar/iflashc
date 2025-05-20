/*
Copyright © 2025 Alireza Arzehgar <alirezaarzehgar82@gmail.com>
*/
package cmd

import (
	"log"
	"os"

	"github.com/alirezaarzehgar/iflashc/internal/gui"
	"github.com/alirezaarzehgar/iflashc/internal/translate"
	"github.com/spf13/cobra"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

var transType string

// translateCmd represents the translate command
var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "translate selected text",
	Run: func(cmd *cobra.Command, args []string) {
		c := clipboard.New(clipboard.ClipboardOptions{Primary: true})
		selectedText, err := c.PasteText()
		if err != nil {
			log.Fatal("unable copying selected text", err)
		}

		cfg := translate.TranslatorConfig{To: "fa", ApiKey: os.Getenv("GROK_API_KEY"), LLMmodel: os.Getenv("GROK_LLM_MODEL")}
		response, err := translate.New(transType, cfg).Translate(selectedText)
		if err != nil {
			log.Fatal("unable translating the word:", err)
		}

		err = gui.ShowWord(selectedText, response)
		if err != nil {
			log.Fatal("unable show message:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.PersistentFlags().StringVar(&transType, "tt", "google", "Set translation type")
}
