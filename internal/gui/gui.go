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

func ShowWord(title, text string) error {
	a := app.New()
	w := a.NewWindow("Integrated Flashcard")
	w.Resize(fyne.NewSize(500, 0))
	w.SetFixedSize(true)

	if len(title) > MaxTitleLen {
		title = title[:MaxTitleLen]
	}

	l := canvas.NewText(title, color.White)
	l.TextSize = 35
	l.Alignment = fyne.TextAlignCenter

	rt := widget.NewRichTextFromMarkdown(text)
	rt.Wrapping = fyne.TextWrapBreak
	btn := widget.NewButton("store", func() {
		a.Quit()
	})

	w.SetContent(container.NewVBox(l, rt, btn))
	w.ShowAndRun()
	return nil
}
