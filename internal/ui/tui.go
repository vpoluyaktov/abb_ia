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
	components []components
	app        *tview.Application
}

type Fn func()

func NewTUI(dispatcher *mq.Dispatcher) *TUI {

	ui := TUI{}
	ui.app = tview.NewApplication()
	ui.app.EnableMouse(true)
	setColorTheme()

	// Set Event Dispatcher
	ui.dispatcher = dispatcher

	// UI components
	header := newHeader(dispatcher)
	footer := newFooter(dispatcher)
	searchPage := newSearchPage(dispatcher)
	downloadPage := newDownloadPage(dispatcher)

	// UI main frame
	frame := newFrame(dispatcher)
	frame.addHeader(header)
	frame.addFooter(footer)
	frame.addPage("SearchPage", searchPage.grid)
	frame.addPage("DownloadPage", downloadPage.grid)

	ui.components = append(ui.components, frame)
	ui.components = append(ui.components, searchPage)
	ui.components = append(ui.components, downloadPage)

	frame.showPage("SearchPage")

	ui.app.SetRoot(frame.grid, true)
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
