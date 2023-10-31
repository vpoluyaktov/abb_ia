package ui

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
	"github.com/vpoluyaktov/abb_ia/internal/utils"
)

type ChaptersPage struct {
	mq                 *mq.Dispatcher
	grid               *tview.Grid
	author             *tview.InputField
	title              *tview.InputField
	series             *tview.InputField
	seriesNo           *tview.InputField
	genre              *tview.DropDown
	narator            *tview.InputField
	cover              *tview.InputField
	descriptionEditor  *tview.TextArea
	chaptersSection    *tview.Grid
	chaptersTable      *table
	ab                 *dto.Audiobook
	searchDescription  string
	replaceDescription string
	searchChapters     string
	replaceChapters    string
	chaptersUndoStack  *UndoStack
}

func newChaptersPage(dispatcher *mq.Dispatcher) *ChaptersPage {
	p := &ChaptersPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.ChaptersPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(9, -1, -1)
	p.grid.SetColumns(0)

	// Ignore mouse events when the grid has no focus
	p.grid.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if p.grid.HasFocus() {
			return action, event
		} else {
			return action, nil
		}
	})

	// book info section
	infoSection := tview.NewGrid()
	infoSection.SetColumns(50, 50, -1, 30)
	infoSection.SetRows(5, 2)
	infoSection.SetBorder(true)
	infoSection.SetTitle(" Audiobook information: ")
	infoSection.SetTitleAlign(tview.AlignLeft)
	f0 := newForm()
	f0.SetBorderPadding(1, 0, 1, 2)
	p.author = f0.AddInputField("Author:", "", 40, nil, func(s string) {
		if p.ab != nil {
			p.ab.Author = s
		}
	})
	p.title = f0.AddInputField("Title:", "", 40, nil, func(s string) {
		if p.ab != nil {
			p.ab.Title = s
		}
	})
	infoSection.AddItem(f0.f, 0, 0, 1, 1, 0, 0, true)
	f1 := newForm()
	f1.SetBorderPadding(1, 0, 2, 2)
	p.series = f1.AddInputField("Series:", "", 40, nil, func(s string) {
		if p.ab != nil {
			p.ab.Series = s
		}
	})
	p.seriesNo = f1.AddInputField("Series No:", "", 5, acceptInt, func(s string) {
		if p.ab != nil {
			p.ab.SeriesNo = s
		}
	})
	infoSection.AddItem(f1.f, 0, 1, 1, 1, 0, 0, true)
	f2 := newForm()
	f2.SetBorderPadding(1, 0, 2, 2)
	p.genre = f2.AddDropdown("Genre:", config.Genres(), 0, func(s string, i int) {
		if p.ab != nil {
			p.ab.Genre = s
		}
	})
	p.narator = f2.AddInputField("Narator:", "", 40, nil, func(s string) {
		if p.ab != nil {
			p.ab.Narator = s
		}
	})
	infoSection.AddItem(f2.f, 0, 2, 1, 1, 0, 0, true)
	f3 := newForm()
	f3.SetBorderPadding(0, 1, 1, 1)
	p.cover = f3.AddInputField("Book cover:", "", 0, nil, func(s string) {
		if p.ab != nil {
			p.ab.CoverURL = s
		}
	})
	infoSection.AddItem(f3.f, 1, 0, 1, 4, 0, 0, true)
	f4 := newForm()
	f4.SetHorizontal(false)
	f4.SetButtonsAlign(tview.AlignRight)
	f4.AddButton("Create Book", p.buildBook)
	f4.AddButton("Cancel", p.stopConfirmation)
	infoSection.AddItem(f4.f, 0, 3, 1, 1, 0, 0, false)
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
	f5 := newForm()
	f5.SetBorder(true)
	f5.SetHorizontal(true)

	f5.AddInputField("Search: ", "", 30, nil, func(s string) { p.searchDescription = s })
	f5.AddInputField("Replace:", "", 30, nil, func(s string) { p.replaceDescription = s })
	f5.AddButton("Replace", p.buildBook)
	f5.AddButton(" Undo  ", p.stopConfirmation)
	f5.SetButtonsAlign(tview.AlignRight)
	descriptionSection.AddItem(f5.f, 0, 1, 1, 1, 0, 0, true)
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
	f6 := newForm()
	f6.SetBorder(true)
	f6.SetHorizontal(true)
	f6.AddInputField("Search: ", "", 30, nil, func(s string) { p.searchChapters = s })
	f6.AddInputField("Replace:", "", 30, nil, func(s string) { p.replaceChapters = s })
	f6.AddButton("Replace", p.searchReplaceChapters)
	f6.AddButton(" Undo  ", p.undoChapters)
	f6.AddButton(" Join Chapters  ", p.joinChapters)
	f6.SetButtonsAlign(tview.AlignRight)
	f6.SetMouseDblClickFunc(func() {})
	p.chaptersSection.AddItem(f6.f, 0, 1, 1, 1, 0, 0, false)
	p.grid.AddItem(p.chaptersSection, 2, 0, 1, 1, 0, 0, true)

	p.chaptersUndoStack = NewUndoStack()

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
	case *dto.AddPartCommand:
		p.addPart(dto.Part)
	case *dto.ChaptersReady:
		p.displayParts(dto.Audiobook)
	case *dto.RefreshChaptersCommand:
		p.refreshChapters(dto.Audiobook)
	default:
		m.UnsupportedTypeError(mq.ChaptersPage)
	}
}

func (p *ChaptersPage) displayBookInfo(ab *dto.Audiobook) {
	p.ab = ab
	p.author.SetText(ab.Author)
	p.title.SetText(ab.Title)
	p.cover.SetText(ab.CoverURL)
	p.descriptionEditor.SetText(ab.IAItem.Description, false)

	p.chaptersTable.clear()
	p.chaptersTable.showHeader()
	p.chaptersTable.ScrollToBeginning()
	p.mq.SendMessage(mq.EncodingPage, mq.TUI, &dto.SetFocusCommand{Primitive: p.chaptersSection}, false)
}

func (p *ChaptersPage) displayParts(ab *dto.Audiobook) {
	if len(ab.Parts) > 1 {
		p.chaptersTable.clear()
		p.chaptersTable.showHeader()
		for _, part := range ab.Parts {
			p.addPart(&part)
		}
	}
	p.mq.SendMessage(mq.ChaptersPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, false)
}

func (p *ChaptersPage) addPart(part *dto.Part) {
	number := strconv.Itoa(part.Number)
	p.chaptersTable.appendSeparator("", "", "", "", "Part Number "+number)
	for _, chapter := range part.Chapters {
		p.addChapter(&chapter)
	}
	p.mq.SendMessage(mq.ChaptersPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, false)
}

func (p *ChaptersPage) addChapter(chapter *dto.Chapter) {
	number := strconv.Itoa(chapter.Number)
	startH := utils.SecondsToTime(chapter.Start)
	endH := utils.SecondsToTime(chapter.End)
	durationH := utils.SecondsToTime(chapter.Duration)
	p.chaptersTable.appendRow(number, startH, endH, durationH, chapter.Name)
	p.chaptersTable.ScrollToBeginning()
}

func (p *ChaptersPage) updateChapterEntry(row int, col int) {
	chapterNo, err := strconv.Atoi(p.chaptersTable.t.GetCell(row, 0).Text)
	if err != nil {
		// Part Number line found
		return
	}

	chapter, _ := p.ab.GetChapter(chapterNo)
	durationH := utils.SecondsToTime(chapter.Duration)
	d := newDialogWindow(p.mq, 11, 80, p.chaptersSection)
	f := newForm()
	f.SetTitle("Update Chapter Name:")
	f.AddTextView("Chapter #:  ", strconv.Itoa(chapter.Number), 5, 1, true, false)
	f.AddTextView("Duration:   ", strings.TrimLeft(durationH, " "), 10, 1, true, false)
	nameF := f.AddInputField("Chapter name:", chapter.Name, 60, nil, nil)
	f.AddButton("Save changes", func() {
		cell := p.chaptersTable.t.GetCell(row, col)
		cell.Text = nameF.GetText()
		chapter.Name = nameF.GetText()
		p.ab.SetChapter(chapterNo, *chapter)
		d.Close()
	})
	f.AddButton("Cancel", func() {
		d.Close()
	})
	d.setForm(f.f)
	d.Show()
}

func (p *ChaptersPage) searchReplaceChapters() {
	if p.searchChapters != "" {
		abCopy, err := p.ab.GetCopy()
		if err != nil {
			logger.Error("Can't create a copy of Audiobook struct: " + err.Error())
		} else {
			p.chaptersUndoStack.Push(abCopy)
			p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.SearchReplaceChaptersCommand{Audiobook: p.ab, SearchStr: p.searchChapters, ReplaceStr: p.replaceChapters}, false)
		}
	}
}

func (p *ChaptersPage) joinChapters() {
	abCopy, err := p.ab.GetCopy()
	if err != nil {
		logger.Error("Can't create a copy of Audiobook struct: " + err.Error())
	} else {
		p.chaptersUndoStack.Push(abCopy)
		p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.JoinChaptersCommand{Audiobook: p.ab}, false)
	}
}

func (p *ChaptersPage) undoChapters() {
	ab, err := p.chaptersUndoStack.Pop()
	if err == nil {
		p.ab.Parts = ab.Parts
		p.refreshChapters(p.ab)
	}
}

func (c *ChaptersPage) refreshChapters(ab *dto.Audiobook) {
	go c.displayParts(ab)
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
	p.mq.SendMessage(mq.ChaptersPage, mq.BuildController, &dto.BuildCommand{Audiobook: p.ab}, true)
	p.mq.SendMessage(mq.ChaptersPage, mq.Frame, &dto.SwitchToPageCommand{Name: "BuildPage"}, false)
}

// Simple Undo stack implementation
type UndoStack struct {
	stack *list.List
}

func NewUndoStack() *UndoStack {
	return &UndoStack{
		stack: list.New(),
	}
}

func (u *UndoStack) Push(obj *dto.Audiobook) {
	u.stack.PushBack(obj)
}

func (u *UndoStack) Pop() (*dto.Audiobook, error) {
	if u.stack.Len() == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	e := u.stack.Back()
	u.stack.Remove(e)
	return e.Value.(*dto.Audiobook), nil
}
