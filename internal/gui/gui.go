package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowWord(title, text string) error {
	a := app.New()
	w := a.NewWindow("Integrated Flashcard")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(300, 300))

	l := canvas.NewText(title, color.White)
	l.TextSize = 35
	l.Alignment = fyne.TextAlignCenter

	rt := widget.NewRichTextFromMarkdown(text)
	btn := widget.NewButton("store", func() {
		a.Quit()
	})

	w.CenterOnScreen()
	w.SetContent(container.NewVBox(l, rt, btn))
	w.ShowAndRun()
	return nil
}
