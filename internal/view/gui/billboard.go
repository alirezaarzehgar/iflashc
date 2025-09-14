package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/alirezaarzehgar/iflashc/internal/view/ui"
)

func (g GUI) ShowText(tb ui.TextBox) {
	if tb.Size.Height+tb.Size.Width > 0 {
		g.win.Resize(fyne.NewSize(tb.Size.Width, tb.Size.Width))
	}

	if len(tb.Title) > MaxTitleLen {
		tb.Title = tb.Title[:MaxTitleLen]
	}

	l := canvas.NewText(tb.Title, color.White)
	l.TextSize = DefaultTitleSize
	l.Alignment = fyne.TextAlignCenter

	rt := widget.NewRichTextFromMarkdown(tb.Text)

	g.win.SetContent(container.NewVBox(l, rt))
	g.win.Show()
	g.win.CenterOnScreen()
}

func (g GUI) ShowError(text string, err error) {
	g.win.Resize(DefaultWindowSize)

	l := canvas.NewText(text, color.White)
	l.TextSize = DefaultTitleSize
	l.Alignment = fyne.TextAlignCenter

	rt := widget.NewRichTextFromMarkdown(err.Error())
	rt.Wrapping = fyne.TextWrapBreak

	g.win.SetContent(container.NewPadded(container.NewVBox(l, rt)))
	g.win.Show()
	g.win.CenterOnScreen()
}
