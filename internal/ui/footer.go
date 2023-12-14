package ui

import (
	"abb_ia/internal/config"
	"abb_ia/internal/dto"
	"abb_ia/internal/mq"

	"github.com/vpoluyaktov/tview"
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
	f.busyIndicator.SetText("")
	f.busyIndicator.SetTextAlign(tview.AlignCenter)
	f.busyIndicator.SetBorder(false)
	f.busyIndicator.SetTextColor(busyIndicatorFgColor)
	f.busyIndicator.SetBackgroundColor(busyIndicatorBgColor)

	f.statusMessage = tview.NewTextView()
	f.statusMessage.SetText("")
	f.statusMessage.SetTextAlign(tview.AlignLeft)
	f.statusMessage.SetBorder(false)
	f.statusMessage.SetTextColor(footerFgColor)
	f.statusMessage.SetBackgroundColor(footerBgColor)

	f.version = tview.NewTextView()
	f.version.SetText("v" + config.Instance().AppVersion() + " (" + config.Instance().GetBuildDate() + ") ")
	f.version.SetTextAlign(tview.AlignRight)
	f.version.SetBorder(false)
	f.version.SetTextColor(footerFgColor)
	f.version.SetBackgroundColor(footerBgColor)

	f.grid = tview.NewGrid()
	f.grid.SetColumns(8, -1)
	f.grid.AddItem(f.busyIndicator, 0, 0, 1, 1, 0, 0, false)
	f.grid.AddItem(f.statusMessage, 0, 1, 1, 1, 0, 0, false)
	f.grid.AddItem(f.version, 0, 2, 1, 1, 0, 0, false)

	return f
}

func (f *footer) checkMQ() {
	m := f.mq.GetMessage(mq.Footer)
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
	f.statusMessage.SetText(s.Message)
	ui.Draw()
}

func (f *footer) toggleBusyIndicator(c *dto.SetBusyIndicator) {
	if c.Busy {
		f.busyFlag = true
		f.busyIndicator.SetTextColor(busyIndicatorFgColor)
		f.busyIndicator.SetBackgroundColor(busyIndicatorBgColor)
		f.busyIndicator.SetText(" Busy> ")
		go f.updateBusyIndicator()
	} else {
		f.busyFlag = false
		f.busyIndicator.SetText("")
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
