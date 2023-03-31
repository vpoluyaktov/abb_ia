package ui

import (
	"code.rocketnine.space/tslocum/cview"
)

type TUI struct {
	// View components.
	app     *cview.Application
	colors  *colors
	header  *header
	pannels *cview.Panels
	footer  *footer
}

type Fn func()

func NewTUI() *TUI {

	ui := TUI{}
	ui.app = cview.NewApplication()
	defer ui.app.HandlePanic()
	ui.colors = newColors()
	ui.app.EnableMouse(true)

	// UI components
	ui.header = newHeader(ui.colors)
	ui.footer = newFooter(ui.colors)

	newText := func(text string) *cview.TextView {
		tv := cview.NewTextView()
		tv.SetTextAlign(cview.AlignCenter)
		tv.SetText(text)
		return tv
	}

	newButton := func(text string, f Fn) *cview.Button {
		bt := cview.NewButton(text)
		bt.SetRect(0, 0, 22, 3)
		bt.SetBorder(true)
		bt.SetSelectedFunc(f)
		return bt
	}

	ui.pannels = cview.NewPanels()

	backGrid := cview.NewGrid()
	backGrid.SetRows(1, 0, 1)
	backGrid.SetColumns(1, 0, 1)
	backGrid.AddItem(ui.header.view, 0, 0, 1, 3, 0, 0, false)
	backGrid.AddItem(ui.footer.view, 2, 0, 1, 3, 0, 0, false)
	ui.pannels.AddPanel("BackPanel", backGrid, true, true)

	intPannels := cview.NewPanels()

	searchGrid := cview.NewGrid()
	// searchGrid.SetBorder(true)
	// searchGrid.SetBorders(true)
	searchGrid.SetRows(1, 1, 0, 3)
	searchGrid.SetColumns(0)
	searchGrid.AddItem(newText("SearchPanel"), 0, 0, 1, 1, 0, 0, true)
	searchGrid.AddItem(newText("Body Search"), 1, 0, 1, 1, 0, 0, true)
	searchGrid.AddItem(newButton("NextPanel", func() { intPannels.SetCurrentPanel("ProcessingPanel") }), 3, 0, 1, 1, 0, 0, true)
	intPannels.AddPanel("SearchPanel", searchGrid, true, true)

	processingGrid := cview.NewGrid()
	processingGrid.SetRows(1, 1, 0, 3)
	processingGrid.SetColumns(0)
	processingGrid.AddItem(newText("ProcessingPanel"), 0, 0, 1, 1, 0, 0, true)
	processingGrid.AddItem(newText("Processing..."), 1, 0, 1, 1, 0, 0, true)
	processingGrid.AddItem(newButton("NextPanel", func() { intPannels.SetCurrentPanel("UploadPanel") }), 3, 0, 1, 1, 0, 0, true)
	intPannels.AddPanel("ProcessingPanel", processingGrid, true, true)

	uploadGrid := cview.NewGrid()
	uploadGrid.SetRows(1, 1, 0, 3)
	uploadGrid.SetColumns(0)
	uploadGrid.AddItem(newText("UploadPanel"), 0, 0, 1, 1, 0, 0, true)
	uploadGrid.AddItem(newText("Uploading..."), 1, 0, 1, 1, 0, 0, true)
	uploadGrid.AddItem(newButton("NextPanel", func() { intPannels.SetCurrentPanel("SearchPanel") }), 3, 0, 1, 1, 0, 0, true)
	intPannels.AddPanel("UploadPanel", uploadGrid, true, true)

	backGrid.AddItem(intPannels, 1, 1, 1, 1, 0, 0, false)
	ui.pannels.SetCurrentPanel("BackPanel")
	intPannels.SetCurrentPanel("SearchPanel")

	ui.app.SetRoot(backGrid, true)
	return &ui
}

func (tui *TUI) Run() {
	if err := tui.app.Run(); err != nil {
		panic(err)
	}
}
