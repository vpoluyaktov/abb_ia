package ui

import (
	"strconv"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type ChaptersPage struct {
	mq                *mq.Dispatcher
	grid              *tview.Grid
	author            *tview.InputField
	title             *tview.InputField
	cover             *tview.InputField
	descriptionEditor *tview.TextArea
	chaptersTable     *table
	ab                *dto.Audiobook
}

func newChaptersPage(dispatcher *mq.Dispatcher) *ChaptersPage {
	p := &ChaptersPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.ChaptersPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(9, -1, -1)
	p.grid.SetColumns(0)

	// book info section
	infoSection := tview.NewGrid()
	infoSection.SetColumns(50, -1, 30)
	infoSection.SetRows(5, 3)
	infoSection.SetBorder(true)
	infoSection.SetTitle(" Audiobook information: ")
	infoSection.SetTitleAlign(tview.AlignLeft)
	f0 := newForm()
	f0.f.SetBorderPadding(1, 0, 1, 1)
	p.author = f0.AddInputField("Author:", "", 40, nil, func(s string) { p.ab.Author = s })
	p.title = f0.AddInputField("Title:", "", 40, nil, func(s string) { p.ab.Title = s })
	infoSection.AddItem(f0.f, 0, 0, 1, 1, 0, 0, true)
	f1 := newForm()
	f1.f.SetBorderPadding(0, 1, 1, 1)
	p.cover = f1.AddInputField("Book cover:", "", 0, nil, func(s string) { p.ab.IAItem.Cover = s })
	infoSection.AddItem(f1.f, 1, 0, 1, 2, 0, 0, true)
	f3 := newForm()
	f3.SetHorizontal(false)
	f3.f.SetButtonsAlign(tview.AlignRight)
	f3.AddButton("Create Book", p.createBook)
	f3.AddButton("Cancel", p.stopConfirmation)
	infoSection.AddItem(f3.f, 0, 2, 1, 1, 0, 0, false)
	p.grid.AddItem(infoSection, 0, 0, 1, 1, 0, 0, false)

	// description section
	descriptionSection := tview.NewGrid()
	descriptionSection.SetColumns(-1, 40)
	descriptionSection.SetBorder(false)
	p.descriptionEditor = newTextArea("")
	p.descriptionEditor.SetBorder(true)
	p.descriptionEditor.SetTitle(" Book description: ")
	p.descriptionEditor.SetTitleAlign(tview.AlignLeft)
	descriptionSection.AddItem(p.descriptionEditor, 0, 0, 1, 1, 0, 0, true)
	f4 := newForm()
	f4.f.SetBorder(true)
	f4.SetHorizontal(true)

	f4.AddInputField("Search: ", "", 30, nil, func(s string) { p.ab.Author = s })
	f4.AddInputField("Replace:", "", 30, nil, func(s string) { p.ab.Author = s })
	f4.AddButton("Replace", p.createBook)
	f4.AddButton(" Undo  ", p.stopConfirmation)
	f4.f.SetButtonsAlign(tview.AlignRight)
	descriptionSection.AddItem(f4.f, 0, 1, 1, 1, 0, 0, false)
	p.grid.AddItem(descriptionSection, 1, 0, 1, 1, 0, 0, true)

	// chapters section
	chaptersSection := tview.NewGrid()
	chaptersSection.SetColumns(-1, 40)
	p.chaptersTable = newTable()
	p.chaptersTable.t.SetBorder(true)
	p.chaptersTable.t.SetTitle(" Book chapters: ")
	p.chaptersTable.t.SetTitleAlign(tview.AlignLeft)
	p.chaptersTable.setHeaders("  # ", "Duration", "Chapter name")
	p.chaptersTable.setWeights(1, 1, 20)
	p.chaptersTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft)
	chaptersSection.AddItem(p.chaptersTable.t, 0, 0, 1, 1, 0, 0, false)
	f5 := newForm()
	f5.f.SetBorder(true)
	f5.SetHorizontal(true)
	f5.AddInputField("Search: ", "", 30, nil, func(s string) { p.ab.Author = s })
	f5.AddInputField("Replace:", "", 30, nil, func(s string) { p.ab.Author = s })
	f5.AddButton("Replace", p.createBook)
	f5.AddButton(" Undo  ", p.stopConfirmation)
	f5.f.SetButtonsAlign(tview.AlignRight)
	chaptersSection.AddItem(f5.f, 0, 1, 1, 1, 0, 0, false)
	p.grid.AddItem(chaptersSection, 2, 0, 1, 1, 0, 0, false)

	return p
}

func (p *ChaptersPage) checkMQ() {
	m := p.mq.GetMessage(mq.ChaptersPage)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *ChaptersPage) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.DisplayBookInfoCommand:
		p.displayBookInfo(dto.Audiobook)
	default:
		m.UnsupportedTypeError(mq.ChaptersPage)
	}
}

func (p *ChaptersPage) displayBookInfo(ab *dto.Audiobook) {
	p.ab = ab
	p.author.SetText(ab.Author)
	p.title.SetText(ab.Title)
	p.cover.SetText(ab.IAItem.Cover)
	p.descriptionEditor.SetText(ab.IAItem.Description, false)


	p.chaptersTable.clear()
	p.chaptersTable.showHeader()
	for i, f := range ab.IAItem.Files {
		p.chaptersTable.appendRow(" "+strconv.Itoa(i+1)+" ", f.LengthH, f.Name)
	}
	p.chaptersTable.t.ScrollToBeginning()
}

func (p *ChaptersPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop editing chapters?", p.stopChapters, func() {})
}

func (p *ChaptersPage) stopChapters() {
	// Stop the chapters here
	p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.StopCommand{Process: "Chapters", Reason: "User request"}, false)
	p.mq.SendMessage(mq.ChaptersPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

func (p *ChaptersPage) createBook() {
	p.mq.SendMessage(mq.DownloadPage, mq.ChaptersPage, &dto.DisplayBookInfoCommand{Audiobook: p.ab}, true)
	p.mq.SendMessage(mq.DownloadPage, mq.ChaptersController, &dto.EncodeCommand{Audiobook: p.ab}, true)
	p.mq.SendMessage(mq.DownloadPage, mq.Frame, &dto.SwitchToPageCommand{Name: "ChaptersPage"}, false)
}
