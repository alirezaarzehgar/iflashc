package ui

import (
	"context"
	"image/color"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
)

const (
	MaxTitleLen = 45
)

var (
	DefaultTitleSize  float32   = 35
	DefaultWindowSize fyne.Size = fyne.NewSize(600, 200)
)

type TextBox struct {
	Title, Text string
}

type gui struct {
	app fyne.App
	win fyne.Window
}

func NewGUI() gui {
	a := app.New()
	w := a.NewWindow("Integrated Flashcard")
	w.Resize(DefaultWindowSize)
	w.SetFixedSize(true)
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if slices.Contains([]fyne.KeyName{fyne.KeyEscape, fyne.KeyQ, fyne.KeyT, fyne.KeyG}, ke.Name) {
			w.Close()
		}
	})
	return gui{win: w, app: a}
}

func (g gui) ShowText(tb TextBox) {
	if len(tb.Title) > MaxTitleLen {
		tb.Title = tb.Title[:MaxTitleLen]
	}

	l := canvas.NewText(tb.Title, color.White)
	l.TextSize = DefaultTitleSize
	l.Alignment = fyne.TextAlignCenter

	rt := widget.NewRichTextFromMarkdown(tb.Text)
	rt.Wrapping = fyne.TextWrapBreak

	g.win.SetContent(container.NewVBox(l, rt))
	g.win.Show()
}

func (g gui) Run() {
	g.app.Run()
}

func (g gui) ShowError(text string, err error) {
	l := canvas.NewText(text, color.White)
	l.TextSize = DefaultTitleSize
	l.Alignment = fyne.TextAlignCenter

	rt := widget.NewRichTextFromMarkdown(err.Error())
	rt.Wrapping = fyne.TextWrapBreak

	g.win.SetContent(container.NewPadded(container.NewVBox(l, rt)))
	g.win.Show()
}

type keyEntry map[string]*widget.Entry

func (g gui) ManageConfigs(q *query.Queries, cfgs config.Config) {
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

func (g gui) Dashboard(q *query.Queries, cfgs config.Config) {

	ctx := context.Background()
	languages, err := q.ListStoredLanguages(ctx)
	if err != nil {
		g.ShowError("failed to get languages from database", err)
		return
	}

	langSelector := widget.NewSelect(languages, func(s string) {
		dialog.ShowInformation(s, s, g.win)
	})

	mainPage := container.NewVBox(langSelector)
	// wordPage := container.NewVBox()

	g.win.SetContent(mainPage)
	g.win.Show()
}
