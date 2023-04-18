package controller

import (
	"sort"
	"strconv"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/utils"
	"github.com/vpoluyaktov/audiobook_creator_IA/pkg/ia_client"
)

const (
	controllerName  = "SearchController"
	uiComponentName = "SearchPanel"
)

type SearchController struct {
	dispatcher *mq.Dispatcher
}

func NewSearchController(dispatcher *mq.Dispatcher) *SearchController {
	sp := &SearchController{}
	sp.dispatcher = dispatcher
	sp.dispatcher.RegisterListener(controllerName, sp.dispatchMessage)
	return sp
}

func (p *SearchController) checkMQ() {
	m := p.dispatcher.GetMessage(controllerName)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *SearchController) sendMessage(from string, to string, dtoType string, dto dto.Dto, async bool) {
	m := &mq.Message{}
	m.From = from
	m.To = to
	m.Type = dtoType
	m.Dto = dto
	m.Async = async
	p.dispatcher.SendMessage(m)
}

func (p *SearchController) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.SearchCommandType:
		if c, ok := m.Dto.(*dto.SearchCommand); ok {
			go p.performSearch(c)
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}

func (p *SearchController) performSearch(c *dto.SearchCommand) {
	logger.Debug(controllerName + ": Received SearchCommand with condition: " + c.SearchCondition)
	ia := ia_client.New()
	resp := ia.Search(c.SearchCondition, "audio")
	if resp == nil {
		logger.Error(controllerName + ": Failed to perform IA search with condition: " + c.SearchCondition)
	}

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
						file.Name = tview.Escape(name)
						file.Size = size
						file.SizeH, _ = utils.BytesToHuman(size)
						file.Length = length
						file.LengthH, _ = utils.SecondToTime(length)
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
			item.TotalLengthH, _ = utils.SecondToTime(totalLength)
			item.FilesCount = len(item.Files)
		}
		if item.FilesCount > 0 {
			p.sendMessage(controllerName, uiComponentName, dto.IAItemType, item, true)
		}
	}
}
