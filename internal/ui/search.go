package ui

import (
	"strconv"

	"code.rocketnine.space/tslocum/cview"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

const (
	uiComponentName = "SearchPanel"
	controllerName  = "SearchController"
)

type searchPanel struct {
	grid                *cview.Grid
	dispatcher          *mq.Dispatcher
	searchSection       *cview.Flex
	searchCriteria      string
	searchResult        []*dto.IAItem
	searchResultSection *cview.Grid
	searchResultTable   *table
	detailsSection      *cview.Grid
	descriptionView     *cview.TextView
	filesTable          *table
}

func newSearchPanel(dispatcher *mq.Dispatcher) *searchPanel {
	p := &searchPanel{}
	p.dispatcher = dispatcher
	p.dispatcher.RegisterListener(uiComponentName, p.dispatchMessage)

	p.grid = cview.NewGrid()
	p.grid.SetRows(5, -1, -1)
	p.grid.SetColumns(0)

	// search section
	p.searchSection = cview.NewFlex()
	p.searchSection.SetDirection(cview.FlexRow)
	p.searchSection.SetBorder(true)
	p.searchSection.SetTitle(" Internet Archive Search ")
	p.searchSection.SetTitleAlign(cview.AlignLeft)
	form := cview.NewForm()
	form.SetHorizontal(true)
	form.AddInputField("Search criteria", "", 40, nil, func(t string) { p.searchCriteria = t })
	form.AddButton("Search", p.runSearch)
	form.AddButton("Clear", p.clearSearchResults)
	p.searchSection.AddItem(form, 0, 1, true)
	p.grid.AddItem(p.searchSection, 0, 0, 1, 1, 0, 0, true)

	// result section
	p.searchResultSection = cview.NewGrid()
	p.searchResultSection.SetColumns(-1)
	p.searchResultSection.SetTitle(" Search result ")
	p.searchResultSection.SetTitleAlign(cview.AlignLeft)
	p.searchResultSection.SetBorder(true)

	p.searchResultTable = newTable()
	p.searchResultTable.setHeaders([]string{"Title", "Files", "Duration (HH:MM:SS)", "Size"})
	p.searchResultTable.setWidths([]int{4, 1, 1, 1})
	p.searchResultTable.setAlign([]uint{cview.AlignLeft, cview.AlignRight, cview.AlignRight, cview.AlignRight})
	p.searchResultTable.t.SetSelectionChangedFunc(p.updateDetails)
	p.searchResultSection.AddItem(p.searchResultTable.t, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(p.searchResultSection, 1, 0, 1, 1, 0, 0, true)

	// details section
	p.detailsSection = cview.NewGrid()
	p.detailsSection.SetRows(-1)
	p.detailsSection.SetColumns(-1, 1, -1)
	// p.detailsSection.SetTitle(" Details ")
	// p.detailsSection.SetTitleAlign(cview.AlignLeft)
	// p.detailsSection.SetBorder(true)

	p.descriptionView = cview.NewTextView()
	p.descriptionView.SetDynamicColors(true)
	p.descriptionView.SetRegions(true)
	p.descriptionView.SetWrap(true)
	p.descriptionView.SetWordWrap(true)
	p.descriptionView.SetBorder(true)
	p.descriptionView.SetTitle(" Description: ")
	p.descriptionView.SetTitleAlign(cview.AlignLeft)
	p.detailsSection.AddItem(p.descriptionView, 0, 0, 1, 1, 0, 0, true)

	p.filesTable = newTable()
	p.filesTable.t.SetBorder(true)
	p.filesTable.t.SetTitle(" Files: ")
	p.filesTable.t.SetTitleAlign(cview.AlignLeft)
	p.filesTable.setHeaders([]string{"File name", "Format", "Duration", "Size"})
	p.filesTable.setWidths([]int{4, 1, 1, 1})
	p.filesTable.setAlign([]uint{cview.AlignLeft, cview.AlignRight, cview.AlignRight, cview.AlignRight})
	p.detailsSection.AddItem(p.filesTable.t, 0, 2, 1, 1, 0, 0, true)

	p.grid.AddItem(p.detailsSection, 2, 0, 1, 1, 0, 0, true)

	return p
}

func (p *searchPanel) readMessages() {
	m := p.dispatcher.GetMessage(uiComponentName)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *searchPanel) sendMessage(from string, to string, dtoType string, dto dto.Dto, async bool) {
	m := &mq.Message{}
	m.From = from
	m.To = to
	m.Type = dtoType
	m.Dto = dto
	m.Async = async
	p.dispatcher.SendMessage(m)
}

func (p *searchPanel) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.IAItemType:
		if r, ok := m.Dto.(*dto.IAItem); ok {
			go p.updateSearchResult(r)
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}

func (p *searchPanel) runSearch() {
	p.clearSearchResults()
	p.searchResultTable.showHeader()
	// Disable Search Button here
	p.sendMessage(uiComponentName, controllerName, dto.SearchCommandType, dto.SearchCommand{SearchCondition: p.searchCriteria}, true)
}

func (p *searchPanel) clearSearchResults() {
	p.searchResult = make([]*dto.IAItem, 0)
	p.searchResultTable.clear()
	p.descriptionView.SetText("")
	p.filesTable.clear()
}

func (p *searchPanel) updateSearchResult(i *dto.IAItem) {
	logger.Debug(uiComponentName + ": Got AI Item: " + i.Title)
	p.searchResult = append(p.searchResult, i)
	p.searchResultTable.appendRow([]string{i.Title, strconv.Itoa(i.FilesCount), i.TotalLengthH, i.TotalSizeH})
	p.sendMessage(uiComponentName, "TUI", dto.CommandType, dto.Command{Command: "RedrawUI"}, true)
}

func (p *searchPanel) updateDetails(row int, col int) {
	if row > 0 && len(p.searchResult) > 0 && len(p.searchResult) >= row {
		d := p.searchResult[row-1].Description
		p.descriptionView.SetText(d)
		p.descriptionView.ScrollToBeginning()

		p.filesTable.clear()
		p.filesTable.showHeader()
		files := p.searchResult[row-1].Files
		for _, f := range files {
			p.filesTable.appendRow([]string{f.Name, f.Format, f.LengthH, f.SizeH})
		}
		p.filesTable.t.Select(1, 0)
	}

}
