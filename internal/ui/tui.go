package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/event"
)

type TUI struct {
	// Message dispatcher
	dispatcher *event.Dispatcher
	// View components.
	app     *cview.Application
	colors  *colors
	header  *header
	footer  *footer
	search  *searchPanel
}

type Fn func()

func NewTUI(dispatcher *event.Dispatcher) *TUI {

	ui := TUI{}
	ui.app = cview.NewApplication()
	defer ui.app.HandlePanic()
	ui.colors = newColors()
	ui.app.EnableMouse(true)

	// Set Event Dispatcher
	ui.dispatcher = dispatcher

	// UI components
	ui.header = newHeader(ui.colors)
	ui.footer = newFooter(ui.colors)
	ui.search = newSearchPanel(dispatcher)

	// UI main frame
	f := newFrame()
	f.addHeader(ui.header)
	f.addFooter(ui.footer)
	f.addPannel("Search", ui.search.grid)

	ui.app.SetRoot(f.grid, true)
	return &ui
}

func (tui *TUI) Run() {
	if err := tui.app.Run(); err != nil {
		panic(err)
	}
}
