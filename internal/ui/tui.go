package ui

import (
	"os"
	"os/exec"
	"time"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
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
	encodingPage := newEncodingPage(dispatcher)
	chaptersPage := newChaptersPage(dispatcher)

	// UI main frame
	frame := newFrame(dispatcher)
	frame.addHeader(header)
	frame.addFooter(footer)
	frame.addPage("SearchPage", searchPage.grid)
	frame.addPage("DownloadPage", downloadPage.grid)
	frame.addPage("EncodingPage", encodingPage.grid)
	frame.addPage("ChaptersPage", chaptersPage.grid)

	ui.components = append(ui.components, frame)
	ui.components = append(ui.components, header)
	ui.components = append(ui.components, footer)
	ui.components = append(ui.components, searchPage)
	ui.components = append(ui.components, downloadPage)
	ui.components = append(ui.components, encodingPage)
	ui.components = append(ui.components, chaptersPage)

	frame.switchToPage("SearchPage")

	ui.app.SetRoot(frame.grid, true)
	return &ui
}

func (ui *TUI) startEventListener() {
	for {
		ui.checkMQ()
		for _, c := range ui.components {
			c.checkMQ()
		}
		time.Sleep(mq.PullFrequency)
	}
}

func (ui *TUI) Run() {
	defer func() {
		c := exec.Command("reset")
		c.Stdout = os.Stdout
		c.Run()
	}()
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
	switch cmd := m.Dto.(type) {
	case *dto.DrawCommand:
		if cmd.Primitive == nil {
			ui.app.Draw()
		} else {
			// ui.app.Draw(cmd.Primitive) // not supported by rivo/tview
		}
	case *dto.SetFocusCommand:
		ui.app.SetFocus(cmd.Primitive)
	default:
		m.UnsupportedTypeError(mq.TUI)
	}
}
