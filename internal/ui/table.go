package ui

import (
	"code.rocketnine.space/tslocum/cview"
)

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

func (t *table) setHeaders(headers []string) {
	t.headers = headers
}

func (t *table) setWidths(widths []int) {
	t.widths = widths
}

func (t *table) setAlign(aligns []uint) {
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

func (t *table) appendRow(cols []string) {
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
