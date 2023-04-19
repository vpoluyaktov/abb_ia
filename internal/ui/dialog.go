package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type dialogWindow struct {
	dispatcher *mq.Dispatcher
	grid       *tview.Grid
	form       *tview.Form
	shadow     *tview.Grid
}

func newDialogWindow(dispatcher *mq.Dispatcher, height int, width int) *dialogWindow {
	d := &dialogWindow{}
	d.dispatcher = dispatcher

	d.grid = tview.NewGrid()
	d.grid.SetRows(-1, height, -1)
	d.grid.SetColumns(-1, width, -1)

	d.form = tview.NewForm()
	d.form.SetBorderColor(black)
	d.form.SetTitleColor(blue)
	d.form.SetLabelColor(blue)
	d.form.SetFieldTextColor(black)
	d.form.SetFieldBackgroundColor(cyan)
	d.form.SetButtonTextColor(white)
	d.form.SetButtonBackgroundColor(blue)
	d.form.SetFieldTextColor(black)
	d.form.SetBackgroundColor(gray)
	d.form.SetBorder(true)
	d.form.SetTitle(" Form ")
	d.form.SetHorizontal(false)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.AddInputField("Book Title", "Test", 40, nil, nil)
	d.form.SetButtonsAlign(tview.AlignCenter)
	d.form.AddButton("Create Audiobook", d.Close)
	d.form.AddButton("Cancel", d.Close)
	d.grid.AddItem(d.form, 1, 1, 1, 1, 0, 0, true)

	d.shadow = tview.NewGrid()
	d.shadow.SetRows(2, -1, height, -1)
	d.shadow.SetColumns(4, -1, width, -1)
	d.shadow.AddItem(tview.NewBox().SetBackgroundColor(black), 2, 2, 1, 1, 0, 0, false)

	return d
}

func (mw *dialogWindow) sendMessage(from string, to string, dtoType string, dto dto.Dto, async bool) {
	m := &mq.Message{}
	m.From = from
	m.To = to
	m.Type = dtoType
	m.Dto = dto
	m.Async = async
	mw.dispatcher.SendMessage(m)
}

func (m *dialogWindow) Show() {
	m.sendMessage(mq.DialogWindow, mq.Frame, dto.AddPageCommandType, &dto.AddPageCommand{Name: "Shadow", Grid: m.shadow}, false)
	m.sendMessage(mq.DialogWindow, mq.Frame, dto.AddPageCommandType, &dto.AddPageCommand{Name: "DialogWindow", Grid: m.grid}, false)
	m.sendMessage(mq.DialogWindow, mq.Frame, dto.ShowPageCommandType, &dto.ShowPageCommand{Name: "DialogWindow"}, false)
	m.sendMessage(mq.DialogWindow, mq.TUI, dto.SetFocusCommandType, &dto.SetFocusCommand{Primitive: m.grid}, true)
}

func (m *dialogWindow) Close() {
	m.sendMessage(mq.DialogWindow, mq.Frame, dto.RemovePageCommandType, &dto.RemovePageCommand{Name: "DialogWindow"}, false)
	m.sendMessage(mq.DialogWindow, mq.Frame, dto.RemovePageCommandType, &dto.RemovePageCommand{Name: "Shadow"}, false)
	m.sendMessage(mq.DialogWindow, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: nil}, true)
}
