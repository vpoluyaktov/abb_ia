package controller

import (
	"path/filepath"

	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/ffmpeg"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type ChaptersController struct {
	mq       *mq.Dispatcher
	ab       *dto.Audiobook
	stopFlag bool
}

/**
 * Creates a new ChaptersController instance.
 * @param dispatcher - The message queue dispatcher.
 * @returns The new ChaptersController instance.
 *
 * This code is useful for creating a new ChaptersController instance and registering it with the message queue dispatcher. This allows the ChaptersController to receive messages from the message queue and dispatch them to the appropriate handler.
 **/
func NewChaptersController(dispatcher *mq.Dispatcher) *ChaptersController {
	dc := &ChaptersController{}
	dc.mq = dispatcher
	dc.mq.RegisterListener(mq.ChaptersController, dc.dispatchMessage)
	return dc
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
	case *dto.StopCommand:
		go c.stopChapters(dto)
	default:
		m.UnsupportedTypeError(mq.ChaptersController)
	}
}

func (c *ChaptersController) stopChapters(cmd *dto.StopCommand) {
	c.stopFlag = true
	logger.Debug(mq.ChaptersController + ": Received StopChapters command")
}

/**
 * @description Splits an audiobook into parts and chapters.
 * @param {dto.ChaptersCreate} cmd - The command to create chapters.
 * @returns {void}
 *
 * This function is useful for splitting an audiobook into parts and chapters.
 * It takes in a command object containing the audiobook and then splits the audiobook into parts and chapters.
 * It then sends messages to the ChaptersPage and Footer to update the status and busy indicator.
 */
func (c *ChaptersController) createChapters(cmd *dto.ChaptersCreate) {

	logger.Debug(mq.ChaptersController + " received " + cmd.String())

	c.ab = cmd.Audiobook
	c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.DisplayBookInfoCommand{Audiobook: c.ab}, true)
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.UpdateStatus{Message: "Calculating book parts and chapters..."}, false)
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)


	// Split the book into parts
	c.ab.Parts = []dto.Part{}
	var partNo int = 1
	var fileNo int = 1
	var chapterNo int = 1
	var offset float64 = 0
	var partSize int64 = 0
	var partDuration float64 = 0
	var partChapters []dto.Chapter = []dto.Chapter{}
	var chapterFiles []dto.Mp3File = []dto.Mp3File{}

	for i, file := range c.ab.Mp3Files {
		filePath := filepath.Join(c.ab.OutputDir, file.FileName)
		mp3, _ := ffmpeg.NewFFProbe(filePath)
		chapterFiles = append(chapterFiles, dto.Mp3File{Number: fileNo, FileName: file.FileName, Size: mp3.Size(), Duration: mp3.Duration()})
		fileNo++
		partSize += mp3.Size()
		partDuration += mp3.Duration()
		chapter := dto.Chapter{Number: chapterNo, Name: c.getMp3Title(mp3.Title()), Size: mp3.Size(), Duration: mp3.Duration(), Start: offset, End: offset + mp3.Duration(), Files: chapterFiles}
		partChapters = append(partChapters, chapter)
		c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.AddChapterCommand{Chapter: &chapter}, true)
		offset += mp3.Duration()
		chapterNo++
		chapterFiles = []dto.Mp3File{}
		if partSize >= config.MaxFileSize() || i == len(c.ab.Mp3Files)-1 {
			part := dto.Part{Number: partNo, Size: partSize, Duration: partDuration, Chapters: partChapters}
			c.ab.Parts = append(c.ab.Parts, part)
			partNo++
			fileNo = 1
			partSize = 0
			partDuration = 0
			offset = 0
			partChapters = []dto.Chapter{}
		}
	}

	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.ChaptersReady{Audiobook: cmd.Audiobook}, true)
}

func (c *ChaptersController) getMp3Title(title string) string {
	// TODO: Concatenate Chapter names here
	return title
}
