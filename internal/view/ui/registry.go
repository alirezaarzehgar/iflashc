package ui

import (
	"github.com/alirezaarzehgar/iflashc/internal/config"
	"github.com/alirezaarzehgar/iflashc/internal/query"
)

var (
	LoadingPageSize = WindowSize{}
)

type WindowSize struct {
	Width  float32
	Height float32
}

type TextBox struct {
	Title, Text string
	Size        WindowSize
}

type UI interface {
	ShowText(tb TextBox)
	GetCommentAndConfirmNote(text string) (string, bool)
	ShowError(text string, err error)
	ManageConfigs(q *query.Queries, cfgs config.Config)
	Dashboard(q *query.Queries, cfgs config.Config)
	Run()
}

var TypeRegistry = make(map[string]func() UI)

// TODO: This over-engineering considered for future usages.
// I believe it is simple enough. Add more ui backend if needed
func New(uiBackends string) UI {
	return TypeRegistry[uiBackends]()
}
