package ui

import (
	"math"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// //////////////////////////////////////////////////////////////
// tview.Table wrapper
// //////////////////////////////////////////////////////////////
type table struct {
	t         *tview.Table
	headers   []string
	colWeight []int
	colWidth  []int
	aligns    []uint
}

func newTable() *table {
	t := &table{}
	t.t = tview.NewTable()
	t.t.SetDrawFunc(t.draw)
	t.t.SetSelectable(true, false)
	t.t.SetSeparator(tview.Borders.Vertical)
	// t.t.SetSortClicked(false)
	// t.t.SetSortFunc() // TODO implement sorting
	// t.t.ShowFocus(true)
	t.t.SetBorder(false)
	t.t.Clear()
	// t.t.SetEvaluateAllRows(true)
	return t
}

// Enter Key
func (t *table) SetSelectedFunc(f func(row, column int)) {
	t.t.SetSelectedFunc(f)
}

// Mouse Double Click
func (t *table) SetMouseDblClickFunc(f func(row, column int)) {
	t.t.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseLeftDoubleClick:
			{
				f(t.t.GetSelection())
			}
		}
		return action, event
	})
}

func (t *table) SetBorder(b bool) {
	t.t.SetBorder(b)
}

func (t *table) SetTitle(s string) {
	t.t.SetTitle(s)
}

func (t *table) SetTitleAlign(a int) {
	t.t.SetTitleAlign(a)
}

func (t *table) draw(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	t.recalculateColumnWidths()
	return t.t.GetInnerRect()
}

func (t *table) setHeaders(headers ...string) {
	t.headers = headers
}

func (t *table) setWeights(weights ...int) {
	t.colWeight = weights
	t.recalculateColumnWidths()
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
		cell.SetExpansion(t.colWeight[c])
		cell.SetMaxWidth(t.colWidth[c])
		cell.NotSelectable = true
		t.t.SetCell(0, c, cell)
	}
	t.t.SetFixed(1, 0)
	t.t.Select(1, 0)
}

func (t *table) appendRow(cols ...string) {
	row := t.t.GetRowCount()
	for col, val := range cols {
		cell := tview.NewTableCell(val)
		cell.SetAlign(int(t.aligns[col]))
		cell.SetExpansion(t.colWeight[col])
		cell.SetMaxWidth(t.colWidth[col])
		t.t.SetCell(row, col, cell)
	}
}

func (t *table) appendSeparator(cols ...string) {
	row := t.t.GetRowCount()
	for col, val := range cols {
		cell := tview.NewTableCell(val)
		cell.SetAlign(int(t.aligns[col]))
		cell.SetExpansion(t.colWeight[col])
		cell.SetMaxWidth(t.colWidth[col])
		// cell.NotSelectable = true
		cell.SetTextColor(tview.Styles.PrimaryTextColor)
		cell.SetBackgroundColor(lightBlue)
		t.t.SetCell(row, col, cell)
	}
}

func (t *table) ScrollToBeginning() {
	t.t.ScrollToBeginning()
}

func (t *table) clear() {
	t.t.Clear()
}

// TODO - implement more accurate calculation
func (t *table) recalculateColumnWidths() {
	if len(t.colWeight) == 0 {
		return
	}
	allWeights := 0
	for _, w := range t.colWeight {
		allWeights += w
	}
	_, _, tw, _ := t.t.GetInnerRect()                           // table weight
	m := (float64(tw-len(t.colWeight)-1) / float64(allWeights)) // multiplier

	t.colWidth = make([]int, len(t.colWeight))
	for c, _ := range t.colWidth {
		t.colWidth[c] = int(math.Round(m * float64(t.colWeight[c])))
	}
}

func (t *table) getColumnWidth(col int) int {
	if len(t.colWidth) == 0 {
		t.recalculateColumnWidths()
	}
	return t.colWidth[col]
}

func (t *table) GetRowCount() int {
	return t.t.GetColumnCount()
}

// //////////////////////////////////////////////////////////////
// tview.Table wrapper (vertical layout)
// //////////////////////////////////////////////////////////////
type infoPanel struct {
	t *tview.Table
}

func newInfoPanel() *infoPanel {
	p := &infoPanel{}
	p.t = tview.NewTable()
	p.t.SetSelectable(false, false)
	p.t.SetBorder(false)
	return p
}

func (p *infoPanel) appendRow(label string, value string) {
	row := p.t.GetRowCount()
	// label
	labelCell := tview.NewTableCell(" " + label)
	labelCell.SetTextColor(yellow)
	p.t.SetCell(row, 0, labelCell)
	// value
	valueCell := tview.NewTableCell(" " + value)
	p.t.SetCell(row, 1, valueCell)
}

func (p *infoPanel) clear() {
	p.t.Clear()
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
	f.f.SetTitle(" " + t + " ")
}

func (f *form) SetBorder(b bool) {
	f.f.SetBorder(b)
}

func (f *form) SetButtonsAlign(a int) {
	f.f.SetButtonsAlign(a)
}

func (f *form) SetBorderPadding(top int, bottom int, left int, right int) {
	f.f.SetBorderPadding(top, bottom, left, right)
}

// Mouse Double Click
func (f *form) SetMouseDblClickFunc(fn func()) {
	f.f.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseLeftDoubleClick:
			{
				fn()
			}
		}
		return action, event
	})
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
func newTextView(text string) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetTextAlign(tview.AlignCenter)
	tv.SetText(text)
	return tv
}

// //////////////////////////////////////////////////////////////
// tview.TextArea wrapper
// //////////////////////////////////////////////////////////////
func newTextArea(text string) *tview.TextArea {
	ta := tview.NewTextArea()
	ta.SetText(text, false)
	return ta
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
