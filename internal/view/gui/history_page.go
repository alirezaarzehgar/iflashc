package gui

import (
	"context"
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

func (g GUI) historyPage(q *query.Queries, cfgs config.Config, pageSelector *widget.Select) *fyne.Container {
	searchQueryParams := query.ListStoredWordsParams{}
	wordListCreator := make(chan any)
	wordList := container.NewVBox()
	historyPage := &fyne.Container{}

	go func() {
		for {
			select {
			case <-wordListCreator:
				ctx := context.Background()
				list, err := q.ListStoredWords(ctx, searchQueryParams)
				if err != nil {
					dialog.ShowError(err, g.win)
				}

				wordList.RemoveAll()
				for _, l := range list {
					wordList.Add(widget.NewButton(l.Word, func() {
						title := canvas.NewText(l.Word, color.White)
						title.TextSize = DefaultTitleSize
						title.Alignment = fyne.TextAlignCenter

						rt := widget.NewRichTextFromMarkdown(l.Exp)
						rt.Wrapping = fyne.TextWrapBreak

						backBtn := widget.NewButton("back", func() {
							g.win.Resize(fyne.NewSize(DefaultDictWindowLen, DefaultDictWindowLen))
							g.win.SetContent(historyPage)
						})

						g.setContent(container.NewVBox(title, rt, backBtn))
						g.win.Resize(fyne.NewSize(DefaultDictWindowLen, 0))
					}))
				}
			}
		}
	}()

	ctx := context.Background()
	languages, err := q.ListStoredLanguages(ctx)
	if err != nil {
		g.ShowError("failed to get languages from database", err)
		return nil
	}
	languages = append([]string{""}, languages...)

	contexts, err := q.ListStoredHistoryContexts(ctx)
	if err != nil {
		g.ShowError("failed to get languages from database", err)
		return nil
	}
	contexts = append([]string{""}, contexts...)

	contextSelector := widget.NewSelect(contexts, func(s string) {
		searchQueryParams.Context = s
		wordListCreator <- struct{}{}
	})
	if len(contexts) >= 1 {
		contextSelector.SetSelected(cfgs[config.DefaultKeys.Context])
	}

	langSelector := widget.NewSelect(languages, func(s string) {
		searchQueryParams.Lang = s
		wordListCreator <- struct{}{}
	})
	if len(languages) >= 1 {
		langSelector.SetSelected(cfgs[config.DefaultKeys.DestLang])
	}

	translators := append([]string{""}, config.ConfigurableTranslators...)
	transSelector := widget.NewSelect(translators, func(s string) {
		searchQueryParams.Translator = s
		wordListCreator <- struct{}{}
	})
	transSelector.SetSelected(cfgs[config.DefaultKeys.Translator])

	searchEntry := widget.NewEntry()
	searchEntry.PlaceHolder = "Search word"
	searchBind := binding.NewString()
	searchEntry.Bind(searchBind)
	searchBind.AddListener(binding.NewDataListener(func() {
		searchQueryParams.WordLike, _ = searchBind.Get()
		wordListCreator <- struct{}{}
	}))

	historyPage = container.NewBorder(
		container.NewGridWithColumns(2,
			container.NewGridWithRows(2, pageSelector, searchEntry),
			container.NewHBox(
				container.NewVBox(widget.NewLabel("Languages:"), langSelector),
				container.NewVBox(widget.NewLabel("Translators:"), transSelector),
				container.NewVBox(widget.NewLabel("Contexts:"), contextSelector),
			),
		),
		nil, nil, nil,
		container.NewVScroll(wordList),
	)

	return historyPage
}
