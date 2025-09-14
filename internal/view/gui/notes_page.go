package gui

import (
	"context"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
)

func (g GUI) notesPage(q *query.Queries, cfgs config.Config, pageSelector *widget.Select) *fyne.Container {
	searchQueryParams := query.ListStoredNotesParams{}
	noteListCreator := make(chan any)
	noteList := container.NewVBox()
	notesPage := &fyne.Container{}

	go func() {
		for {
			select {
			case <-noteListCreator:
				ctx := context.Background()
				list, err := q.ListStoredNotes(ctx, searchQueryParams)
				if err != nil {
					dialog.ShowError(err, g.win)
				}

				noteList.RemoveAll()
				for _, l := range list {
					preview := l.Note
					if len(preview) >= 20 {
						preview = preview[:20]
					}
					preview = fmt.Sprintf("%s | %d", preview, l.Occurrence)

					noteList.Add(widget.NewButton(preview, func() {
						title := canvas.NewText("Your Note", color.White)
						title.TextSize = DefaultTitleSize
						title.Alignment = fyne.TextAlignCenter

						note := widget.NewRichTextFromMarkdown(l.Note)
						note.Wrapping = fyne.TextWrapBreak

						comment := widget.NewRichTextFromMarkdown(l.Comment)
						comment.Wrapping = fyne.TextWrapBreak

						backBtn := widget.NewButton("back", func() {
							g.win.Resize(fyne.NewSize(DefaultDictWindowLen, DefaultDictWindowLen))
							g.setContent(notesPage)
						})

						g.setContent(container.NewVBox(title, note, comment, backBtn))
						g.win.Resize(fyne.NewSize(DefaultDictWindowLen, 0))
					}))
				}
			}
		}
	}()

	ctx := context.Background()
	contexts, err := q.ListStoredNoteContexts(ctx)
	if err != nil {
		g.ShowError("failed to get languages from database", err)
		return nil
	}
	contexts = append([]string{""}, contexts...)

	contextSelector := widget.NewSelect(contexts, func(s string) {
		searchQueryParams.Context = s
		noteListCreator <- struct{}{}
	})
	if len(contexts) >= 1 {
		contextSelector.SetSelected(cfgs[config.DefaultKeys.Context])
	}

	searchEntry := widget.NewEntry()
	searchEntry.PlaceHolder = "Search word"
	searchBind := binding.NewString()
	searchEntry.Bind(searchBind)
	searchBind.AddListener(binding.NewDataListener(func() {
		searchQueryParams.Column1, _ = searchBind.Get()
		noteListCreator <- struct{}{}
	}))

	notesPage = container.NewBorder(
		container.NewGridWithColumns(2,
			container.NewGridWithRows(2, pageSelector, searchEntry),
			container.NewHBox(
				container.NewVBox(widget.NewLabel("Contexts:"), contextSelector),
			),
		),
		nil, nil, nil,
		container.NewVScroll(noteList),
	)

	return notesPage
}
