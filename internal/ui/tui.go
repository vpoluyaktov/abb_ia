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
	// Message mq
	mq         *mq.Dispatcher
	components []components
	app        *tview.Application
}

type Fn func()

func NewTUI(dispatcher *mq.Dispatcher) *TUI {

	ui := TUI{}
	ui.app = tview.NewApplication()
	ui.app.EnableMouse(true)
	setColorTheme()

	// Set Event Dispatcher and recepient
	ui.mq = dispatcher

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
	ui.components = append(ui.components, header)
	ui.components = append(ui.components, footer)
	ui.components = append(ui.components, searchPage)
	ui.components = append(ui.components, downloadPage)

	frame.showPage("SearchPage")

	ui.app.SetRoot(frame.grid, true)
	return &ui
}

func (ui *TUI) startEventListener() {
	for {
		ui.checkMQ()
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

func (ui *TUI) checkMQ() {
	m := ui.mq.GetMessage(mq.TUI)
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
				// ui.app.Draw(c.Primitive) // not supported by rivo/tview
			}
		} else {
			m.DtoCastError(mq.TUI)
		}
	case dto.SetFocusCommandType:
		if c, ok := m.Dto.(*dto.SetFocusCommand); ok {
			ui.app.SetFocus(c.Primitive)
		} else {
			m.DtoCastError(mq.TUI)
		}

	default:
		m.UnsupportedTypeError(mq.TUI)
	}
}
