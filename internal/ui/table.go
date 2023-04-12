package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
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
		textColor := tcell.ColorYellow.TrueColor()
		bgColor := tcell.ColorBlue.TrueColor()
		cell.SetTextColor(textColor)
		cell.SetBackgroundColor(bgColor)
		cell.SetAlign(int(t.aligns[c]))
		cell.SetExpansion(t.widths[c])
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
	// p.searchResultTable.Select(r, 0)
}

func (t *table) clear() {
	t.t.Clear()
}
