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
	"context"
	"database/sql"
	"os"
	"path"
	"strings"

	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
	"github.com/alirezaarzehgar/iflashc/internal/translate"
	"github.com/alirezaarzehgar/iflashc/internal/ui"
	"github.com/spf13/cobra"
	"github.com/tiagomelo/go-clipboard/clipboard"

	_ "embed"

	_ "modernc.org/sqlite"
)

var (
	TranslateConfig struct {
		dbPath            string
		noDB              bool
		SchemaDataQueries string
	}

	DefaultConfigs = config.Defaults{
		Translator: translate.TypeDictionaryApi,
		DestLang:   "fa",
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iflashc",
	Short: "translate selected text",
	Run: func(cmd *cobra.Command, args []string) {
		gui := ui.NewGUI()
		defer gui.Run()

		c := clipboard.New(clipboard.ClipboardOptions{Primary: true})
		selectedText, err := c.PasteText()
		if err != nil {
			gui.ShowText(ui.TextBox{Title: "failed to paste", Text: err.Error()})
			return
		}

		ctx := context.Background()
		db, err := sql.Open("sqlite", TranslateConfig.dbPath)
		if err != nil {
			gui.ShowText(ui.TextBox{Title: "failed to open local database", Text: err.Error()})
			return
		}
		defer db.Close()

		if _, err := os.Stat(TranslateConfig.dbPath); os.IsNotExist(err) {
			schema, err := config.GetSchema(DefaultConfigs)
			if err != nil {
				gui.ShowText(ui.TextBox{Title: "failed to generate default config", Text: err.Error()})
				return
			}
			_, err = db.ExecContext(ctx, schema)
			if err != nil {
				gui.ShowText(ui.TextBox{Title: "failed to migrate local database", Text: err.Error()})
				return
			}
		}

		q := query.New(db)

		kv, _ := q.GetConfigs(ctx)
		configs := config.ConfigToMap(kv)

		selectedText = strings.ToLower(selectedText)

		cfgTranslator := configs[config.DefaultKeys.Translator]
		cfgLang := configs[config.DefaultKeys.DestLang]
		cfgCtx := configs[config.DefaultKeys.Context]
		explaination, err := q.FindMatchedWord(ctx, query.FindMatchedWordParams{Word: selectedText, Translator: cfgTranslator, Lang: cfgLang})
		if err == nil {
			gui.ShowText(ui.TextBox{Title: selectedText, Text: explaination})
			return
		}

		cfg := translate.TranslatorConfig{
			To:        configs[config.DefaultKeys.DestLang],
			GroqModel: configs[config.DefaultKeys.GroqModel],
			ApiKey:    configs[config.DefaultKeys.GroqApiKey],
		}
		translator := translate.New(translate.TransType(cfgTranslator), cfg)
		explaination, err = translator.Translate(selectedText)
		if err != nil {
			gui.ShowText(ui.TextBox{Title: "failed to translate selected text", Text: err.Error()})
			return
		}

		err = q.SaveWord(ctx, query.SaveWordParams{
			Word:       selectedText,
			Exp:        explaination,
			Translator: cfgTranslator,
			Lang:       cfgLang,
			Context:    sql.NullString{String: cfgCtx, Valid: len(cfgCtx) > 0},
		})
		if err != nil {
			gui.ShowText(ui.TextBox{Title: "failed to save explanation", Text: err.Error()})
			return
		}

		gui.ShowText(ui.TextBox{Title: selectedText, Text: explaination})
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&TranslateConfig.dbPath, "db", path.Join(os.Getenv("HOME"), ".iflashc.db"), "local database path")
	// rootCmd.PersistentFlags().BoolVar(&translateConfig.noDB, "nodb", false, "disable database actions and operate using default values")
}
