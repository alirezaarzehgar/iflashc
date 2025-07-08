package gui

import (
	"context"
	"fmt"
	"image/color"
	"slices"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
	"github.com/alirezaarzehgar/iflashc/internal/view/ui"
)

func init() {
	ui.TypeRegistry[Backend] = New
}

const (
	MaxTitleLen = 45
	Backend     = "gui"
)

var (
	DefaultTitleSize     float32   = 35
	DefaultWindowSize    fyne.Size = fyne.NewSize(600, 200)
	DefaultDictWindowLen float32   = 800
)

type GUI struct {
	app  fyne.App
	win  fyne.Window
	lock *sync.Mutex
}

func New() ui.UI {
	ui.LoadingPageSize = ui.WindowSize{Width: 400, Height: 200}

	a := app.New()
	w := a.NewWindow("Integrated Flashcard")
	w.Resize(DefaultWindowSize)
	w.SetFixedSize(true)
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if slices.Contains([]fyne.KeyName{fyne.KeyEscape, fyne.KeyQ, fyne.KeyT, fyne.KeyG}, ke.Name) {
			a.Quit()
		}
	})
	return GUI{win: w, app: a, lock: &sync.Mutex{}}
}

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

func (g GUI) Run() {
	g.app.Run()
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

type keyEntry map[string]*widget.Entry

func (g GUI) ManageConfigs(q *query.Queries, cfgs config.Config) {
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

func (g GUI) setContent(page *fyne.Container) {
	g.lock.Lock()
	g.win.SetContent(page)
	g.lock.Unlock()
}

func (g GUI) Dashboard(q *query.Queries, cfgs config.Config) {
	g.win.Resize(fyne.NewSize(DefaultDictWindowLen, DefaultDictWindowLen))

	historyPage := &fyne.Container{}
	notesPage := &fyne.Container{}

	pageSelector := widget.NewSelect([]string{"history", "notes"}, func(s string) {
		switch s {
		case "history":
			g.setContent(historyPage)
		case "notes":
			g.setContent(notesPage)
		}
	})
	pageSelector.SetSelected("history")

	historyPage = g.historyPage(q, cfgs, pageSelector)
	notesPage = g.notesPage(q, cfgs, pageSelector)
	g.setContent(historyPage)
	g.win.Show()
}

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
