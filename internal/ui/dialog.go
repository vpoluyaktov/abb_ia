package ui

import (
	"abb_ia/internal/dto"
	"abb_ia/internal/mq"

	"github.com/gdamore/tcell/v2"
	"github.com/vpoluyaktov/tview"
)

type dialogWindow struct {
	mq         *mq.Dispatcher
	grid       *tview.Grid
	form       *tview.Form
	shadow     *tview.Grid
	background *tview.Grid
	height     int
	width      int
	focus      tview.Primitive // set focus after the form close
}

func newDialogWindow(dispatcher *mq.Dispatcher, height int, width int, focus tview.Primitive) *dialogWindow {
	d := &dialogWindow{}
	d.mq = dispatcher
	d.height = height
	d.width = width
	d.focus = focus

	// shadow background
	d.shadow = tview.NewGrid()
	d.shadow.SetRows(2, -1, d.height, -1)
	d.shadow.SetColumns(4, -1, d.width, -1)
	d.shadow.AddItem(tview.NewBox().SetBackgroundColor(black), 2, 2, 1, 1, 0, 0, false)

	// gray background
	d.background = tview.NewGrid()
	d.background.SetRows(-1, d.height, -1)
	d.background.SetColumns(-1, d.width, -1)
	d.background.AddItem(tview.NewBox().SetBackgroundColor(gray), 1, 1, 1, 1, 0, 0, false)

	// transparent background
	d.grid = tview.NewGrid()
	d.grid.SetRows(-1, d.height, -1)
	d.grid.SetColumns(-1, d.width, -1)

	d.grid.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		return action, event
	})

	d.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			{
				d.Close()
			}
		}
		return event
	})

	return d
}

func (d *dialogWindow) Show() {
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, &dto.AddPageCommand{Name: "Shadow", Grid: d.shadow}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, &dto.AddPageCommand{Name: "Background", Grid: d.background}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, &dto.AddPageCommand{Name: "DialogWindow", Grid: d.grid}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, &dto.ShowPageCommand{Name: "DialogWindow"}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.TUI, &dto.SetFocusCommand{Primitive: d.form}, false)
	ui.Draw()
}

func (d *dialogWindow) Close() {
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, &dto.RemovePageCommand{Name: "DialogWindow"}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, &dto.RemovePageCommand{Name: "Background"}, false)
	d.mq.SendMessage(mq.DialogWindow, mq.Frame, &dto.RemovePageCommand{Name: "Shadow"}, false)
	if d.focus != nil {
		ui.SetFocus(d.focus)
	}
	ui.Draw()
}

func (d *dialogWindow) setForm(f *tview.Form) {
	d.form = f
	d.setFormAttributes()
	d.grid.AddItem(d.form, 1, 1, 1, 1, 0, 0, true)
}

func (d *dialogWindow) setFormAttributes() {
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

type OkFunc func()

func newMessageDialog(dispatcher *mq.Dispatcher, title string, message string, focus tview.Primitive, okFunc OkFunc) {
	d := newDialogWindow(dispatcher, 12, 80, focus)
	f := newForm()
	f.SetTitle(title)
	tv := tview.NewTextView()
	tv.SetWrap(true)
	tv.SetWordWrap(true)
	tv.SetDynamicColors(true)
	tv.SetText(message)
	tv.SetTextAlign(tview.AlignCenter)
	f.AddFormItem(tv)
	f.AddButton("Ok", func() {
		okFunc()
		d.Close()
	})
	d.setForm(f.Form)
	d.Show()
}

type YesNoFunc func()

func newYesNoDialog(dispatcher *mq.Dispatcher, title string, message string, focus tview.Primitive, yesFunc YesNoFunc, noFunc YesNoFunc) {
	d := newDialogWindow(dispatcher, 11, 60, focus)
	f := newForm()
	f.SetTitle(title)
	tv := tview.NewTextView()
	tv.SetWrap(true)
	tv.SetWordWrap(true)
	tv.SetDynamicColors(true)
	tv.SetText(message)
	tv.SetTextAlign(tview.AlignCenter)
	f.AddFormItem(tv)
	f.AddButton("Yes", func() {
		yesFunc()
		d.Close()
	})
	f.AddButton("No", func() {
		noFunc()
		d.Close()
	})
	d.setForm(f.Form)
	d.Show()
}
