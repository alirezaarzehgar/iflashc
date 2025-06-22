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
	"fmt"

	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
	"github.com/spf13/cobra"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "manage notes and bookmarks",
	Run: func(cmd *cobra.Command, args []string) {
		c := clipboard.New(clipboard.ClipboardOptions{Primary: true})
		selectedText, err := c.PasteText()
		if err != nil {
			app.gui.ShowError("failed to paste", err)
			app.gui.Run()
			return
		}

		app.queries.SaveNote(app.ctx, query.SaveNoteParams{
			Note:    selectedText,
			Comment: "",
			Context: app.configs[config.DefaultKeys.Context],
		})
		if err != nil {
			app.gui.ShowError("failed to save note", err)
			app.gui.Run()
			return
		}

		fmt.Println("Saved note")
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
