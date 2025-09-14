package gui

import (
	"slices"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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

func (g GUI) Run() {
	g.app.Run()
}

func (g GUI) setContent(page *fyne.Container) {
	g.lock.Lock()
	g.win.SetContent(page)
	g.lock.Unlock()
}
