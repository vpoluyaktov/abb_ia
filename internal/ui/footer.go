package ui

import (
	"time"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
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
	f.busyIndicator.SetTextColor(footerFgColor)
	f.busyIndicator.SetBackgroundColor(footerBgColor)

	f.statusMessage = tview.NewTextView()
	f.statusMessage.SetText("")
	f.statusMessage.SetTextAlign(tview.AlignLeft)
	f.statusMessage.SetBorder(false)
	f.statusMessage.SetTextColor(footerFgColor)
	f.statusMessage.SetBackgroundColor(footerBgColor)

	f.version = tview.NewTextView()
	f.version.SetText("v.0.0.1")
	f.version.SetTextAlign(tview.AlignRight)
	f.version.SetBorder(false)
	f.version.SetTextColor(footerFgColor)
	f.version.SetBackgroundColor(footerBgColor)

	f.grid = tview.NewGrid()
	f.grid.SetColumns(10, -1, 10)
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
	switch t := m.Type; t {
	case dto.UpdateStatusType:
		if s, ok := m.Dto.(*dto.UpdateStatus); ok {
			f.updateStatus(s)
		} else {
			m.DtoCastError(mq.Footer)
		}
	case dto.SetBusyIndicatorType:
		if c, ok := m.Dto.(*dto.SetBusyIndicator); ok {
			f.toggleBusyIndicator(c)
		} else {
			m.DtoCastError(mq.Footer)
		}

	default:
		m.UnsupportedTypeError(mq.Footer)
	}
}

func (f *footer) updateStatus(s *dto.UpdateStatus) {
	f.statusMessage.SetText(s.Message)
	f.mq.SendMessage(mq.Footer, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: nil}, true)
}

func (f *footer) toggleBusyIndicator(c *dto.SetBusyIndicator) {
	if c.Busy {
		f.busyFlag = true
		go f.updateBusyIndicator()
	} else {
		f.busyFlag = false

	}
}

func (f *footer) updateBusyIndicator() {
	// busyChars := []string{"[>    ]", "[ >   ]", "[  >  ]", "[   > ]", "[    >]", "[    <]", "[   < ]", "[  <  ]", "[ <   ]", "[<    ]"}
	busyChars := []string{"[O----]", "[-O---]", "[--O--]", "[---O-]", "[----O]", "[---O-]", "[--O--]", "[-O---]"}
	for f.busyFlag {
		for i := 0; i < len(busyChars); i++ {
			f.busyIndicator.SetText(busyChars[i])
			f.mq.SendMessage(mq.Footer, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: nil}, true)
			time.Sleep(250 * time.Millisecond)
			if !f.busyFlag {
				break
			}
		}
	}
	f.busyIndicator.SetText("")
	f.mq.SendMessage(mq.Footer, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: nil}, true)
}
