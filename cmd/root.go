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
	"fmt"
	"os"
	"path"

	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/gui"
	"github.com/alirezaarzehgar/iflashc/internal/query"
	"github.com/alirezaarzehgar/iflashc/internal/translate"
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
		// Get selected text
		c := clipboard.New(clipboard.ClipboardOptions{Primary: true})
		selectedText, err := c.PasteText()
		if err != nil {
			gui.ShowText(gui.TextBox{Title: "failed to paste", Text: err.Error()})
			return
		}
		_ = selectedText

		// Check & open database
		ctx := context.Background()
		db, err := sql.Open("sqlite", TranslateConfig.dbPath)
		if err != nil {
			gui.ShowText(gui.TextBox{Title: "failed to open local database", Text: err.Error()})
			return
		}
		defer db.Close()

		// Migrate if needed
		if _, err := os.Stat(TranslateConfig.dbPath); os.IsNotExist(err) {
			schema, err := config.GetSchema(DefaultConfigs)
			if err != nil {
				gui.ShowText(gui.TextBox{Title: "failed to generate default config", Text: err.Error()})
				return
			}
			_, err = db.ExecContext(ctx, schema)
			if err != nil {
				gui.ShowText(gui.TextBox{Title: "failed to migrate local database", Text: err.Error()})
				return
			}
			fmt.Println("migrate")
		}

		// init sqlc queries
		q := query.New(db)

		// Get configuration
		kv, _ := q.GetConfigs(ctx)
		conf := config.ConfigToMap(kv)
		fmt.Println(conf)

		// Convert text to lower case
		// Search in db and show first translator result on GUI; EXIT
		// Translate text using all apis in parallel
		// Save each response to database
		// Block for configured translator response
		// show first translator result on GUI

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
