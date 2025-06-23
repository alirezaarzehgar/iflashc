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
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
	"github.com/alirezaarzehgar/iflashc/internal/setproxy"
	"github.com/alirezaarzehgar/iflashc/internal/translate"
	"github.com/alirezaarzehgar/iflashc/internal/ui"
	"github.com/spf13/cobra"
	"github.com/tiagomelo/go-clipboard/clipboard"
	"golang.org/x/net/proxy"

	_ "embed"

	_ "modernc.org/sqlite"
)

var (
	TranslateConfig struct {
		dbPath            string
		noDB              bool
		SchemaDataQueries string
	}

	app struct {
		gui     ui.GUI
		ctx     context.Context
		db      *sql.DB
		queries *query.Queries
		configs config.Config
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iflashc",
	Short: "translate selected text",
	Run: func(cmd *cobra.Command, args []string) {
		defer app.gui.Run()

		c := clipboard.New(clipboard.ClipboardOptions{Primary: true})
		selectedText, err := c.PasteText()
		if err != nil {
			app.gui.ShowError("failed to paste", err)
			return
		}

		selectedText = strings.ToLower(selectedText)

		cfgTranslator := app.configs[config.DefaultKeys.Translator]
		cfgLang := app.configs[config.DefaultKeys.DestLang]
		cfgCtx := app.configs[config.DefaultKeys.Context]
		explaination, err := app.queries.FindMatchedWord(app.ctx, query.FindMatchedWordParams{Word: selectedText, Translator: cfgTranslator, Lang: cfgLang})
		if err == nil {
			app.gui.ShowText(ui.TextBox{Title: selectedText, Text: explaination})
			return
		}

		translator := translate.New(config.TransType(cfgTranslator), app.configs)
		explaination, err = translator.Translate(selectedText)
		if err != nil {
			app.gui.ShowError("failed to translate selected text", err)
			return
		}

		if !TranslateConfig.noDB {
			err = app.queries.SaveWord(app.ctx, query.SaveWordParams{
				Word:       selectedText,
				Exp:        explaination,
				Translator: cfgTranslator,
				Lang:       cfgLang,
				Context:    cfgCtx,
			})
			if err != nil {
				app.gui.ShowError("failed to save explanation", err)
				return
			}
		}

		app.gui.ShowText(ui.TextBox{Title: selectedText, Text: explaination})
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
	rootCmd.PersistentFlags().BoolVar(&TranslateConfig.noDB, "nodb", false, "disable database actions and operate using default values")

	var err error
	app.gui = ui.NewGUI()

	app.ctx = context.Background()
	app.db, err = sql.Open("sqlite", TranslateConfig.dbPath)
	if err != nil {
		app.gui.ShowError("failed to open local database", err)
		return
	}

	if _, err = os.Stat(TranslateConfig.dbPath); os.IsNotExist(err) {
		schema, err := config.GetSchema()
		if err != nil {
			app.gui.ShowError("failed to generate default config", err)
			return
		}
		_, err = app.db.ExecContext(app.ctx, schema)
		if err != nil {
			app.gui.ShowError("failed to migrate local database", err)
			return
		}
	}

	app.queries = query.New(app.db)
	kv, _ := app.queries.GetConfigs(app.ctx)
	app.configs = config.ConfigToMap(kv)

	if v, ok := app.configs[config.DefaultKeys.Socks5]; ok && v != "" {
		httpClient, err := setproxy.NewSocks5Client(v, &proxy.Auth{
			User:     app.configs[config.DefaultKeys.Socks5Username],
			Password: app.configs[config.DefaultKeys.Socks5Password],
		})
		if err != nil {
			app.gui.ShowError("failed connecting to proxy", err)
			return
		}
		http.DefaultClient = httpClient
	}
}
