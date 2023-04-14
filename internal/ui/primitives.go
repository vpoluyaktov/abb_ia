package ui

import (
	"code.rocketnine.space/tslocum/cview"
)

// //////////////////////////////////////////////////////////////
// cview.Table wrapper
// //////////////////////////////////////////////////////////////
type table struct {
	t       *cview.Table
	headers []string
	widths  []int
	aligns  []uint
}

func newTable() *table {
	t := &table{}
	t.t = cview.NewTable()
	t.t.SetSelectable(true, false)
	t.t.SetSeparator(cview.Borders.Vertical)
	t.t.SetSortClicked(false)
	// t.t.SetSortFunc() // TODO
	t.t.ShowFocus(true)
	t.t.SetBorder(false)
	return t
}

func (t *table) setHeaders(headers ...string) {
	t.headers = headers
}

func (t *table) setWidths(widths ...int) {
	t.widths = widths
}

func (t *table) setAlign(aligns ...uint) {
	t.aligns = aligns
}

func (t *table) showHeader() {
	for c, h := range t.headers {
		cell := cview.NewTableCell(h)
		cell.SetTextColor(yellow)
		cell.SetBackgroundColor(blue)
		cell.SetAlign(cview.AlignCenter)
		cell.SetExpansion(t.widths[c])
		cell.NotSelectable = true
		t.t.SetCell(0, c, cell)
	}
	t.t.SetFixed(1, 0)
	t.t.Select(1, 0)
}

func (t *table) appendRow(cols ...string) {
	r := t.t.GetRowCount()
	for c, col := range cols {
		cell := cview.NewTableCell(col)
		cell.SetAlign(int(t.aligns[c]))
		cell.SetExpansion(t.widths[c])
		t.t.SetCell(r, c, cell)
	}
}

func (t *table) clear() {
	t.t.Clear()
}

// //////////////////////////////////////////////////////////////
// cview.Form wrapper
// //////////////////////////////////////////////////////////////
type form struct {
	f *cview.Form
}

func newForm() *form {
	f := &form{}
	f.f = cview.NewForm()
	f.f.SetFieldTextColor(black)
	f.f.SetFieldTextColorFocused(black)
	f.f.SetButtonTextColor(black)
	f.f.SetButtonTextColorFocused(black)
	return f
}

func (f *form) SetHorizontal(b bool) {
	f.f.SetHorizontal(b)
}

func (f *form) AddInputField(label, value string, fieldWidth int, accept func(textToCheck string, lastChar rune) bool, changed func(text string)) (*cview.InputField) {
	f.f.AddInputField(label, value, fieldWidth, accept, changed)
	// return just created input field
	return f.f.GetFormItem(f.f.GetFormItemCount()-1).(*cview.InputField)
}

func (f *form) AddButton(label string, selected func()) (*cview.Button) {
	f.f.AddButton(label, selected)
	// return just created button
	return f.f.GetButton(f.f.GetButtonCount()-1)
}

func newText(text string) *cview.TextView {
	tv := cview.NewTextView()
	tv.SetTextAlign(cview.AlignCenter)
	tv.SetText(text)
	return tv
}

func newButton(text string, f Fn) *cview.Button {
	bt := cview.NewButton(text)
	bt.SetRect(0, 0, 20, 1)
	bt.SetBorder(false)
	bt.SetSelectedFunc(f)
	return bt
}

func box(title string) *cview.Box {
	b := cview.NewBox()
	b.SetBorder(true)
	b.SetTitle(title)
	return b
}
