package controller

import (
	"fmt"
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

	"github.com/vpoluyaktov/tview"
)

var (
	// mp3 format list ranged by priority
	Mp3Formats = []string{"16Kbps MP3", "24Kbps MP3", "32Kbps MP3", "40Kbps MP3", "48Kbps MP3", "56Kbps MP3", "64Kbps MP3", "80Kbps MP3", "96Kbps MP3", "112Kbps MP3", "128Kbps MP3", "144Kbps MP3", "160Kbps MP3", "224Kbps MP3", "256Kbps MP3", "320Kbps MP3", "VBR MP3"}
	// audiobook cover formats
	CoverFormats = []string{"JPEG", "JPEG Thumb"}
)

type SearchController struct {
	mq                *mq.Dispatcher
	ia                *ia_client.IAClient
	totalItemsFetched int
}

func NewSearchController(dispatcher *mq.Dispatcher) *SearchController {
	c := &SearchController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.SearchController, c.dispatchMessage)
	return c
}

func (c *SearchController) checkMQ() {
	m, err := c.mq.GetMessage(mq.SearchController)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get message for SearchController: %v", err))
		return
	}
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *SearchController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.SearchCommand:
		go c.search(dto)
	case *dto.GetNextPageCommand:
		go c.getGetNextPage(dto)
	default:
		m.UnsupportedTypeError(mq.SearchController)
	}
}

func (c *SearchController) search(cmd *dto.SearchCommand) {
	logger.Info(fmt.Sprintf("Searching for: %s - %s", cmd.Condition.Author, cmd.Condition.Title))
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.UpdateStatus{Message: "Fetching Internet Archive items..."}, mq.PriorityNormal)
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, mq.PriorityNormal)
	c.totalItemsFetched = 0
	c.ia = ia_client.New(config.Instance().GetRowsPerPage(), config.Instance().IsUseMock(), config.Instance().IsSaveMock())
	resp := c.ia.Search(cmd.Condition.Author, cmd.Condition.Title, "audio", cmd.Condition.SortBy, cmd.Condition.SortOrder)
	if resp == nil {
		logger.Error(mq.SearchController + ": Failed to perform IA search with condition: " + cmd.Condition.Author + " - " + cmd.Condition.Title)
	}
	itemsFetched, err := c.fetchDetails(resp)
	if err != nil {
		logger.Error(mq.SearchController + ": Failed to fetch item details: " + err.Error())
	}
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, mq.PriorityNormal)
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.UpdateStatus{Message: ""}, mq.PriorityNormal)
	c.mq.SendMessage(mq.SearchController, mq.SearchPage, &dto.SearchComplete{Condition: cmd.Condition}, mq.PriorityNormal)
	if itemsFetched == 0 {
		logger.Info("Nothing found")
		c.mq.SendMessage(mq.SearchController, mq.SearchPage, &dto.NothingFoundError{Condition: cmd.Condition}, mq.PriorityNormal)
	} else {
		logger.Info(fmt.Sprintf("Items fetched: %d", c.totalItemsFetched))
	}
}

func (c *SearchController) getGetNextPage(cmd *dto.GetNextPageCommand) {
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.UpdateStatus{Message: "Fetching Internet Archive items..."}, mq.PriorityNormal)
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, mq.PriorityNormal)
	resp := c.ia.GetNextPage(cmd.Condition.Author, cmd.Condition.Title, "audio", cmd.Condition.SortBy, cmd.Condition.SortOrder)
	if resp == nil {
		logger.Error(mq.SearchController + ": Failed to perform IA search with condition: " + cmd.Condition.Author + " - " + cmd.Condition.Title)
	}
	itemsFetched, err := c.fetchDetails(resp)
	if err != nil {
		logger.Error(mq.SearchController + ": Failed to fetch item details: " + err.Error())
	}
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, mq.PriorityNormal)
	c.mq.SendMessage(mq.SearchController, mq.Footer, &dto.UpdateStatus{Message: ""}, mq.PriorityNormal)
	c.mq.SendMessage(mq.SearchController, mq.SearchPage, &dto.SearchComplete{Condition: cmd.Condition}, mq.PriorityNormal)
	if itemsFetched == 0 {
		logger.Info("Last page reached")
		c.mq.SendMessage(mq.SearchController, mq.SearchPage, &dto.LastPageMessage{Condition: cmd.Condition}, mq.PriorityNormal)
	} else {
		logger.Info(fmt.Sprintf("Items fetched: %d", c.totalItemsFetched))
	}
}

func (c *SearchController) fetchDetails(resp *ia_client.SearchResponse) (int, error) {
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
		d := c.ia.GetItemDetails(doc.Identifier)
		if d != nil {
			item.Server = d.Server
			item.Dir = d.Dir
			if len(doc.Creator) > 0 && doc.Creator[0] != "" {
				item.Creator = doc.Creator[0]
			} else if len(d.Metadata.Creator) > 0 && d.Metadata.Creator[0] != "" {
				item.Creator = d.Metadata.Creator[0]
			} else if len(d.Metadata.Artist) > 0 && d.Metadata.Artist[0] != "" {
				item.Creator = d.Metadata.Artist[0]
			} else {
				item.Creator = "Internet Archive"
			}

			if len(d.Metadata.Description) > 0 {
				item.Description = tview.Escape(c.ia.Html2Text(d.Metadata.Description[0]))
			}

			for name, metadata := range d.Files {
				format := metadata.Format
				// collect mp3 files

				if utils.Contains(Mp3Formats, format) {
					size, sErr := strconv.ParseInt(metadata.Size, 10, 64)
					length, lErr := utils.TimeToSeconds(metadata.Length)
					if sErr != nil || lErr != nil {
						logger.Error("Can't parse the file metadata: " + name)
						return 0, fmt.Errorf("can't parse file metadata: %s", name)
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
								oldFilePriority := utils.GetIndex(Mp3Formats, oldFile.Format)
								newFilePriority := utils.GetIndex(Mp3Formats, file.Format)
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
				if utils.Contains(CoverFormats, format) {
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
				c.totalItemsFetched++
				sp := &dto.SearchProgress{ItemsTotal: itemsTotal, ItemsFetched: c.totalItemsFetched}
				c.mq.SendMessage(mq.SearchController, mq.SearchPage, sp, mq.PriorityNormal)
				c.mq.SendMessage(mq.SearchController, mq.SearchPage, item, mq.PriorityNormal)
			}
		}
		logger.Debug(mq.SearchController + " fetched first " + strconv.Itoa(c.totalItemsFetched) + " items from " + strconv.Itoa(itemsTotal) + " total")
	}
	return itemsFetched, nil
}
