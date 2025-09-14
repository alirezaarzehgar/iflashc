package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g GUI) GetCommentAndConfirmNote(text string) (string, bool) {
	g.win.Resize(fyne.NewSize(600, 250))

	rt := widget.NewMultiLineEntry()
	rt.SetText(text)
	rt.SetPlaceHolder("Write your text")

	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder("Write a comment for your text")

	confirm := false
	okBtn := widget.NewButton("Save Note", func() {
		confirm = true
		g.app.Quit()
	})
	cancelBtn := widget.NewButton("Ignore", func() {
		g.app.Quit()
	})

	g.win.SetContent(container.NewGridWithRows(
		2,
		rt,
		container.NewVBox(entry, container.NewHSplit(okBtn, cancelBtn)),
	))
	g.win.ShowAndRun()
	return entry.Text, confirm
}
