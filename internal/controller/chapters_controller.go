package controller

import (
	"path/filepath"
	"regexp"
	"strings"

	"abb_ia/internal/dto"
	"abb_ia/internal/ffmpeg"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
)

type ChaptersController struct {
	mq       *mq.Dispatcher
	ab       *dto.Audiobook
	stopFlag bool
}

func NewChaptersController(dispatcher *mq.Dispatcher) *ChaptersController {
	c := &ChaptersController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.ChaptersController, c.dispatchMessage)
	return c
}

func (c *ChaptersController) checkMQ() {
	m := c.mq.GetMessage(mq.ChaptersController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *ChaptersController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.ChaptersCreate:
		go c.createChapters(dto)
	case *dto.SearchReplaceDescriptionCommand:
		go c.searchReplaceDescription(dto)
	case *dto.SearchReplaceChaptersCommand:
		go c.searchReplaceChapters(dto)
	case *dto.JoinChaptersCommand:
		go c.joinChapters(dto)
	case *dto.RecalculatePartsCommand:
		go c.recalculateParts(dto)
	case *dto.StopCommand:
		go c.stopChapters(dto)
	case *dto.UseMP3NamesCommand:
		go c.useMP3Names(dto)
	default:
		m.UnsupportedTypeError(mq.ChaptersController)
	}
}

func (c *ChaptersController) stopChapters(cmd *dto.StopCommand) {
	c.stopFlag = true
	logger.Debug(mq.ChaptersController + ": Received StopChapters command")
}

func (c *ChaptersController) createChapters(cmd *dto.ChaptersCreate) {
	logger.Debug(mq.ChaptersController + " received " + cmd.String())
	c.stopFlag = false
	c.ab = cmd.Audiobook

	if c.ab.Config.IsShortenTitle() {
		for _, pair := range c.ab.Config.ShortenPairs {
			c.ab.Title = strings.ReplaceAll(c.ab.Title, pair.Search, pair.Replace)
			c.ab.Author = strings.ReplaceAll(c.ab.Author, pair.Search, pair.Replace)
		}
	}

	c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.DisplayBookInfoCommand{Audiobook: c.ab}, true)
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.UpdateStatus{Message: "Calculating book parts and chapters..."}, false)
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)

	// Split the book into parts
	c.ab.Parts = []dto.Part{}
	var partNo int = 1
	var fileNo int = 1
	var chapterNo int = 1
	var offset float64 = 0
	var abSize int64 = 0
	var abDuration float64 = 0
	var partSize int64 = 0
	var partDuration float64 = 0
	var partChapters []dto.Chapter = []dto.Chapter{}
	var chapterFiles []dto.Mp3File = []dto.Mp3File{}

	for i, file := range c.ab.Mp3Files {
		filePath := filepath.Join(c.ab.OutputDir, file.FileName)
		mp3, _ := ffmpeg.NewFFProbe(filePath)
		chapterFiles = append(chapterFiles, dto.Mp3File{Number: fileNo, FileName: file.FileName, Size: mp3.Size(), Duration: mp3.Duration()})
		fileNo++
		abSize += mp3.Size()
		abDuration += mp3.Duration()
		partSize += mp3.Size()
		partDuration += mp3.Duration()
		chapter := dto.Chapter{Number: chapterNo, Name: mp3.Title(), Size: mp3.Size(), Duration: mp3.Duration(), Start: offset, End: offset + mp3.Duration(), Files: chapterFiles}
		partChapters = append(partChapters, chapter)
		c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.AddChapterCommand{Chapter: &chapter}, true)
		offset += mp3.Duration()
		chapterNo++
		chapterFiles = []dto.Mp3File{}
		if partSize >= int64(c.ab.Config.GetMaxFileSizeMb())*1024*1024 || i == len(c.ab.Mp3Files)-1 {
			part := dto.Part{Number: partNo, Format: mp3.Format(), Size: partSize, Duration: partDuration, Chapters: partChapters}
			c.ab.Parts = append(c.ab.Parts, part)
			partNo++
			fileNo = 1
			partSize = 0
			partDuration = 0
			offset = 0
			partChapters = []dto.Chapter{}
		}
	}

	// update the audiobook size and duration
	c.ab.TotalSize = abSize
	c.ab.TotalDuration = abDuration

	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	if !c.stopFlag {
		c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.ChaptersReady{Audiobook: cmd.Audiobook}, true)
	}
	c.stopFlag = true
}

func (c *ChaptersController) searchReplaceDescription(cmd *dto.SearchReplaceDescriptionCommand) {
	ab := cmd.Audiobook
	searchStr := cmd.SearchStr
	replaceStr := cmd.ReplaceStr
	re, err := regexp.Compile(searchStr)
	if err != nil {
		return
	}

	description := re.ReplaceAllString(ab.Description, replaceStr)
	ab.Description = description
	c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.RefreshDescriptionCommand{Audiobook: cmd.Audiobook}, true)
}

func (c *ChaptersController) searchReplaceChapters(cmd *dto.SearchReplaceChaptersCommand) {
	ab := cmd.Audiobook
	searchStr := cmd.SearchStr
	replaceStr := cmd.ReplaceStr
	re, err := regexp.Compile(searchStr)
	if err != nil {
		return
	}

	for partNo, p := range ab.Parts {
		for chapterNo := range p.Chapters {
			chapter := &ab.Parts[partNo].Chapters[chapterNo]
			title := chapter.Name
			title = re.ReplaceAllString(title, replaceStr)
			chapter.Name = strings.TrimSpace(title)
		}
	}
	c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.RefreshChaptersCommand{Audiobook: cmd.Audiobook}, true)
}

// Join Chapters having the same name
func (c *ChaptersController) joinChapters(cmd *dto.JoinChaptersCommand) {
	ab := cmd.Audiobook
	chapterNo := 1
	for partNo := range ab.Parts {
		part := &ab.Parts[partNo]
		chapters := []dto.Chapter{}
		var chapter *dto.Chapter = nil
		previousChapterName := ""
		for chNo := range part.Chapters {
			ch := &part.Chapters[chNo]
			// always add first chapter in a part
			if chNo == 0 {
				chapter = ch
				chapter.Number = chapterNo
				previousChapterName = chapter.Name
			} else {
				if ch.Name == previousChapterName {
					// the same name - extend current chapter
					chapter.Duration += ch.Duration
					chapter.Size += ch.Size
					chapter.End = ch.End
					chapter.Files = append(chapter.Files, ch.Files...)
				} else {
					// new chapter
					chapters = append(chapters, *chapter)
					chapterNo++
					chapter = ch
					chapter.Number = chapterNo
					previousChapterName = chapter.Name

				}
			}
			// add last chapter in a part
			if chNo == len(part.Chapters)-1 {
				chapters = append(chapters, *chapter)
				chapterNo++
			}
		}
		part.Chapters = chapters
	}
	c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.RefreshChaptersCommand{Audiobook: cmd.Audiobook}, true)
}

// Recalculate Parts using new PartSize
func (c *ChaptersController) recalculateParts(cmd *dto.RecalculatePartsCommand) {
	ab := cmd.Audiobook
	c.createChapters(&dto.ChaptersCreate{Audiobook: ab})
}

// useMP3Names replaces chapter names with their corresponding MP3 file names
func (c *ChaptersController) useMP3Names(cmd *dto.UseMP3NamesCommand) {
	ab := cmd.Audiobook

	for partNo, p := range ab.Parts {
		for chapterNo := range p.Chapters {
			chapter := &ab.Parts[partNo].Chapters[chapterNo]
			// Since each chapter can have multiple files, we'll use the name of the first file
			if len(chapter.Files) > 0 {
				// Get just the filename without path and extension
				baseName := filepath.Base(chapter.Files[0].FileName)
				fileName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
				chapter.Name = fileName
			}
		}
	}

	c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.RefreshChaptersCommand{Audiobook: cmd.Audiobook}, true)
}
