package ui

import (
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type components interface {
	// Check for messages in message queue
	checkMQ()
}

type TUI struct {
	// Message dispatcher
	dispatcher *mq.Dispatcher
	// UI components
	components []components
	app        *cview.Application
	header     *header
	footer     *footer
	search     *searchPanel
}

type Fn func()

func NewTUI(dispatcher *mq.Dispatcher) *TUI {

	ui := TUI{}
	ui.app = cview.NewApplication()
	defer ui.app.HandlePanic()
	ui.app.EnableMouse(true)
	setColorTheme()

	// Set Event Dispatcher
	ui.dispatcher = dispatcher

	// UI components
	ui.header = newHeader(dispatcher)
	ui.footer = newFooter(dispatcher)
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
		ui.pullMq()
		for _, c := range ui.components {
			c.checkMQ()
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

func (ui *TUI) pullMq() {
	m := ui.dispatcher.GetMessage("TUI")
	if m != nil {
		ui.dispatchMessage(m)
	}
}

func (ui *TUI) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.CommandType:
		if c, ok := m.Dto.(dto.Command); ok {
			if c.Command == "RedrawUI" {
				ui.app.Draw()
			}
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}
