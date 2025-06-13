package ui

import (
	"context"
	"image/color"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
			a.Quit()
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
	g.win.Resize(fyne.NewSize(800, 800))

	searchQueryParams := query.ListStoredWordsParams{}
	wordListCreator := make(chan any)
	wordList := container.NewVBox()

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
					wordList.Add(widget.NewLabel(l.Word))
				}
			}
		}
	}()

	ctx := context.Background()
	languages, err := q.ListStoredLanguages(ctx)
	if err != nil {
		g.ShowError("failed to get languages from database", err)
		return
	}

	contexts, err := q.ListStoredContexts(ctx)
	if err != nil {
		g.ShowError("failed to get languages from database", err)
		return
	}

	contextSelector := widget.NewSelect(contexts, func(s string) {
		searchQueryParams.Context = s
		wordListCreator <- 0
	})
	if len(contexts) >= 1 {
		contextSelector.SetSelected(contexts[0])
	}

	langSelector := widget.NewSelect(languages, func(s string) {
		searchQueryParams.Lang = s
		wordListCreator <- 0
	})
	if len(languages) >= 1 {
		langSelector.SetSelected(languages[0])
	}

	transSelector := widget.NewSelect(config.ConfigurableTranslators, func(s string) {
		searchQueryParams.Translator = s
		wordListCreator <- 0
	})
	transSelector.SetSelected(config.ConfigurableTranslators[0])

	searchEntry := widget.NewEntry()
	searchEntry.PlaceHolder = "Search word"
	searchBind := binding.NewString()
	searchEntry.Bind(searchBind)
	searchBind.AddListener(binding.NewDataListener(func() {
		searchQueryParams.WordLike, _ = searchBind.Get()
		wordListCreator <- 0
	}))

	mainPage := container.NewBorder(
		container.NewGridWithColumns(2,
			container.NewGridWithRows(2, widget.NewLabel(""), searchEntry),
			container.NewHBox(
				container.NewVBox(widget.NewLabel("Languages:"), langSelector),
				container.NewVBox(widget.NewLabel("Translators:"), transSelector),
				container.NewVBox(widget.NewLabel("Contexts:"), contextSelector),
			),
		),
		nil, nil, nil,
		container.NewVScroll(wordList),
	)

	g.win.SetContent(mainPage)
	g.win.Show()
}
