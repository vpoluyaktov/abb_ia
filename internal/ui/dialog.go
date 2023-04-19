package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type dialogWindow struct {
	mq     *mq.Dispatcher
	grid   *tview.Grid
	form   *tview.Form
	shadow *tview.Grid
	height int
	width  int
}

func newDialogWindow(dispatcher *mq.Dispatcher, height int, width int) *dialogWindow {
	d := &dialogWindow{}
	d.mq = dispatcher
	d.height = height
	d.width = width

	// shadow background
	d.shadow = tview.NewGrid()
	d.shadow.SetRows(2, -1, d.height, -1)
	d.shadow.SetColumns(4, -1, d.width, -1)
	d.shadow.AddItem(tview.NewBox().SetBackgroundColor(black), 2, 2, 1, 1, 0, 0, false)

	// transparent background
	d.grid = tview.NewGrid()
	d.grid.SetRows(-1, d.height, -1)
	d.grid.SetColumns(-1, d.width, -1)

	return d
}

func (d *dialogWindow) Show() {
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, dto.AddPageCommandType, &dto.AddPageCommand{Name: "Shadow", Grid: d.shadow}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, dto.AddPageCommandType, &dto.AddPageCommand{Name: "DialogWindow", Grid: d.grid}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, dto.ShowPageCommandType, &dto.ShowPageCommand{Name: "DialogWindow"}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.TUI, dto.SetFocusCommandType, &dto.SetFocusCommand{Primitive: d.form}, true)
}

func (d *dialogWindow) Close() {
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, dto.RemovePageCommandType, &dto.RemovePageCommand{Name: "DialogWindow"}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, dto.RemovePageCommandType, &dto.RemovePageCommand{Name: "Shadow"}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: nil}, true)
}

func (d *dialogWindow) setForm(f *tview.Form) {
	d.form = f
	d.setFormAttr()
	d.grid.AddItem(d.form, 1, 1, 1, 1, 0, 0, true)
}

func (d *dialogWindow) setFormAttr() {
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
	d.form.SetHorizontal(false)
	d.form.SetButtonsAlign(tview.AlignCenter)
}

func messageDialog(dispatcher *mq.Dispatcher, title string, message string) {
	d := newDialogWindow(dispatcher, 12, 80)
	f := tview.NewForm()
	f.SetTitle(title)
	tv := tview.NewTextView()
	tv.SetText(message)
	tv.SetTextAlign(tview.AlignCenter)
	f.AddFormItem(tv)
	f.AddButton("Ok", func() {
		d.Close()
	})
	d.setForm(f)
	d.Show()

}
