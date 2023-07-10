package controller

import (
	"sort"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/ia_client"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
	"github.com/vpoluyaktov/abb_ia/internal/utils"
)

type SearchController struct {
	mq *mq.Dispatcher
}

func NewSearchController(dispatcher *mq.Dispatcher) *SearchController {
	sc := &SearchController{}
	sc.mq = dispatcher
	sc.mq.RegisterListener(mq.SearchController, sc.dispatchMessage)
	return sc
}

func (c *SearchController) checkMQ() {
	m := c.mq.GetMessage(mq.SearchController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *SearchController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.SearchCommand:
		go c.performSearch(dto)
	default:
		m.UnsupportedTypeError(mq.SearchController)
	}
}

func (c *SearchController) performSearch(cmd *dto.SearchCommand) {
	logger.Debug(mq.SearchController + ": Received SearchCommand with condition: " + cmd.SearchCondition)
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.UpdateStatus{Message: "Fetching Internet Archive items..."}, false)
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)
	ia := ia_client.New(config.IsUseMock(), config.IsSaveMock())
	resp := ia.Search(cmd.SearchCondition, "audio")
	if resp == nil {
		logger.Error(mq.SearchController + ": Failed to perform IA search with condition: " + cmd.SearchCondition)
	}

	itemsTotal := resp.Response.NumFound
	itemsFetched := 0

	docs := resp.Response.Docs
	for _, doc := range docs {
		item := &dto.IAItem{}
		item.ID = doc.Identifier
		item.Title = tview.Escape(doc.Title)

		// collect mp3 files
		item.FilesCount = 0
		item.Files = make([]dto.File, 0)
		var totalSize int64 = 0
		var totalLength float64 = 0.0
		d := ia.GetItemDetails(doc.Identifier)
		if d != nil {
			item.Server = d.Server
			item.Dir = d.Dir
			if len(d.Metadata.Creator) > 0 && d.Metadata.Creator[0] != "" {
				item.Creator = d.Metadata.Creator[0]
			} else if len(d.Metadata.Artist) > 0 && d.Metadata.Artist[0] != "" {
				item.Creator = d.Metadata.Artist[0]
			} else {
				item.Creator = "Internet Archive"
			}
			if len(d.Metadata.Description) > 0 {
				item.Description = tview.Escape(ia.Html2Text(d.Metadata.Description[0]))
			}
			for name, metadata := range d.Files {
				format := metadata.Format
				if utils.Contains(dto.FormatList, format) {
					size, sErr := strconv.ParseInt(metadata.Size, 10, 64)
					length, lErr := utils.TimeToSeconds(metadata.Length)
					if sErr == nil && lErr == nil {
						file := dto.File{}
						file.Name = strings.TrimPrefix(name, "/")
						file.Size = size
						file.SizeH, _ = utils.BytesToHuman(size)
						file.Length = length
						file.LengthH, _ = utils.SecondsToTime(length)
						file.Format = metadata.Format
						totalSize += size
						totalLength += length
						item.Files = append(item.Files, file)
					}
				}
			}
			// sort files by name
			sort.Slice(item.Files, func(i, j int) bool { return item.Files[i].Name < item.Files[j].Name })
			item.TotalSize = totalSize
			item.TotalSizeH, _ = utils.BytesToHuman(totalSize)
			item.TotalLength = totalLength
			item.TotalLengthH, _ = utils.SecondsToTime(totalLength)
			item.FilesCount = len(item.Files)
		}
		if item.FilesCount > 0 {
			itemsFetched++
			sp := &dto.SearchProgress{ItemsTotal: itemsTotal, ItemsFetched: itemsFetched}
			c.mq.SendMessage(mq.SearchController, mq.SearchPage, sp, true)
			c.mq.SendMessage(mq.SearchController, mq.SearchPage, item, true)
		}
	}
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
}
