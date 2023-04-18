package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type modalWindow struct {
	dispatcher *mq.Dispatcher
	grid       *tview.Grid
	form       *tview.Form
}

func newModal(dispatcher *mq.Dispatcher) *modalWindow {
	m := &modalWindow{}
	m.dispatcher = dispatcher

	m.grid = tview.NewGrid()
	m.grid.SetRows(-1, -1, -1)
	m.grid.SetColumns(-1, -1, -1)

	m.form = tview.NewForm()
	m.form.SetButtonBackgroundColor(cyan)
	m.form.SetFieldTextColor(black)
	m.form.SetBackgroundColor(gray)
	m.form.SetBorder(true)
	m.form.SetTitle("Form")
	m.form.SetHorizontal(true)
	m.form.AddInputField("Book Title", "", 40, nil, nil)
	m.form.AddButton("Create Audiobook", m.Close)
	// m.form.SetDoneFunc(func(){})

	m.grid.AddItem(m.form, 1, 1, 1, 1, 0, 0, false)

	return m
}

func (mw *modalWindow) sendMessage(from string, to string, dtoType string, dto dto.Dto, async bool) {
	m := &mq.Message{}
	m.From = from
	m.To = to
	m.Type = dtoType
	m.Dto = dto
	m.Async = async
	mw.dispatcher.SendMessage(m)
}

func (m *modalWindow) Show() {
	m.sendMessage(mq.ModalWindow, mq.Frame, dto.AddPageCommandType, &dto.AddPageCommand{Name: "Modal", Grid: m.grid}, false)
	m.sendMessage(mq.ModalWindow, mq.Frame, dto.ShowPageCommandType, &dto.ShowPageCommand{Name: "Modal"}, false)
	m.sendMessage(mq.ModalWindow, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: m.form}, false)
	m.sendMessage(mq.ModalWindow, mq.TUI, dto.SetFocusCommandType, &dto.SetFocusCommand{Primitive: m.form}, false)
}

func (m *modalWindow) Close() {
	m.sendMessage(mq.ModalWindow, mq.Frame, dto.RemovePageCommandType, &dto.RemovePageCommand{Name: "Modal"}, false)
	m.sendMessage(mq.ModalWindow, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: nil}, true)
}
