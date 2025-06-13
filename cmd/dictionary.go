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

	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
	"github.com/alirezaarzehgar/iflashc/internal/ui"
	"github.com/spf13/cobra"
)

var dictionaryCmd = &cobra.Command{
	Use:   "dictionary",
	Short: "visit dictionary",
	Run: func(cmd *cobra.Command, args []string) {
		gui := ui.NewGUI()
		defer gui.Run()

		ctx := context.Background()
		db, err := sql.Open("sqlite", TranslateConfig.dbPath)
		if err != nil {
			gui.ShowError("failed to open local database", err)
			return
		}
		defer db.Close()

		if _, err := os.Stat(TranslateConfig.dbPath); os.IsNotExist(err) {
			schema, err := config.GetSchema()
			if err != nil {
				gui.ShowError("failed to generate default config", err)
				return
			}
			_, err = db.ExecContext(ctx, schema)
			if err != nil {
				gui.ShowError("failed to migrate local database", err)
				return
			}
		}

		q := query.New(db)
		kv, _ := q.GetConfigs(ctx)
		gui.Dashboard(
			q,
			config.ConfigToMap(kv),
		)
	},
}

func init() {
	rootCmd.AddCommand(dictionaryCmd)
}
