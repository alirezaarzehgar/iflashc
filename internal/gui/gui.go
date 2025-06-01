package gui

import (
	"image/color"

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
	DefaultTitleSize   float32   = 35
	DefaultButtonLabel string    = "Store"
	DefaultWindowSize  fyne.Size = fyne.NewSize(600, 0)
)

type TextBox struct {
	Title, Text string
	HaveBtn     bool
}

func ShowText(tb TextBox) error {
	a := app.New()
	w := a.NewWindow("Integrated Flashcard")
	w.Resize(DefaultWindowSize)
	w.SetFixedSize(true)

	if len(tb.Title) > MaxTitleLen {
		tb.Title = tb.Title[:MaxTitleLen]
	}

	l := canvas.NewText(tb.Title, color.White)
	l.TextSize = DefaultTitleSize
	l.Alignment = fyne.TextAlignCenter

	rt := widget.NewRichTextFromMarkdown(tb.Text)
	rt.Wrapping = fyne.TextWrapBreak

	btn := widget.NewButton(DefaultButtonLabel, func() {
		a.Quit()
	})
	qbtn := widget.NewButton("Quit", func() {
		a.Quit()
	})

	vbox := container.NewVBox(l, rt, qbtn)

	if tb.HaveBtn {
		vbox.Add(btn)
	}

	w.SetContent(vbox)
	w.ShowAndRun()
	return nil
}
