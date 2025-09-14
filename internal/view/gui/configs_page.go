package gui

import (
	"context"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
)

type keyEntry map[string]*widget.Entry

func (g GUI) ManageConfigs(q *query.Queries, cfgs config.Config) {
	hboxConfig := container.NewVBox()

	label := canvas.NewText("Configuration Manager", color.White)
	label.TextSize = DefaultTitleSize
	label.Alignment = fyne.TextAlignCenter
	hboxConfig.Add(label)

	keyEntries := keyEntry{}

	for _, k := range config.ConfigurableKeys {
		entry := widget.NewEntry()
		entry.Text = cfgs[k]
		label := widget.NewLabel(k)
		hbox := container.NewGridWithColumns(2, label, entry)
		hboxConfig.Add(hbox)
		keyEntries[k] = entry
	}

	hboxConfig.Add(widget.NewButton("Save & Exit", func() {
		for _, k := range config.ConfigurableKeys {
			e := keyEntries[k]
			if e.Text == cfgs[k] {
				continue
			}

			ctx := context.Background()
			err := q.ChangeConfig(ctx, query.ChangeConfigParams{Key: k, Value: e.Text})
			if err != nil {
				g.ShowError("failed to change config", err)
			}
		}

		g.app.Quit()
	}))

	g.win.SetContent(hboxConfig)
	g.win.Show()
}
