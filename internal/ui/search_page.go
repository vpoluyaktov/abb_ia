package ui

import (
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

type SearchPage struct {
	mq              *mq.Dispatcher
	mainGrid        *grid
	searchCondition dto.SearchCondition
	isSearchRunning bool
	searchResult    []*dto.IAItem

	searchSection         *grid
	author                *tview.InputField
	title                 *tview.InputField
	SortBy                *tview.DropDown
	sortOrder             *tview.DropDown
	searchButton          *tview.Button
	clearButton           *tview.Button
	createAudioBookButton *tview.Button
	SettingsButton        *tview.Button

	resultSection *grid
	resultTable   *table
	urlSection    *grid
	urlField      *tview.TextView

	detailsSection  *grid
	descriptionView *tview.TextView
	filesTable      *table
}

func newSearchPage(dispatcher *mq.Dispatcher) *SearchPage {
	p := &SearchPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.SearchPage, p.dispatchMessage)

	p.searchCondition.Author = config.Instance().GetDefaultAuthor()
	p.searchCondition.Title = config.Instance().GetDefaultTitle()

	p.mainGrid = newGrid()
	p.mainGrid.SetRows(9, -1, -1, 3)
	p.mainGrid.SetColumns(0)

	// search section
	p.searchSection = newGrid()
	p.searchSection.SetColumns(50, -1, -1)
	p.searchSection.SetBorder(true)
	p.searchSection.SetTitle(" Internet Archive Search ")
	p.searchSection.SetTitleAlign(tview.AlignLeft)
	f := newForm()
	f.SetHorizontal(false)
	p.author = f.AddInputField("Creator", config.Instance().GetDefaultAuthor(), 40, nil, func(t string) { p.searchCondition.Author = t })
	p.title = f.AddInputField("Title", config.Instance().GetDefaultTitle(), 40, nil, func(t string) { p.searchCondition.Title = t })

	p.searchButton = f.AddButton("Search", p.runSearch)
	p.clearButton = f.AddButton("Clear", p.clearEverything)
	f.SetButtonsAlign(tview.AlignRight)
	p.searchSection.AddItem(f, 0, 0, 1, 1, 0, 0, true)
	f = newForm()
	p.SortBy = f.AddDropdown("Sort by:", utils.AddSpaces(config.Instance().GetSortByOptions()), utils.GetIndex(config.Instance().GetSortByOptions(), config.Instance().GetSortBy()), func(o string, i int) { p.searchCondition.SortBy = p.mapSortBy(o) })
	p.sortOrder = f.AddDropdown("Sort order:", utils.AddSpaces(config.Instance().GetSortOrderOptions()), utils.GetIndex(config.Instance().GetSortOrderOptions(), config.Instance().GetSortOrder()), func(o string, i int) { p.searchCondition.SortOrder = p.mapSortOrder(o) })
	p.searchSection.AddItem(f, 0, 1, 1, 1, 0, 0, true)
	g := newGrid()
	g.SetRows(-1, -1)
	g.SetColumns(0)
	f = newForm()
	f.SetHorizontal(false)
	f.SetButtonsAlign(tview.AlignRight)
	p.createAudioBookButton = f.AddButton("Create Audiobook", p.createBook)
	g.AddItem(f, 0, 0, 1, 1, 1, 1, true)
	f = newForm()
	f.SetHorizontal(false)
	f.SetButtonsAlign(tview.AlignRight)
	g.AddItem(f, 1, 0, 1, 1, 1, 1, true)
	p.SettingsButton = f.AddButton("Settings", p.updateConfig)
	p.searchSection.AddItem(g, 0, 2, 1, 1, 0, 0, true)

	p.mainGrid.AddItem(p.searchSection.Grid, 0, 0, 1, 1, 0, 0, true)

	// result section
	p.resultSection = newGrid()
	p.resultSection.SetColumns(-1)
	p.resultSection.SetTitle(" Search result: ")
	p.resultSection.SetTitleAlign(tview.AlignLeft)
	p.resultSection.SetBorder(true)

	p.resultTable = newTable()
	p.resultTable.setHeaders(" # ", "Creator", "Title", "Files", "Duration", "Size")
	p.resultTable.setWeights(1, 3, 6, 1, 2, 2)
	p.resultTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignRight)
	p.resultTable.SetSelectionChangedFunc(p.updateDetails)
	p.resultTable.SetSelectedFunc(p.itemSelected)
	p.resultTable.SetMouseDblClickFunc(p.itemSelected)
	p.resultTable.setLastRowEvent(p.lastRowEvent)
	p.resultSection.AddItem(p.resultTable.Table, 0, 0, 1, 1, 0, 0, true)
	p.mainGrid.AddItem(p.resultSection.Grid, 1, 0, 1, 1, 0, 0, true)

	// details section
	p.detailsSection = newGrid()
	p.detailsSection.SetRows(-1)
	p.detailsSection.SetColumns(-1, 1, -1)

	p.descriptionView = tview.NewTextView()
	p.descriptionView.SetWrap(true)
	p.descriptionView.SetWordWrap(true)
	p.descriptionView.SetBorder(true)
	p.descriptionView.SetTitle(" Description: ")
	p.descriptionView.SetTitleAlign(tview.AlignLeft)
	p.detailsSection.AddItem(p.descriptionView, 0, 0, 1, 1, 0, 0, true)

	p.filesTable = newTable()
	p.filesTable.SetBorder(true)
	p.filesTable.SetTitle(" Files: ")
	p.filesTable.SetTitleAlign(tview.AlignLeft)
	p.filesTable.setHeaders("File name", "Format", "Duration", "Size")
	p.filesTable.setWeights(3, 2, 2, 2)
	p.filesTable.setAlign(tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignRight)
	p.detailsSection.AddItem(p.filesTable.Table, 0, 2, 1, 1, 0, 0, true)

	p.mainGrid.AddItem(p.detailsSection.Grid, 2, 0, 1, 1, 0, 0, true)

	// URL section
	p.urlSection = newGrid()
	p.urlSection.SetRows(1)
	p.urlSection.SetColumns(-1)
	p.urlSection.SetBorder(true)
	p.urlSection.SetTitle(" Internet Archive item url: ")
	p.urlSection.SetTitleAlign(tview.AlignLeft)
	p.urlField = tview.NewTextView()
	p.urlSection.AddItem(p.urlField, 0, 0, 1, 1, 0, 0, false)
	p.mainGrid.AddItem(p.urlSection.Grid, 3, 0, 1, 1, 0, 0, true)

	p.mq.SendMessage(mq.SearchPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
	ui.SetFocus(p.searchSection.Grid)

	p.mainGrid.Focus(func(pr tview.Primitive) {
		if p.resultTable.GetRowCount() == 0 {
			ui.SetFocus(p.searchSection.Grid)
		} else {
			ui.SetFocus(p.resultSection.Grid)
		}
	})

	// screen navigation order
	p.mainGrid.SetNavigationOrder(
		p.author,
		p.title,
		p.searchButton,
		p.clearButton,
		p.SortBy,
		p.sortOrder,
		p.createAudioBookButton,
		p.SettingsButton,
		p.resultTable.Table,
		p.descriptionView,
		p.filesTable.Table,
	)

	return p
}

func (p *SearchPage) checkMQ() {
	m := p.mq.GetMessage(mq.SearchPage)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *SearchPage) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.UpdateSearchConfigCommand:
		p.updateSearchConfig(dto)
	case *dto.IAItem:
		p.updateResult(dto)
	case *dto.SearchProgress:
		p.updateTitle(dto)
	case *dto.SearchComplete:
		p.isSearchRunning = false
	case *dto.NothingFoundError:
		p.showNothingFoundError(dto)
	case *dto.LastPageMessage:
		p.showLastPageMessage(dto)
	case *dto.NewAppVersionFound:
		p.showNewVersionMessage(dto)
	case *dto.FFMPEGNotFoundError:
		p.showFFMPEGNotFoundError(dto)
	default:
		m.UnsupportedTypeError(mq.SearchPage)
	}
}

func (p *SearchPage) runSearch() {
	if p.isSearchRunning {
		return
	}
	p.isSearchRunning = true
	p.clearSearchResults()
	p.resultTable.showHeader()
	p.mq.SendMessage(mq.SearchPage, mq.SearchController, &dto.SearchCommand{Condition: p.searchCondition}, false)
	p.mq.SendMessage(mq.SearchPage, mq.TUI, &dto.SetFocusCommand{Primitive: p.resultTable.Table}, true)
}

func (p *SearchPage) clearSearchResults() {
	p.searchResult = make([]*dto.IAItem, 0)
	p.resultSection.SetTitle(" Search result: ")
	p.resultTable.Clear()
	p.descriptionView.SetText("")
	p.filesTable.Clear()
}

func (p *SearchPage) clearEverything() {
	p.author.SetText("")
	p.title.SetText("")
	p.clearSearchResults()
	p.urlField.SetText("")
}

func (p *SearchPage) updateSearchConfig(c *dto.UpdateSearchConfigCommand) {
	if p.author.GetText() == "" && p.title.GetText() == "" {
		p.author.SetText(c.Config.GetDefaultAuthor())
		p.title.SetText(c.Config.GetDefaultTitle())
		p.SortBy.SetCurrentOption(utils.GetIndex(c.Config.GetSortByOptions(), c.Config.GetSortBy()))
		p.sortOrder.SetCurrentOption(utils.GetIndex(c.Config.GetSortOrderOptions(), c.Config.GetSortOrder()))
		ui.Draw()
	}
}

func (p *SearchPage) lastRowEvent() {
	if p.isSearchRunning {
		return
	}
	p.isSearchRunning = true
	p.mq.SendMessage(mq.SearchPage, mq.SearchController, &dto.GetNextPageCommand{Condition: p.searchCondition}, false)
}

func (p *SearchPage) updateResult(i *dto.IAItem) {
	logger.Debug(mq.SearchPage + ": Got AI Item: " + i.Title)
	p.searchResult = append(p.searchResult, i)
	row, col := p.resultTable.GetSelection()
	p.resultTable.appendRow(strconv.Itoa(p.resultTable.GetRowCount()), i.Creator, i.Title, strconv.Itoa(len(i.AudioFiles)), utils.SecondsToTime(i.TotalLength), utils.BytesToHuman(i.TotalSize))
	p.resultTable.Select(row, col)
	ui.Draw()
}

func (p *SearchPage) updateTitle(sp *dto.SearchProgress) {
	p.resultSection.SetTitle(fmt.Sprintf(" Search result (fetched %d items from %d total): ", sp.ItemsFetched, sp.ItemsTotal))
}

func (p *SearchPage) updateDetails(row int, col int) {
	if row > 0 && len(p.searchResult) > 0 && len(p.searchResult) >= row {
		item := p.searchResult[row-1]
		d := item.Description
		p.descriptionView.SetText(d)
		p.descriptionView.ScrollToBeginning()

		p.filesTable.Clear()
		p.filesTable.showHeader()
		files := item.AudioFiles
		for _, f := range files {
			p.filesTable.appendRow(f.Name, f.Format, utils.SecondsToTime(f.Length), utils.BytesToHuman(f.Size))
		}
		p.filesTable.ScrollToBeginning()

		p.urlField.SetText(" " + item.IaURL)
	}
}

// process Enter and DoubleClick on the result table
func (p *SearchPage) itemSelected(row int, col int) {
	if row > 0 && len(p.searchResult) > 0 && row <= len(p.searchResult) {
		p.createBook()
	}
}

func (p *SearchPage) createBook() {
	// get selectet row from the results table
	row, _ := p.resultTable.GetSelection()
	if row <= 0 || len(p.searchResult) <= 0 || len(p.searchResult) < row {
		newMessageDialog(p.mq, "Error", "\nPlease conduct a search beforehand.", p.searchSection.Grid, func() {})
	} else if !(utils.CommandExists("ffmpeg") && utils.CommandExists("ffprobe")) {
		p.showFFMPEGNotFoundError(&dto.FFMPEGNotFoundError{})
	} else {
		item := p.searchResult[row-1]
		// create new audiobook object
		ab := &dto.Audiobook{}
		ab.IAItem = item
		c := config.Instance().GetCopy()
		ab.Config = &c

		d := newDialogWindow(p.mq, 17, 55, p.resultSection.Grid)
		f := newForm()
		f.SetTitle("Create Audiobook")
		f.AddInputField("Concurrent Downloaders:", utils.ToString(ab.Config.GetConcurrentDownloaders()), 8, acceptInt, func(t string) { ab.Config.SetConcurrentDownloaders(utils.ToInt(t)) })
		f.AddInputField("Concurrent Encoders:", utils.ToString(ab.Config.GetConcurrentEncoders()), 8, acceptInt, func(t string) { ab.Config.SetConcurrentEncoders(utils.ToInt(t)) })
		f.AddCheckbox("Re-encode .mp3 files to the same Bit Rate?", ab.Config.IsReEncodeFiles(), func(t bool) { ab.Config.SetReEncodeFiles(t) })
		f.AddInputField("Bit Rate (Kbps):", utils.ToString(ab.Config.GetBitRate()), 8, acceptInt, func(t string) { ab.Config.SetBitRate(utils.ToInt(t)) })
		f.AddInputField("Sample Rate (Hz):", utils.ToString(ab.Config.GetSampleRate()), 8, acceptInt, func(t string) { ab.Config.SetSampleRate(utils.ToInt(t)) })
		f.AddInputField("Audiobook part max file size (Mb):", utils.ToString(ab.Config.GetMaxFileSizeMb()), 8, acceptInt, func(t string) { ab.Config.SetMaxFileSizeMb(utils.ToInt(t)) })

		f.AddButton("Create Audiobook", func() {
			p.startDownload(ab)
			d.Close()
		})
		f.AddButton("Cancel", func() {
			d.Close()
		})
		d.setForm(f.Form)
		d.Show()
	}
}

func (p *SearchPage) startDownload(ab *dto.Audiobook) {
	p.mq.SendMessage(mq.SearchPage, mq.DownloadController, &dto.DownloadCommand{Audiobook: ab}, true)
	p.mq.SendMessage(mq.SearchPage, mq.Frame, &dto.SwitchToPageCommand{Name: "DownloadPage"}, false)
}

func (p *SearchPage) updateConfig() {
	p.mq.SendMessage(mq.SearchPage, mq.ConfigPage, &dto.DisplayConfigCommand{Config: config.Instance().GetCopy()}, true)
	p.mq.SendMessage(mq.SearchPage, mq.Frame, &dto.SwitchToPageCommand{Name: "ConfigPage"}, false)
}

func (p *SearchPage) showNothingFoundError(dto *dto.NothingFoundError) {
	newMessageDialog(p.mq, "Error",
		"\nNo results were found for your search term:\n"+
			"Creator: [darkblue]'"+dto.Condition.Author+"'[black] Title: [darkblue]'"+dto.Condition.Title+"'[black].\n"+
			"Please revise your search criteria.",
		p.searchSection.Grid, func() {})
}

func (p *SearchPage) showLastPageMessage(dto *dto.LastPageMessage) {
	newMessageDialog(p.mq, "Notification",
		"No more items were found for your search term: \n"+
			"Creator: [darkblue]'"+dto.Condition.Author+"'[black] Title: [darkblue]'"+dto.Condition.Title+"'[black].\n"+
			"This is the last page.\n",
		p.resultSection.Grid, func() {})
}

func (p *SearchPage) showFFMPEGNotFoundError(dto *dto.FFMPEGNotFoundError) {
	newMessageDialog(p.mq, "Error",
		"This application requires the utilities [darkblue]ffmpeg[black] and [darkblue]ffprobe[black].\n"+
			"Please install both [darkblue]ffmpeg[black] and [darkblue]ffprobe[black] by following the instructions provided on FFMPEG website\n"+
			"[darkblue]https://ffmpeg.org/download.html",
		p.searchSection.Grid, func() {})
}

func (p *SearchPage) showNewVersionMessage(dto *dto.NewAppVersionFound) {
	newMessageDialog(p.mq, "Notification",
		"New version of the Audiobook Builder has been released: [darkblue]"+dto.NewVersion+"[black]\n"+
			"Your current version is [darkblue]v"+dto.CurrentVersion+"[black]\n"+
			"You can download the new version of the application from:\n[darkblue]https://github.com/"+config.Instance().GetRepoOwner()+"/"+config.Instance().GetRepoName()+"/releases",
		p.searchSection.Grid, func() {})
}

func (p *SearchPage) mapSortBy(source string) string {
	switch s := strings.TrimSpace(source); s {
	case "Creator":
		return "creator"
	case "Title":
		return "title"
	case "Date":
		return "date"
	case "Size":
		return "item_size"
	default:
		return source
	}
}

func (p *SearchPage) mapSortOrder(source string) string {
	switch s := strings.TrimSpace(source); s {
	case "Ascending":
		return "asc"
	case "Descending":
		return "desc"
	default:
		return source
	}
}
