package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
)

func (g GUI) Dashboard(q *query.Queries, cfgs config.Config) {
	g.win.Resize(fyne.NewSize(DefaultDictWindowLen, DefaultDictWindowLen))

	historyPage := &fyne.Container{}
	notesPage := &fyne.Container{}

	pageSelector := widget.NewSelect([]string{"history", "notes"}, func(s string) {
		switch s {
		case "history":
			g.setContent(historyPage)
		case "notes":
			g.setContent(notesPage)
		}
	})
	pageSelector.SetSelected("history")

	historyPage = g.historyPage(q, cfgs, pageSelector)
	notesPage = g.notesPage(q, cfgs, pageSelector)
	g.setContent(historyPage)
	g.win.Show()
}
