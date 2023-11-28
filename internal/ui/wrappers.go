package ui

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// // //////////////////////////////////////////////////////////////
// // tview.Grid wrapper
// // //////////////////////////////////////////////////////////////
type grid struct {
	*tview.Grid
	navigationOrder []tview.Primitive
}

func newGrid() *grid {
	g := &grid{}
	g.Grid = tview.NewGrid()
	g.navigationOrder = []tview.Primitive{}

	// Ignore mouse events when the grid has no focus
	g.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if g.HasFocus() {
			return action, event
		} else if len(g.navigationOrder) == 0 {
			return action, event
		} else {
			return action, nil
		}
	})

	// Tab / Backtab navigation
	g.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if g.Grid.HasFocus() {
			switch event.Key() {
			case tcell.KeyTab:
				if len(g.navigationOrder) == 0 {
					return event
				}
				currentElement, err := g.getFocusIndex()
				if err != nil {
					return event
				}
				nextElement := currentElement + 1
				if nextElement > len(g.navigationOrder)-1 {
					nextElement = 0
				}
				nextPrimitive := g.navigationOrder[nextElement]
				ui.SetFocus(nextPrimitive)
				ui.Draw()
				return nil

			case tcell.KeyBacktab:
				if len(g.navigationOrder) == 0 {
					return event
				}
				currentElement, err := g.getFocusIndex()
				if err != nil {
					return event
				}
				previousElement := currentElement - 1
				if previousElement < 0 {
					previousElement = len(g.navigationOrder) - 1
				}
				previousPrimitive := g.navigationOrder[previousElement]
				ui.SetFocus(previousPrimitive)
				ui.Draw()
				return nil
			default:
				return event
			}
		} else {
			return event
		}
	})

	g.SetFocusFunc(func() {
		// g.SetBorderColor(yellow)
	})

	g.SetBlurFunc(func() {
		// g.SetBorderColor(black)
	})

	return g
}

func (g *grid) SetNavigationOrder(navigationOrder ...tview.Primitive) {
	g.navigationOrder = navigationOrder
}

// return an index of the focus element in the navigation list
func (g *grid) getFocusIndex() (int, error) {
	index := 0
	found := false
	focus := ui.GetFocus()
	for i, v := range g.navigationOrder {
		if focus == v {
			found = true
			index = i
			break
		}
	}
	if found {
		return index, nil
	} else {
		return index, fmt.Errorf("focus element not found")
	}
}

// //////////////////////////////////////////////////////////////
// tview.Table wrapper
// //////////////////////////////////////////////////////////////
type table struct {
	*tview.Table
	headers   []string
	colWeight []int
	colWidth  []int
	aligns    []uint
}

func newTable() *table {
	t := &table{}
	t.Table = tview.NewTable()
	t.Table.SetDrawFunc(t.draw)
	t.Table.SetSelectable(true, false)
	t.Table.SetSeparator(tview.Borders.Vertical)
	// t.Table.SetSortClicked(false)
	// t.Table.SetSortFunc() // TODO implement sorting??
	// t.Table.ShowFocus(true)
	t.Table.SetBorder(false)
	t.Table.Clear()
	// t.Table.SetEvaluateAllRows(true)
	return t
}

// Mouse Double Click
func (t *table) SetMouseDblClickFunc(f func(row, column int)) {
	t.Table.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseLeftDoubleClick:
			{
				f(t.Table.GetSelection())
			}
		}
		return action, event
	})
}

func (t *table) draw(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	t.recalculateColumnWidths()
	return t.Table.GetInnerRect()
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
		t.SetCell(0, c, cell)
	}
	t.SetFixed(1, 0)
	t.Select(1, 0)
}

func (t *table) appendRow(cols ...string) {
	row := t.GetRowCount()
	for col, val := range cols {
		cell := tview.NewTableCell(val)
		cell.SetAlign(int(t.aligns[col]))
		cell.SetExpansion(t.colWeight[col])
		cell.SetMaxWidth(t.colWidth[col])
		t.SetCell(row, col, cell)
	}
}

func (t *table) appendSeparator(cols ...string) {
	row := t.GetRowCount()
	for col, val := range cols {
		cell := tview.NewTableCell(val)
		cell.SetAlign(int(t.aligns[col]))
		cell.SetExpansion(t.colWeight[col])
		cell.SetMaxWidth(t.colWidth[col])
		// cell.NotSelectable = true
		cell.SetTextColor(tview.Styles.PrimaryTextColor)
		cell.SetBackgroundColor(lightBlue)
		t.SetCell(row, col, cell)
	}
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
	_, _, tw, _ := t.Table.GetInnerRect()                       // table weight
	m := (float64(tw-len(t.colWeight)-1) / float64(allWeights)) // multiplier

	t.colWidth = make([]int, len(t.colWeight))
	for c := range t.colWidth {
		t.colWidth[c] = int(math.Round(m * float64(t.colWeight[c])))
	}
}

func (t *table) getColumnWidth(col int) int {
	if len(t.colWidth) == 0 {
		t.recalculateColumnWidths()
	}
	return t.colWidth[col]
}

// //////////////////////////////////////////////////////////////
// tview.Table wrapper (vertical layout)
// //////////////////////////////////////////////////////////////
type infoPanel struct {
	*tview.Table
}

func newInfoPanel() *infoPanel {
	p := &infoPanel{}
	p.Table = tview.NewTable()
	p.SetSelectable(false, false)
	p.SetBorder(false)
	return p
}

func (p *infoPanel) appendRow(label string, value string) {
	row := p.GetRowCount()
	// label
	labelCell := tview.NewTableCell(" " + label)
	labelCell.SetTextColor(yellow)
	p.SetCell(row, 0, labelCell)
	// value
	valueCell := tview.NewTableCell(" " + value)
	p.SetCell(row, 1, valueCell)
}

func (p *infoPanel) clear() {
	p.Clear()
}

// //////////////////////////////////////////////////////////////
// tview.Form wrapper
// //////////////////////////////////////////////////////////////
type form struct {
	*tview.Form
	mu sync.Mutex
}

func newForm() *form {
	f := &form{}
	f.Form = tview.NewForm()
	f.SetFieldTextColor(black)
	// f.f.SetFieldBackgroundColor(black)
	f.SetButtonTextColor(black)
	// f.f.SetButtonBackgroundColor()
	return f
}

func (f *form) SetTitle(t string) {
	f.Form.SetTitle(" " + t + " ")
}

// Mouse Double Click
func (f *form) SetMouseDblClickFunc(fn func()) {
	f.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
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
	// right aligment for numberic values
	// if utils.IsNumber(value) {
	// 	value = fmt.Sprintf("%*s", fieldWidth, value)
	// }

	f.Form.AddInputField(label, value, fieldWidth, accept, changed)
	// return just created input field
	obj := f.GetFormItem(f.GetFormItemCount() - 1).(*tview.InputField)
	f.mu.Unlock()
	return obj
}

func acceptInt(textToCheck string, lastChar rune) bool {
	_, err := strconv.Atoi(textToCheck)
	return err == nil
}

func (f *form) AddButton(label string, selected func()) *tview.Button {
	f.mu.Lock()
	f.Form.AddButton(label, selected)
	// return just created button
	obj := f.GetButton(f.GetButtonCount() - 1)
	f.mu.Unlock()
	return obj
}

func (f *form) AddCheckbox(label string, checked bool, changed func(checked bool)) *tview.Checkbox {
	f.mu.Lock()
	f.Form.AddCheckbox(label, checked, changed)
	// return just created checkbox
	obj := f.GetFormItem(f.GetFormItemCount() - 1).(*tview.Checkbox)
	f.mu.Unlock()
	return obj
}

func (f *form) AddDropdown(label string, options []string, initialOption int, selected func(option string, optionIndex int)) *tview.DropDown {
	f.mu.Lock()
	f.AddDropDown(label, options, initialOption, selected)
	// return just created dropdown
	obj := f.GetFormItem(f.GetFormItemCount() - 1).(*tview.DropDown)
	f.mu.Unlock()
	return obj
}

func (f *form) AddPasswordField(label, value string, fieldWidth int, mask rune, changed func(text string)) *tview.InputField {
	f.mu.Lock()
	f.Form.AddPasswordField(label, value, fieldWidth, mask, changed)
	// return just created InputField
	obj := f.GetFormItem(f.GetFormItemCount() - 1).(*tview.InputField)
	f.mu.Unlock()
	return obj
}

func (f *form) AddTextArea(label, text string, fieldWidth, fieldHeight, maxLength int, changed func(text string)) *tview.TextArea {
	f.mu.Lock()
	f.Form.AddTextArea(label, text, fieldWidth, fieldHeight, maxLength, changed)
	// return just created InputField
	obj := f.GetFormItem(f.GetFormItemCount() - 1).(*tview.TextArea)
	f.mu.Unlock()
	return obj
}

func (f *form) AddTextView(label, text string, fieldWidth, fieldHeight int, dynamicColors, scrollable bool) *tview.TextView {
	f.mu.Lock()
	f.Form.AddTextView(label, text, fieldWidth, fieldHeight, dynamicColors, scrollable)
	// return just created InputField
	obj := f.GetFormItem(f.GetFormItemCount() - 1).(*tview.TextView)
	f.mu.Unlock()
	return obj
}

func (f *form) AddFormItem(item tview.FormItem) *tview.FormItem {
	f.mu.Lock()
	f.Form.AddFormItem(item)
	// return just created FormItem
	obj := f.GetFormItem(f.GetFormItemCount() - 1)
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
