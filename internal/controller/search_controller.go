package controller

import (
	"net/url"
	"sort"
	"strconv"
	"strings"

	"abb_ia/internal/config"
	"abb_ia/internal/dto"
	"abb_ia/internal/ia"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"
	"github.com/rivo/tview"
)

type SearchController struct {
	mq *mq.Dispatcher
}

func NewSearchController(dispatcher *mq.Dispatcher) *SearchController {
	c := &SearchController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.SearchController, c.dispatchMessage)
	return c
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
	logger.Info(mq.SearchController + " received " + cmd.String())
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.UpdateStatus{Message: "Fetching Internet Archive items..."}, false)
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)
	ia := ia_client.New(config.Instance().GetSearchRowsMax(), config.Instance().IsUseMock(), config.Instance().IsSaveMock())
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
		item.IaURL = ia_client.IA_BASE_URL + "/details/" + doc.Identifier
		item.LicenseUrl = doc.Licenseurl

		item.AudioFiles = make([]dto.AudioFile, 0)
		var totalSize int64 = 0
		var totalLength float64 = 0.0
		d := ia.GetItemDetails(doc.Identifier)
		if d != nil {
			item.Server = d.Server
			item.Dir = d.Dir
			if len(doc.Creator) > 0 && d.Metadata.Creator[0] != "" {
				item.Creator = doc.Creator[0]
			} else if len(d.Metadata.Creator) > 0 && d.Metadata.Creator[0] != "" {
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
				// collect mp3 files

				if utils.Contains(dto.Mp3Formats, format) {
					size, sErr := strconv.ParseInt(metadata.Size, 10, 64)
					length, lErr := utils.TimeToSeconds(metadata.Length)
					if sErr != nil || lErr != nil {
						logger.Error("Can't parse the file metadata: " + name)
					} else {
						file := dto.AudioFile{}
						file.Name = strings.TrimPrefix(name, "/")
						if metadata.Title != "" {
							file.Title = metadata.Title
						} else {
							file.Title = utils.SanitizeMp3FileName(file.Name)
						}
						file.Size = size
						file.Length = length
						file.Format = metadata.Format
						// check if there is a file with the same title but different bitrate. Keep highest bitrate only
						// see https://archive.org/details/voyage_moon_1512_librivox or https://archive.org/details/OTRR_Blair_of_the_Mounties_Singles for ex.
						addNewFile := true
						for i, oldFile := range item.AudioFiles {
							if file.Title == oldFile.Title {
								oldFilePriority := utils.GetIndex(dto.Mp3Formats, oldFile.Format)
								newFilePriority := utils.GetIndex(dto.Mp3Formats, file.Format)
								if newFilePriority > oldFilePriority {
									// remove old file from the list
									item.AudioFiles = append(item.AudioFiles[:i], item.AudioFiles[i+1:]...)
									totalSize -= oldFile.Size
									totalLength -= oldFile.Length
									// and add new one
									addNewFile = true
								} else if newFilePriority == oldFilePriority {
									// means multiple files have the same title
									addNewFile = true
								} else {
									addNewFile = false
								}
								break
							}
						}
						if addNewFile {
							item.AudioFiles = append(item.AudioFiles, file)
							totalSize += size
							totalLength += length
						}
					}
				}

				// collect image files
				if utils.Contains(dto.CoverFormats, format) {
					size, err := strconv.ParseInt(metadata.Size, 10, 64)
					if err == nil {
						file := dto.ImageFile{}
						file.Name = strings.TrimPrefix(name, "/")
						file.Size = size
						file.Format = metadata.Format
						item.ImageFiles = append(item.ImageFiles, file)
					}
				}
			}

			// sort mp3 files by name TODO: Check if sort is needed
			sort.Slice(item.AudioFiles, func(i, j int) bool { return item.AudioFiles[i].Name < item.AudioFiles[j].Name })
			item.TotalSize = totalSize
			item.TotalLength = totalLength

			// if len(d.Misc.Image) > 0 { // _thumb images are too small. Have to collect and sort my size all item images below
			// 	item.CoverUrl = d.Misc.Image
			// }
			// find biggest image by size (TODO: Need to find better solution. Maybe analyze if the image is colorful?)
			if len(item.ImageFiles) > 0 {
				biggestImage := item.ImageFiles[0]
				for i := 1; i < len(item.ImageFiles); i++ {
					if item.ImageFiles[i].Size > biggestImage.Size {
						biggestImage = item.ImageFiles[i]
					}
				}
				item.CoverUrl = (&url.URL{Scheme: "https", Host: item.Server, Path: item.Dir + "/" + biggestImage.Name}).String()
			} else {
				item.CoverUrl = "No cover available!"
			}

			if len(item.AudioFiles) > 0 {
				itemsFetched++
				sp := &dto.SearchProgress{ItemsTotal: itemsTotal, ItemsFetched: itemsFetched}
				c.mq.SendMessage(mq.SearchController, mq.SearchPage, sp, false)
				c.mq.SendMessage(mq.SearchController, mq.SearchPage, item, false)
			}
		}
		logger.Debug(mq.SearchController + " fetched first " + strconv.Itoa(itemsFetched) + " items from " + strconv.Itoa(itemsTotal) + " total")
	}
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)

	if itemsFetched == 0 {
		c.mq.SendMessage(mq.SearchController, mq.SearchPage, &dto.NothingFoundError{SearchCondition: cmd.SearchCondition}, false)
	}
}