package ui

import (
	"code.rocketnine.space/tslocum/cview"
)

type frame struct {
	grid    *cview.Grid
	pannels *cview.Panels
}

func newFrame() *frame {
	f := &frame{}
	f.grid = cview.NewGrid()
	f.grid.SetRows(1, 0, 1)
	f.grid.SetColumns(1, 0, 1)

	return f
}

func (f *frame) addHeader(header *header) {
	f.grid.AddItem(header.view, 0, 1, 1, 1, 0, 0, false)
}

func (f *frame) addFooter(footer *footer) {
	f.grid.AddItem(footer.view, 2, 1, 1, 1, 0, 0, false)
}

func (f *frame) addPannel(name string, g *cview.Grid) {
	if f.pannels == nil {
		f.pannels = cview.NewPanels()
		f.grid.AddItem(f.pannels, 1, 1, 1, 1, 0, 0, false)
	}
	f.pannels.AddPanel(name, g, true, true)
}
