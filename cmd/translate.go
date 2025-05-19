/*
Copyright © 2025 Alireza Arzehgar <alirezaarzehgar82@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/alirezaarzehgar/iflashc/internal/gui"
	"github.com/alirezaarzehgar/iflashc/internal/translate"
	"github.com/spf13/cobra"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

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

		response, err := translate.New(translate.TypeGrok, "fa").Translate(selectedText)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// translateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// translateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
