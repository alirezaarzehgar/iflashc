package ui

import (
	"context"
	"image/color"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
)

const (
	MaxTitleLen = 45
)

var (
	DefaultTitleSize  float32   = 35
	DefaultWindowSize fyne.Size = fyne.NewSize(600, 0)
)

type TextBox struct {
	Title, Text string
}

type gui struct {
	app fyne.App
}

func NewGUI() gui {
	a := app.New()
	return gui{app: a}
}

func (g gui) ShowText(tb TextBox) {
	w := g.app.NewWindow("Integrated Flashcard")
	w.Resize(DefaultWindowSize)
	w.SetFixedSize(true)
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if slices.Contains([]fyne.KeyName{fyne.KeyEscape, fyne.KeyQ, fyne.KeyT, fyne.KeyG}, ke.Name) {
			w.Close()
		}
	})

	if len(tb.Title) > MaxTitleLen {
		tb.Title = tb.Title[:MaxTitleLen]
	}

	l := canvas.NewText(tb.Title, color.White)
	l.TextSize = DefaultTitleSize
	l.Alignment = fyne.TextAlignCenter

	rt := widget.NewRichTextFromMarkdown(tb.Text)
	rt.Wrapping = fyne.TextWrapBreak

	w.SetContent(container.NewVBox(l, rt))
	w.Show()
}

func (g gui) Run() {
	g.app.Run()
}

func (g gui) Dashboard(q *query.Queries, cfgs config.Config) {
	w := g.app.NewWindow("Dashboard")
	w.Resize(fyne.NewSize(600, 600))
	w.Show()
}

var (
	configEntries = []string{
		config.DefaultKeys.Translator,
		config.DefaultKeys.GroqApiKey,
		config.DefaultKeys.GroqModel,
		config.DefaultKeys.DestLang,
		config.DefaultKeys.Socks5,
	}
)

func (g gui) ManageConfigs(q *query.Queries, cfgs config.Config) {
	w := g.app.NewWindow("Configration Manager")
	w.Resize(DefaultWindowSize)
	w.SetFixedSize(true)

	status := widget.NewLabel("nothing changed")
	hboxConfig := container.NewVBox()

	l := canvas.NewText("Configuration Manager", color.White)
	l.TextSize = DefaultTitleSize
	l.Alignment = fyne.TextAlignCenter

	hboxConfig.Add(l)
	hboxConfig.Add(container.NewCenter(status))

	for _, k := range []string{
		config.DefaultKeys.Translator,
		config.DefaultKeys.DestLang,
		config.DefaultKeys.GroqApiKey,
		config.DefaultKeys.GroqModel,
		config.DefaultKeys.Socks5,
	} {
		e := widget.NewEntry()
		e.Text = cfgs[k]
		b := widget.NewButton(k, func() {
			ctx := context.Background()
			err := q.ChangeConfig(ctx, query.ChangeConfigParams{Key: k, Value: e.Text})
			if err != nil {
				status.SetText("failed to save config: " + err.Error())
			} else {
				status.SetText("config changed successfuly")
			}
		})

		hbox := container.NewGridWithColumns(2, b, e)
		hboxConfig.Add(hbox)
	}

	hboxConfig.Add(widget.NewButton("Quit", func() {
		g.app.Quit()
	}))

	w.SetContent(hboxConfig)
	w.Show()
}
