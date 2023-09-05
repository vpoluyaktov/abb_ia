package ui

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
	"github.com/vpoluyaktov/abb_ia/internal/utils"
)

type ChaptersPage struct {
	mq                *mq.Dispatcher
	grid              *tview.Grid
	author            *tview.InputField
	title             *tview.InputField
	cover             *tview.InputField
	descriptionEditor *tview.TextArea
	chaptersSection   *tview.Grid
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

	// Ignore mouse events when has no focus
	p.grid.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if p.grid.HasFocus() {
			return action, event
		} else {
			return action, nil
		}
	})

	// book info section
	infoSection := tview.NewGrid()
	infoSection.SetColumns(50, -1, 30)
	infoSection.SetRows(5, 3)
	infoSection.SetBorder(true)
	infoSection.SetTitle(" Audiobook information: ")
	infoSection.SetTitleAlign(tview.AlignLeft)
	f0 := newForm()
	f0.SetBorderPadding(1, 0, 1, 1)
	p.author = f0.AddInputField("Author:", "", 40, nil, func(s string) { p.ab.Author = s })
	p.title = f0.AddInputField("Title:", "", 40, nil, func(s string) { p.ab.Title = s })
	infoSection.AddItem(f0.f, 0, 0, 1, 1, 0, 0, true)
	f1 := newForm()
	f1.SetBorderPadding(0, 1, 1, 1)
	p.cover = f1.AddInputField("Book cover:", "", 0, nil, func(s string) { p.ab.IAItem.Cover = s })
	infoSection.AddItem(f1.f, 1, 0, 1, 2, 0, 0, true)
	f3 := newForm()
	f3.SetHorizontal(false)
	f3.SetButtonsAlign(tview.AlignRight)
	f3.AddButton("Create Book", p.buildBook)
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
	f4.SetBorder(true)
	f4.SetHorizontal(true)

	f4.AddInputField("Search: ", "", 30, nil, func(s string) { p.ab.Author = s })
	f4.AddInputField("Replace:", "", 30, nil, func(s string) { p.ab.Author = s })
	f4.AddButton("Replace", p.buildBook)
	f4.AddButton(" Undo  ", p.stopConfirmation)
	f4.SetButtonsAlign(tview.AlignRight)
	descriptionSection.AddItem(f4.f, 0, 1, 1, 1, 0, 0, true)
	p.grid.AddItem(descriptionSection, 1, 0, 1, 1, 0, 0, true)

	// chapters section
	p.chaptersSection = tview.NewGrid()
	p.chaptersSection.SetColumns(-1, 40)
	p.chaptersTable = newTable()
	p.chaptersTable.SetBorder(true)
	p.chaptersTable.SetTitle(" Book chapters: ")
	p.chaptersTable.SetTitleAlign(tview.AlignLeft)
	p.chaptersTable.setHeaders("  # ", "Start", "End", "Duration", "Chapter name")
	p.chaptersTable.setWeights(1, 1, 1, 1, 20)
	p.chaptersTable.setAlign(tview.AlignRight, tview.AlignRight, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.chaptersTable.SetSelectedFunc(p.updateChapterEntry)
	p.chaptersTable.SetMouseDblClickFunc(p.updateChapterEntry)
	p.chaptersSection.AddItem(p.chaptersTable.t, 0, 0, 1, 1, 0, 0, true)
	f5 := newForm()
	f5.SetBorder(true)
	f5.SetHorizontal(true)
	f5.AddInputField("Search: ", "", 30, nil, func(s string) { p.ab.Author = s })
	f5.AddInputField("Replace:", "", 30, nil, func(s string) { p.ab.Author = s })
	f5.AddButton("Replace", p.buildBook)
	f5.AddButton(" Undo  ", p.stopConfirmation)
	f5.SetButtonsAlign(tview.AlignRight)
	f5.SetMouseDblClickFunc(func() {})
	p.chaptersSection.AddItem(f5.f, 0, 1, 1, 1, 0, 0, false)
	p.grid.AddItem(p.chaptersSection, 2, 0, 1, 1, 0, 0, true)

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
	case *dto.AddChapterCommand:
		p.addChapter(dto.Chapter)
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
	p.chaptersTable.t.ScrollToBeginning()
	p.mq.SendMessage(mq.EncodingPage, mq.TUI, &dto.SetFocusCommand{Primitive: p.chaptersSection}, false)
}

func (p *ChaptersPage) addChapter(chapter *dto.Chapter) {
	number := strconv.Itoa(chapter.Number)
	startH, _ := utils.SecondsToTime(chapter.Start)
	endH, _ := utils.SecondsToTime(chapter.End)
	durationH, _ := utils.SecondsToTime(chapter.Duration)

	p.chaptersTable.appendRow(number, startH, endH, durationH, chapter.Name)
	p.chaptersTable.t.ScrollToBeginning()
}

func (p *ChaptersPage) updateChapterEntry(row int, col int) {
	chapter := p.ab.Chapters[row-1]
	durationH, _ := utils.SecondsToTime(chapter.Duration)
	d := newDialogWindow(p.mq, 11, 80, p.chaptersSection)
	f := newForm()
	f.SetTitle("Update Chapter Name:")
	f.AddTextView("Chapter #:  ", strconv.Itoa(chapter.Number), 5, 1, true, false)
	f.AddTextView("Duration:   ", strings.TrimLeft(durationH, " "), 10, 1, true, false)
	nameF := f.AddInputField("Chapter name:", chapter.Name, 60, nil, nil)
	f.AddButton("Save changes", func() {
		cell := p.chaptersTable.t.GetCell(row, col)
		cell.Text = nameF.GetText()
		p.ab.Chapters[row-1].Name = nameF.GetText()
		d.Close()
	})
	f.AddButton("Cancel", func() {
		d.Close()
	})
	d.setForm(f.f)
	d.Show()
}

func (p *ChaptersPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop editing chapters?", p.chaptersSection, p.stopChapters, func() {})
}

func (p *ChaptersPage) stopChapters() {
	// Stop the chapters here
	p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.StopCommand{Process: "Chapters", Reason: "User request"}, false)
	p.mq.SendMessage(mq.ChaptersPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

func (p *ChaptersPage) buildBook() {
	p.mq.SendMessage(mq.ChaptersPage, mq.BuildPage, &dto.DisplayBookInfoCommand{Audiobook: p.ab}, true)
	p.mq.SendMessage(mq.ChaptersPage, mq.BuildController, &dto.BuildCommand{Audiobook: p.ab}, true)
	p.mq.SendMessage(mq.ChaptersPage, mq.Frame, &dto.SwitchToPageCommand{Name: "BuildPage"}, false)
}
