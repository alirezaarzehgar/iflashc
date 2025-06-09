package ui

import (
	"image/color"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
