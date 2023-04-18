package ui

import (
	"time"

	"github.com/rivo/tview"
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
	app        *tview.Application
	frame      *frame
	header     *header
	footer     *footer
	search     *searchPanel
}

type Fn func()

func NewTUI(dispatcher *mq.Dispatcher) *TUI {

	ui := TUI{}
	ui.app = tview.NewApplication()
	// defer ui.app.pani
	ui.app.EnableMouse(true)
	setColorTheme()

	// Set Event Dispatcher
	ui.dispatcher = dispatcher

	// UI components
	ui.header = newHeader(dispatcher)
	ui.footer = newFooter(dispatcher)
	ui.search = newSearchPanel(dispatcher)

	// UI main frame
	ui.frame = newFrame(dispatcher)
	ui.frame.addHeader(ui.header)
	ui.frame.addFooter(ui.footer)
	ui.frame.addPannel("Search", ui.search.grid)
	ui.components = append(ui.components, ui.frame)
	ui.components = append(ui.components, ui.search)

	ui.app.SetRoot(ui.frame.grid, true)
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
	case dto.DrawCommandType:
		if c, ok := m.Dto.(*dto.DrawCommand); ok {
			if c.Primitive == nil {
				ui.app.Draw()
			} else {
				// ui.app.Draw(c.Primitive)
			}
		} else {
			m.DtoCastError()
		}
	case dto.SetFocusCommandType:
		if c, ok := m.Dto.(*dto.SetFocusCommand); ok {
			ui.app.SetFocus(c.Primitive)
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}
