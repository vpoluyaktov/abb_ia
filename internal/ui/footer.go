package ui

import (
	"fmt"

	"abb_ia/internal/config"
	"abb_ia/internal/dto"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"

	"github.com/vpoluyaktov/tview"
)

const (
	idleMessage          = " IDLE "
	busyMessage          = " BUSY "
	defaultStatusMessage = " [black]Keys and shortcuts: [yellow]Arrows, Tab, Shift+Tab[darkblue]: navigation,  [yellow]Enter[darkblue]: action,  [yellow]Ctrl+U[darkblue]: clear a field,  [yellow]Ctrl+C[darkblue]: exit"
)

type footer struct {
	mq            *mq.Dispatcher
	grid          *tview.Grid
	busyFlag      bool
	busyIndicator *tview.TextView
	statusMessage *tview.TextView
	version       *tview.TextView
}

func newFooter(dispatcher *mq.Dispatcher) *footer {
	f := &footer{}
	f.mq = dispatcher
	f.mq.RegisterListener(mq.Footer, f.dispatchMessage)

	f.busyIndicator = tview.NewTextView()
	f.busyIndicator.SetText(idleMessage)
	f.busyIndicator.SetTextAlign(tview.AlignCenter)
	f.busyIndicator.SetBorder(false)
	f.busyIndicator.SetTextColor(footerFgColor)
	f.busyIndicator.SetBackgroundColor(busyIndicatorBgColor)
	f.busyIndicator.SetDynamicColors(true)

	f.statusMessage = tview.NewTextView()
	f.statusMessage.SetText(defaultStatusMessage)
	f.statusMessage.SetTextAlign(tview.AlignLeft)
	f.statusMessage.SetBorder(false)
	f.statusMessage.SetTextColor(footerFgColor)
	f.statusMessage.SetBackgroundColor(footerBgColor)
	f.statusMessage.SetDynamicColors(true)

	f.version = tview.NewTextView()
	f.version.SetText("v" + config.Instance().AppVersion() + " (" + config.Instance().GetBuildDate() + ") ")
	f.version.SetTextAlign(tview.AlignRight)
	f.version.SetBorder(false)
	f.version.SetTextColor(footerFgColor)
	f.version.SetBackgroundColor(footerBgColor)

	f.grid = tview.NewGrid()
	f.grid.SetColumns(8, -1, 25)
	f.grid.AddItem(f.busyIndicator, 0, 0, 1, 1, 0, 0, false)
	f.grid.AddItem(f.statusMessage, 0, 1, 1, 1, 0, 0, false)
	f.grid.AddItem(f.version, 0, 2, 1, 1, 0, 0, false)

	return f
}

func (f *footer) checkMQ() {
	m, err := f.mq.GetMessage(mq.Footer)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get message for Footer: %v", err))
		return
	}
	if m != nil {
		f.dispatchMessage(m)
	}
}

func (f *footer) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.UpdateStatus:
		f.updateStatus(dto)
	case *dto.SetBusyIndicator:
		f.toggleBusyIndicator(dto)
	default:
		m.UnsupportedTypeError(mq.Footer)
	}
}

func (f *footer) updateStatus(s *dto.UpdateStatus) {
	if s.Message != "" {
		f.statusMessage.SetTextColor(black)
		f.statusMessage.SetText(" " + s.Message)
	} else {
		f.statusMessage.SetTextColor(footerFgColor)
		f.statusMessage.SetText(defaultStatusMessage)
	}
	ui.Draw()
}

func (f *footer) toggleBusyIndicator(c *dto.SetBusyIndicator) {
	if c.Busy {
		f.busyFlag = true
		f.busyIndicator.SetTextColor(yellow)
		f.busyIndicator.SetBackgroundColor(darkRed)
		f.busyIndicator.SetText(busyMessage)
		go f.updateBusyIndicator()
	} else {
		f.busyFlag = false
		f.busyIndicator.SetText(idleMessage)
		f.busyIndicator.SetTextColor(footerFgColor)
		f.busyIndicator.SetBackgroundColor(footerBgColor)
		ui.Draw()
	}
}

func (f *footer) updateBusyIndicator() {
	// busyChars := []string{"[>    ]", "[ >   ]", "[  >  ]", "[   > ]", "[    >]", "[    <]", "[   < ]", "[  <  ]", "[ <   ]", "[<    ]"}
	// busyChars := []string{"[O-----]", "[-O----]", "[--O---]", "[---O--]", "[----O-]", "[-----O]", "[----O-]", "[---O--]", "[--O---]", "[-O----]"}
	// busyChars := []string{"█▒▒▒▒▒", "▒█▒▒▒▒", "▒▒█▒▒▒", "▒▒▒█▒▒", "▒▒▒▒█▒", "▒▒▒▒▒█", "▒▒▒▒█▒", "▒▒▒█▒▒", "▒▒█▒▒▒", "▒█▒▒▒▒"}
	// busyChars := []string{"█     ", " █    ", "  █   ", "   █  ", "    █ ", "     █", "    █ ", "   █  ", "  █   ", " █    "}
	// for f.busyFlag {
	// 	for i := 0; i < len(busyChars); i++ {
	// 		f.busyIndicator.SetText(busyChars[i])
	// 		ui.Draw()
	// 		time.Sleep(200 * time.Millisecond)
	// 		if !f.busyFlag {
	// 			break
	// 		}
	// 	}
	// }
}
