package ui

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"

	"abb_ia/internal/config"
	"abb_ia/internal/dto"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"

	"github.com/vpoluyaktov/tview"
)

type ChaptersPage struct {
	mq                       *mq.Dispatcher
	mainGrid                 *grid
	ab                       *dto.Audiobook
	inputAuthor              *tview.InputField
	inputTitle               *tview.InputField
	inputSeries              *tview.InputField
	inputSeriesNo            *tview.InputField
	inputGenre               *tview.DropDown
	inputNarrator            *tview.InputField
	inputCover               *tview.InputField
	buttonCreateBook         *tview.Button
	buttonCancel             *tview.Button
	textAreaDescription      *tview.TextArea
	inputSearchDescription   *tview.InputField
	inputReplaceDescription  *tview.InputField
	buttonDescriptionReplace *tview.Button
	buttonDescriptionUndo    *tview.Button
	chaptersSection          *grid
	chaptersTable            *table
	inputSearchChapters      *tview.InputField
	inputReplaceChapters     *tview.InputField
	inputPartSize            *tview.InputField
	buttonChaptersReplace    *tview.Button
	buttonChaptersUndo       *tview.Button
	buttonChaptersJoin       *tview.Button
	buttonRecalculateParts   *tview.Button
	searchDescription        string
	replaceDescription       string
	searchChapters           string
	replaceChapters          string
	partSize                 string
	chaptersUndoStack        *UndoStack
	descriptionUndoStack     *UndoStack
}

func newChaptersPage(dispatcher *mq.Dispatcher) *ChaptersPage {
	p := &ChaptersPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.ChaptersPage, p.dispatchMessage)

	p.mainGrid = newGrid()
	p.mainGrid.SetRows(9, -1, -1)
	p.mainGrid.SetColumns(0)

	// book info section
	infoSection := newGrid()
	infoSection.SetColumns(40, 30, -1, 30)
	infoSection.SetRows(5, 2)
	infoSection.SetBorder(true)
	infoSection.SetTitle(" Audiobook information: ")
	infoSection.SetTitleAlign(tview.AlignLeft)
	f0 := newForm()
	f0.SetBorderPadding(1, 0, 1, 2)
	p.inputAuthor = f0.AddInputField("Author:", "", 30, nil, func(s string) {
		if p.ab != nil {
			p.ab.Author = s
		}
	})
	p.inputTitle = f0.AddInputField("Title:", "", 30, nil, func(s string) {
		if p.ab != nil {
			p.ab.Title = s
		}
	})
	infoSection.AddItem(f0.Form, 0, 0, 1, 1, 0, 0, true)
	f1 := newForm()
	f1.SetBorderPadding(1, 0, 2, 2)
	p.inputSeries = f1.AddInputField("Series:", "", 20, nil, func(s string) {
		if p.ab != nil {
			p.ab.Series = s
		}
	})
	p.inputSeriesNo = f1.AddInputField("Series No:", "", 5, acceptInt, func(s string) {
		if p.ab != nil {
			p.ab.SeriesNo = s
		}
	})
	infoSection.AddItem(f1.Form, 0, 1, 1, 1, 0, 0, true)
	f2 := newForm()
	f2.SetBorderPadding(1, 0, 2, 2)
	p.inputGenre = f2.AddDropdown("Genre:", utils.AddSpaces(config.Instance().GetGenres()), 0, func(s string, i int) {
		if p.ab != nil {
			p.ab.Genre = strings.TrimSpace(s)
		}
	})
	p.inputNarrator = f2.AddInputField("Narrator:", "", 20, nil, func(s string) {
		if p.ab != nil {
			p.ab.Narrator = s
		}
	})
	infoSection.AddItem(f2.Form, 0, 2, 1, 1, 0, 0, true)
	f3 := newForm()
	f3.SetBorderPadding(0, 1, 1, 1)
	p.inputCover = f3.AddInputField("Book cover:", "", 0, nil, func(s string) {
		if p.ab != nil {
			p.ab.CoverURL = s
		}
	})
	infoSection.AddItem(f3.Form, 1, 0, 1, 4, 0, 0, true)
	f4 := newForm()
	f4.SetHorizontal(true)
	f4.SetButtonsAlign(tview.AlignRight)
	p.buttonCreateBook = f4.AddButton("Create Book", p.buildBook)
	p.buttonCancel = f4.AddButton("Cancel", p.stopConfirmation)
	infoSection.AddItem(f4.Form, 0, 3, 1, 1, 0, 0, false)
	p.mainGrid.AddItem(infoSection.Grid, 0, 0, 1, 1, 0, 0, false)

	// description section
	descriptionSection := newGrid()
	descriptionSection.SetColumns(-1, 33)
	descriptionSection.SetBorder(false)
	p.textAreaDescription = newTextArea("")
	p.textAreaDescription.SetChangedFunc(p.updateDescription)
	p.textAreaDescription.SetBorder(true)
	p.textAreaDescription.SetTitle(" Book description: ")
	p.textAreaDescription.SetTitleAlign(tview.AlignLeft)
	descriptionSection.AddItem(p.textAreaDescription, 0, 0, 1, 1, 0, 0, true)
	f5 := newForm()
	f5.SetBorder(true)
	f5.SetHorizontal(true)

	p.inputSearchDescription = f5.AddInputField("Search: ", "", 20, nil, func(s string) { p.searchDescription = s })
	p.inputReplaceDescription = f5.AddInputField("Replace:", "", 20, nil, func(s string) { p.replaceDescription = s })
	p.buttonDescriptionReplace = f5.AddButton("Replace", p.searchReplaceDescription)
	p.buttonDescriptionUndo = f5.AddButton(" Undo  ", p.undoDescription)
	f5.SetButtonsAlign(tview.AlignRight)
	descriptionSection.AddItem(f5.Form, 0, 1, 1, 1, 0, 0, true)
	p.mainGrid.AddItem(descriptionSection.Grid, 1, 0, 1, 1, 0, 0, true)

	// chapters section
	p.chaptersSection = newGrid()
	p.chaptersSection.SetColumns(-1, 33)
	p.chaptersTable = newTable()
	p.chaptersTable.SetBorder(true)
	p.chaptersTable.SetTitle(" Book chapters: ")
	p.chaptersTable.SetTitleAlign(tview.AlignLeft)
	p.chaptersTable.setHeaders("  # ", "Start", "End", "Duration", "Chapter name")
	p.chaptersTable.setWeights(1, 1, 1, 2, 10)
	p.chaptersTable.setAlign(tview.AlignRight, tview.AlignRight, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.chaptersTable.SetSelectedFunc(p.updateChapterEntry)
	p.chaptersTable.SetMouseDblClickFunc(p.updateChapterEntry)
	p.chaptersSection.AddItem(p.chaptersTable.Table, 0, 0, 1, 1, 0, 0, true)

	chaptersControls := newGrid()
	chaptersControls.SetColumns(-1)
	chaptersControls.SetRows(-2, -1)
	chaptersControls.SetBorder(true)

	f6 := newForm()
	f6.SetBorder(false)
	f6.SetHorizontal(true)
	p.inputSearchChapters = f6.AddInputField("Search: ", "", 20, nil, func(s string) { p.searchChapters = s })
	p.inputReplaceChapters = f6.AddInputField("Replace:", "", 20, nil, func(s string) { p.replaceChapters = s })
	p.buttonChaptersReplace = f6.AddButton("Replace", p.searchReplaceChapters)
	p.buttonChaptersUndo = f6.AddButton(" Undo  ", p.undoChapters)
	p.buttonChaptersJoin = f6.AddButton(" Join Similar Chapters ", p.joinChapters)
	f6.AddButton(" Use MP3 Names ", p.useMP3Names)
	f6.SetButtonsAlign(tview.AlignRight)
	f6.SetMouseDblClickFunc(func() {})
	chaptersControls.AddItem(f6.Form, 0, 0, 1, 1, 0, 0, false)

	f7 := newForm()
	f7.SetBorder(false)
	f7.SetHorizontal(true)
	p.inputPartSize = f7.AddInputField("Part size (Mb): ", "", 6, acceptInt, func(s string) { p.partSize = s })
	p.buttonRecalculateParts = f7.AddButton(" Recalculate Parts ", p.recalculateParts)
	f7.SetButtonsAlign(tview.AlignRight)
	chaptersControls.AddItem(f7.Form, 1, 0, 1, 1, 0, 0, false)

	p.chaptersSection.AddItem(chaptersControls.Grid, 0, 1, 1, 1, 0, 0, false)
	p.mainGrid.AddItem(p.chaptersSection.Grid, 2, 0, 1, 1, 0, 0, true)

	p.chaptersUndoStack = NewUndoStack()
	p.descriptionUndoStack = NewUndoStack()

	// screen elements navigation order
	p.mainGrid.SetNavigationOrder(
		p.inputAuthor,
		p.inputTitle,
		p.inputSeries,
		p.inputSeriesNo,
		p.inputGenre,
		p.inputNarrator,
		p.inputCover,
		p.buttonCreateBook,
		p.buttonCancel,
		p.textAreaDescription,
		p.inputSearchDescription,
		p.inputReplaceDescription,
		p.buttonDescriptionReplace,
		p.buttonDescriptionUndo,
		p.chaptersTable,
		p.inputSearchChapters,
		p.inputReplaceChapters,
		p.buttonChaptersReplace,
		p.buttonChaptersUndo,
		p.buttonChaptersJoin,
		p.inputPartSize,
		p.buttonRecalculateParts,
	)

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
	case *dto.RefreshDescriptionCommand:
		p.refreshDescription(dto.Audiobook)
	case *dto.RefreshChaptersCommand:
		p.refreshChapters(dto.Audiobook)
	default:
		m.UnsupportedTypeError(mq.ChaptersPage)
	}
}

func (p *ChaptersPage) displayBookInfo(ab *dto.Audiobook) {
	p.ab = ab
	p.inputAuthor.SetText(ab.Author)
	p.inputTitle.SetText(ab.Title)
	p.inputCover.SetText(ab.CoverURL)
	p.textAreaDescription.SetText(ab.Description, false)

	p.chaptersTable.Clear()
	p.chaptersTable.showHeader()
	p.chaptersTable.ScrollToBeginning()
	p.inputPartSize.SetText(strconv.Itoa(ab.Config.GetMaxFileSizeMb()))
	ui.SetFocus(p.chaptersSection.Grid)
}

func (p *ChaptersPage) displayParts(ab *dto.Audiobook) {
	p.chaptersTable.Clear()
	p.chaptersTable.showHeader()
	for _, part := range ab.Parts {
		p.addPart(&part)
	}
	ui.Draw()
}

func (p *ChaptersPage) addPart(part *dto.Part) {
	if len(p.ab.Parts) > 1 {
		number := strconv.Itoa(part.Number)
		p.chaptersTable.appendSeparator("", "", "", "", "Part # "+number+". Size: "+utils.BytesToHuman(part.Size))
	}
	for _, chapter := range part.Chapters {
		p.addChapter(&chapter)
	}
	ui.Draw()
}

func (p *ChaptersPage) addChapter(chapter *dto.Chapter) {
	number := strconv.Itoa(chapter.Number)
	startH := utils.SecondsToTime(chapter.Start)
	endH := utils.SecondsToTime(chapter.End)
	durationH := utils.SecondsToTime(chapter.Duration)
	p.chaptersTable.appendRow(number, startH, endH, durationH, chapter.Name)
	p.chaptersTable.ScrollToBeginning()
}

func (p *ChaptersPage) updateDescription() {
	if p.ab != nil {
		p.ab.Description = p.textAreaDescription.GetText()
	}
}

func (p *ChaptersPage) updateChapterEntry(row int, col int) {
	chapterNo, err := strconv.Atoi(p.chaptersTable.GetCell(row, 0).Text)
	if err != nil {
		// Part Number line found
		return
	}

	chapter, _ := p.ab.GetChapter(chapterNo)
	durationH := utils.SecondsToTime(chapter.Duration)
	d := newDialogWindow(p.mq, 11, 78, p.chaptersSection.Grid)
	f := newForm()
	f.SetTitle("Update Chapter Name:")
	f.AddTextView("Chapter #:  ", strconv.Itoa(chapter.Number), 5, 1, true, false)
	f.AddTextView("Duration:   ", strings.TrimLeft(durationH, " "), 10, 1, true, false)
	nameF := f.AddInputField("Chapter name:", chapter.Name, 60, nil, nil)
	f.AddButton("Save changes", func() {
		cell := p.chaptersTable.GetCell(row, col)
		cell.Text = nameF.GetText()
		chapter.Name = nameF.GetText()
		p.ab.SetChapter(chapterNo, *chapter)
		d.Close()
	})
	f.AddButton("Cancel", func() {
		d.Close()
	})
	d.setForm(f.Form)
	d.Show()
}

func (p *ChaptersPage) searchReplaceDescription() {
	if p.searchDescription != "" {
		abCopy, err := p.ab.GetCopy()
		if err != nil {
			logger.Error("Can't create a copy of Audiobook struct: " + err.Error())
		} else {
			p.descriptionUndoStack.Push(abCopy)
			p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.SearchReplaceDescriptionCommand{Audiobook: p.ab, SearchStr: p.searchDescription, ReplaceStr: p.replaceDescription}, true)
		}
	}
}

func (p *ChaptersPage) undoDescription() {
	ab, err := p.descriptionUndoStack.Pop()
	if err == nil {
		p.ab.Description = ab.Description
		p.textAreaDescription.SetText(p.ab.Description, false)
	}
	ui.Draw()
}

func (p *ChaptersPage) refreshDescription(ab *dto.Audiobook) {
	p.textAreaDescription.SetText(ab.Description, false)
	ui.Draw()
}

func (p *ChaptersPage) searchReplaceChapters() {
	if p.searchChapters != "" {
		abCopy, err := p.ab.GetCopy()
		if err != nil {
			logger.Error("Can't create a copy of Audiobook struct: " + err.Error())
		} else {
			p.chaptersUndoStack.Push(abCopy)
			p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.SearchReplaceChaptersCommand{Audiobook: p.ab, SearchStr: p.searchChapters, ReplaceStr: p.replaceChapters}, true)
		}
	}
}

func (p *ChaptersPage) joinChapters() {
	abCopy, err := p.ab.GetCopy()
	if err != nil {
		logger.Error("Can't create a copy of Audiobook struct: " + err.Error())
	} else {
		p.chaptersUndoStack.Push(abCopy)
		p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.JoinChaptersCommand{Audiobook: p.ab}, true)
	}
}

func (p *ChaptersPage) recalculateParts() {
	p.ab.Config.SetMaxFileSizeMb(utils.ToInt(p.partSize))
	abCopy, err := p.ab.GetCopy()
	if err != nil {
		logger.Error("Can't create a copy of Audiobook struct: " + err.Error())
	} else {
		p.chaptersUndoStack.Push(abCopy)
		p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.RecalculatePartsCommand{Audiobook: p.ab}, true)
	}
}

func (p *ChaptersPage) undoChapters() {
	ab, err := p.chaptersUndoStack.Pop()
	if err == nil {
		p.ab.Parts = ab.Parts
		p.refreshChapters(p.ab)
	}
}

func (p *ChaptersPage) refreshChapters(ab *dto.Audiobook) {
	go p.displayParts(ab)
}

func (p *ChaptersPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop editing chapters?", p.chaptersSection.Grid, p.stopChapters, func() {})
}

func (p *ChaptersPage) stopChapters() {
	p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.StopCommand{Process: "Chapters", Reason: "User request"}, false)
	p.mq.SendMessage(mq.ChaptersPage, mq.CleanupController, &dto.CleanupCommand{Audiobook: p.ab}, true)
	p.mq.SendMessage(mq.ChaptersPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

func (p *ChaptersPage) useMP3Names() {
	abCopy, err := p.ab.GetCopy()
	if err != nil {
		logger.Error("Can't create a copy of Audiobook struct: " + err.Error())
	} else {
		p.chaptersUndoStack.Push(abCopy)
		p.mq.SendMessage(mq.ChaptersPage, mq.ChaptersController, &dto.UseMP3NamesCommand{Audiobook: p.ab}, true)
	}
}

func (p *ChaptersPage) buildBook() {
	// update ab fields just to ensure (they are not updated automatically if a value wasn't change)
	p.ab.Author = p.inputAuthor.GetText()
	p.ab.Title = p.inputTitle.GetText()
	p.ab.Series = p.inputSeries.GetText()
	p.ab.SeriesNo = p.inputSeriesNo.GetText()
	p.ab.Narrator = p.inputNarrator.GetText()
	_, p.ab.Genre = p.inputGenre.GetCurrentOption()

	p.mq.SendMessage(mq.ChaptersPage, mq.BuildController, &dto.BuildCommand{Audiobook: p.ab}, true)
	p.mq.SendMessage(mq.ChaptersPage, mq.Frame, &dto.SwitchToPageCommand{Name: "BuildPage"}, true)
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
