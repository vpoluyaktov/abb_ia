package ui

import (
	"sync"

	"github.com/rivo/tview"
	// "github.com/vpoluyaktov/audiobook_creator_IA/internal/utils"
)

// //////////////////////////////////////////////////////////////
// tview.Table wrapper
// //////////////////////////////////////////////////////////////
type table struct {
	t         *tview.Table
	headers   []string
	widths    []int
	allWidths int
	aligns    []uint
}

func newTable() *table {
	t := &table{}
	t.t = tview.NewTable()
	t.t.SetSelectable(true, false)
	t.t.SetSeparator(tview.Borders.Vertical)
	// t.t.SetSortClicked(false)
	// t.t.SetSortFunc() // TODO
	// t.t.ShowFocus(true)
	t.t.SetBorder(false)
	return t
}

func (t *table) setHeaders(headers ...string) {
	t.headers = headers
}

func (t *table) setWidths(widths ...int) {
	t.widths = widths
	for _, w := range t.widths {
		t.allWidths += w
	}
}

func (t *table) setAlign(aligns ...uint) {
	t.aligns = aligns
}

func (t *table) showHeader() {
	for c, h := range t.headers {
		cell := tview.NewTableCell(h)
		cell.SetTextColor(yellow)
		cell.SetBackgroundColor(blue)
		cell.SetAlign(tview.AlignCenter)
		cell.SetExpansion(t.widths[c])
		cell.NotSelectable = true
		t.t.SetCell(0, c, cell)
	}
	t.t.SetFixed(1, 0)
	t.t.Select(1, 0)
}

func (t *table) appendRow(cols ...string) {
	row := t.t.GetRowCount()
	for col, val := range cols {
		cell := tview.NewTableCell(t.adjustLength(val, col))
		cell.SetAlign(int(t.aligns[col]))
		cell.SetExpansion(t.widths[col])
		t.t.SetCell(row, col, cell)
	}
}

func (t *table) clear() {
	t.t.Clear()
}

// TODO - implement more accurate calculation
func (t *table) adjustLength(val string, col int) string {
	// _, _, w, _ := t.t.GetRect()                            // table weight
	// m := float32(w) / float32(t.allWidths) * 1.3           // multiplier
	// val = utils.FirstN(val, int(m*float32(t.widths[col]))) // cut string
	return val
}

// //////////////////////////////////////////////////////////////
// tview.Table wrapper (vertical layout)
// //////////////////////////////////////////////////////////////
type infoTable struct {
	t         *tview.Table
}

func newInfoTable() *infoTable {
	t := &infoTable{}
	t.t = tview.NewTable()
	t.t.SetSelectable(false, false)
	t.t.SetBorder(false)
	return t
}

func (t *infoTable) appendRow(label string, value string) {
	row := t.t.GetRowCount()
	// label
	labelCell := tview.NewTableCell(" " + label)
	labelCell.SetTextColor(yellow)
	t.t.SetCell(row, 0, labelCell)
	// value
	valueCell := tview.NewTableCell(" " + value)
	t.t.SetCell(row, 1, valueCell)
}

func (t *infoTable) clear() {
	t.t.Clear()
}

// //////////////////////////////////////////////////////////////
// tview.Form wrapper
// //////////////////////////////////////////////////////////////
type form struct {
	f  *tview.Form
	mu sync.Mutex
}

func newForm() *form {
	f := &form{}
	f.f = tview.NewForm()
	f.f.SetFieldTextColor(black)
	// f.f.SetFieldBackgroundColor(black)
	f.f.SetButtonTextColor(black)
	// f.f.SetButtonBackgroundColor()
	return f
}

func (f *form) SetHorizontal(b bool) {
	f.f.SetHorizontal(b)
}

func (f *form) SetTitle(t string) {
	f.f.SetTitle(" " + t  + " ")
}

func (f *form) AddInputField(label, value string, fieldWidth int, accept func(textToCheck string, lastChar rune) bool, changed func(text string)) *tview.InputField {
	f.mu.Lock()
	f.f.AddInputField(label, value, fieldWidth, accept, changed)
	// return just created input field
	obj := f.f.GetFormItem(f.f.GetFormItemCount() - 1).(*tview.InputField)
	f.mu.Unlock()
	return obj
}

func (f *form) AddButton(label string, selected func()) *tview.Button {
	f.mu.Lock()
	f.f.AddButton(label, selected)
	// return just created button
	obj := f.f.GetButton(f.f.GetButtonCount() - 1)
	f.mu.Unlock()
	return obj
}

func (f *form) AddCheckbox(label string, checked bool, changed func(checked bool)) *tview.Checkbox {
	f.mu.Lock()
	f.f.AddCheckbox(label, checked, changed)
	// return just created checkbox
	obj := f.f.GetFormItem(f.f.GetFormItemCount() - 1).(*tview.Checkbox)
	f.mu.Unlock()
	return obj
}

func (f *form) AddDropdown(label string, options []string, initialOption int, selected func(option string, optionIndex int)) *tview.DropDown {
	f.mu.Lock()
	f.f.AddDropDown(label, options, initialOption, selected)
	// return just created dropdown
	obj := f.f.GetFormItem(f.f.GetFormItemCount() - 1).(*tview.DropDown)
	f.mu.Unlock()
	return obj
}

func (f *form) AddPasswordField(label, value string, fieldWidth int, mask rune, changed func(text string)) *tview.InputField {
	f.mu.Lock()
	f.f.AddPasswordField(label, value, fieldWidth, mask, changed)
	// return just created InputField
	obj := f.f.GetFormItem(f.f.GetFormItemCount() - 1).(*tview.InputField)
	f.mu.Unlock()
	return obj
}

func (f *form) AddTextArea(label, text string, fieldWidth, fieldHeight, maxLength int, changed func(text string)) *tview.TextArea {
	f.mu.Lock()
	f.f.AddTextArea(label, text, fieldWidth, fieldHeight, maxLength, changed)
	// return just created InputField
	obj := f.f.GetFormItem(f.f.GetFormItemCount() - 1).(*tview.TextArea)
	f.mu.Unlock()
	return obj
}

func (f *form) AddTextView(label, text string, fieldWidth, fieldHeight int, dynamicColors, scrollable bool) *tview.TextView {
	f.mu.Lock()
	f.f.AddTextView(label, text, fieldWidth, fieldHeight, dynamicColors, scrollable)
	// return just created InputField
	obj := f.f.GetFormItem(f.f.GetFormItemCount() - 1).(*tview.TextView)
	f.mu.Unlock()
	return obj
}

func (f *form) AddFormItem(item tview.FormItem) *tview.FormItem {
	f.mu.Lock()
	f.f.AddFormItem(item)
	// return just created FormItem
	obj := f.f.GetFormItem(f.f.GetFormItemCount() - 1)
	f.mu.Unlock()
	return &obj
}

// //////////////////////////////////////////////////////////////
// tview.TextView wrapper
// //////////////////////////////////////////////////////////////
func newText(text string) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetTextAlign(tview.AlignCenter)
	tv.SetText(text)
	return tv
}

// //////////////////////////////////////////////////////////////
// tview.Button wrapper
// //////////////////////////////////////////////////////////////
func newButton(text string, f Fn) *tview.Button {
	bt := tview.NewButton(text)
	bt.SetRect(0, 0, 20, 1)
	bt.SetBorder(false)
	bt.SetSelectedFunc(f)
	return bt
}

// //////////////////////////////////////////////////////////////
// tview.Box wrapper
// //////////////////////////////////////////////////////////////
func box(title string) *tview.Box {
	b := tview.NewBox()
	b.SetBorder(true)
	b.SetTitle(title)
	return b
}
