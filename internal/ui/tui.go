package ui

import (
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type components interface {
	readMessages()
}

type TUI struct {
	// Message dispatcher
	dispatcher *mq.Dispatcher
	// UI components
	components []components
	app        *cview.Application
	colors     *colors
	header     *header
	footer     *footer
	search     *searchPanel
}

type Fn func()

func NewTUI(dispatcher *mq.Dispatcher) *TUI {

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
	ui.components = append(ui.components, ui.search)

	// UI main frame
	f := newFrame()
	f.addHeader(ui.header)
	f.addFooter(ui.footer)
	f.addPannel("Search", ui.search.grid)

	ui.app.SetRoot(f.grid, true)
	return &ui
}

func (ui *TUI) startEventListener() {
	for {
		for _, c := range ui.components {
			c.readMessages()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (ui *TUI) Run() {
	go ui.startEventListener()
	if err := ui.app.Run(); err != nil {
		panic(err)
	}
}
