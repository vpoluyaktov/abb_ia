package controller

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"abb_ia/internal/audiobookshelf"
	"abb_ia/internal/dto"
	"abb_ia/internal/utils"

	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
)

/**
 * CopyController doesn't have its own UI page.
 * Instead it uses button half of the BuildPage
 **/
type CopyController struct {
	mq        *mq.Dispatcher
	ab        *dto.Audiobook
	startTime time.Time
	stopFlag  bool

	// progress tracking arrays
	filesCopy []fileCopy
}

type fileCopy struct {
	fileId      int
	fileSize    int64
	bytesCopied int64
	progress    int
}

// Progress Reader for file copy progress
type Fn func(fileId int, fileName string, size int64, pos int64, percent int)
type ProgressReader struct {
	FileId   int
	FileName string
	Reader   io.Reader
	Size     int64
	Pos      int64
	Percent  int
	Callback Fn
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if err == nil {
		pr.Pos += int64(n)
		pr.Percent = int(float64(pr.Pos) / float64(pr.Size) * 100)
		pr.Callback(pr.FileId, pr.FileName, pr.Size, pr.Pos, pr.Percent)
	}
	return n, err
}

func NewCopyController(dispatcher *mq.Dispatcher) *CopyController {
	c := &CopyController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.CopyController, c.dispatchMessage)
	return c
}

func (c *CopyController) checkMQ() {
	m, err := c.mq.GetMessage(mq.CopyController)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get message for CopyController: %v", err))
		return
	}
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *CopyController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.CopyCommand:
		go c.startCopy(dto)
	case *dto.StopCommand:
		go c.stopCopy(dto)
	default:
		m.UnsupportedTypeError(mq.CopyController)
	}
}

func (c *CopyController) startCopy(cmd *dto.CopyCommand) {
	c.startTime = time.Now()
	c.ab = cmd.Audiobook

	// update part size and whole audiobook size after re-encoding
	abSize := int64(0)
	for i := range c.ab.Parts {
		part := &c.ab.Parts[i]
		fileInfo, err := os.Stat(part.M4BFile)
		if err != nil {
			logger.Error("Can't open .mb4 file: " + err.Error())
			return
		}
		// Get file size in bytes
		part.Size = fileInfo.Size()
		abSize += part.Size
	}
	c.ab.TotalSize = abSize

	c.mq.SendMessage(mq.CopyController, mq.Footer, &dto.UpdateStatus{Message: "Copying audiobook to the output directory..."}, mq.PriorityNormal)
	c.mq.SendMessage(mq.CopyController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, mq.PriorityNormal)

	logger.Info(fmt.Sprintf("Copying the audiobook: %s - %s to %s/...", c.ab.Author, c.ab.Title, c.ab.Config.OutputDir))

	c.stopFlag = false
	c.filesCopy = make([]fileCopy, len(c.ab.Parts))
	jd := utils.NewJobDispatcher(c.ab.Config.GetConcurrentDownloaders())
	for i := range c.ab.Parts {
		jd.AddJob(i, c.copyAudiobookPart, c.ab, i)
	}
	go c.updateTotalCopyProgress()
	jd.Start()

	c.mq.SendMessage(mq.CopyController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, mq.PriorityNormal)
	c.mq.SendMessage(mq.CopyController, mq.Footer, &dto.UpdateStatus{Message: ""}, mq.PriorityNormal)
	if !c.stopFlag {
		c.mq.SendMessage(mq.CopyController, mq.BuildPage, &dto.CopyComplete{Audiobook: cmd.Audiobook}, mq.PriorityHigh)
	}
	c.stopFlag = true
}

func (c *CopyController) stopCopy(_ *dto.StopCommand) {
	c.stopFlag = true
	logger.Debug(mq.CopyController + ": Received StopCopy command")
}

func (c *CopyController) copyAudiobookPart(ab *dto.Audiobook, partId int) {

	part := &ab.Parts[partId]

	file, err := os.Open(part.M4BFile)
	if err != nil {
		logger.Error("Can't open .mb4 file: " + err.Error())
		return
	}
	fileReader := bufio.NewReader(file)
	defer file.Close()

	// Calculate Audiobookshelf directory structure (see: https://www.audiobookshelf.org/docs#book-directory-structure)
	destPath := audiobookshelf.GetDestignationPath(ab.Config.GetOutputDir(), ab.Series, ab.Author)
	destDir := audiobookshelf.GetDestignationDir(ab.Series, ab.SeriesNo, ab.Title, ab.Narrator)

	filePath := filepath.Clean(filepath.Join(destPath, destDir, filepath.Base(part.M4BFile)))
	fullPath := filepath.Dir(filePath)

	if err := os.MkdirAll(fullPath, 0750); err != nil {
		logger.Error("Can't create output directory: " + err.Error())
		return
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal("Can't create Audiobookshelf .m4b file: " + err.Error())
		return
	}
	defer f.Close()

	progressReader := &ProgressReader{
		FileId:   partId,
		FileName: part.M4BFile,
		Reader:   fileReader,
		Size:     part.Size,
		Callback: c.updateFileCopyProgress,
	}

	if _, err := io.Copy(f, progressReader); err != nil {
		logger.Error("Error while copying .m4b file: " + err.Error())
	}
}

func (c *CopyController) updateFileCopyProgress(fileId int, fileName string, size int64, pos int64, percent int) {
	if c.filesCopy[fileId].progress != percent {

		// wrong calculation protection
		if percent < 0 {
			percent = 0
		} else if percent > 100 {
			percent = 100
		}

		// sent a message only if progress changed
		c.mq.SendMessage(mq.CopyController, mq.BuildPage, &dto.CopyFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, mq.PriorityNormal)
	}
	c.filesCopy[fileId].fileId = fileId
	c.filesCopy[fileId].fileSize = size
	c.filesCopy[fileId].bytesCopied = pos
	c.filesCopy[fileId].progress = percent
}

func (c *CopyController) updateTotalCopyProgress() {
	var percent int = -1

	for !c.stopFlag && percent <= 100 {
		var totalSize = c.ab.TotalSize
		var totalBytesCopied int64 = 0
		filesCopied := 0
		for _, f := range c.filesCopy {
			totalBytesCopied += f.bytesCopied
			if f.progress == 100 {
				filesCopied++
			}
		}

		var p int = 0
		if totalSize > 0 {
			p = int(float64(totalBytesCopied) / float64(totalSize) * 100)
		}

		// fix wrong incorrect calculation
		if filesCopied == len(c.filesCopy) {
			p = 100
			totalBytesCopied = c.ab.TotalSize
		}

		if percent != p {
			// sent a message only if progress changed
			percent = p

			// wrong calculation protection
			if percent < 0 {
				percent = 0
			} else if percent > 100 {
				percent = 100
			}

			elapsed := time.Since(c.startTime).Seconds()
			speed := int64(float64(totalBytesCopied) / elapsed)
			eta := (100 / (float64(percent) / elapsed)) - elapsed
			if eta < 0 || eta > (60*60*24*365) {
				eta = 0
			}

			elapsedH := utils.SecondsToTime(elapsed)
			bytesH := utils.BytesToHuman(totalBytesCopied)
			filesH := fmt.Sprintf("%d/%d", filesCopied, len(c.ab.Parts))
			speedH := utils.SpeedToHuman(speed)
			etaH := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.CopyController, mq.BuildPage, &dto.CopyProgress{Elapsed: elapsedH, Percent: percent, Files: filesH, Bytes: bytesH, Speed: speedH, ETA: etaH}, mq.PriorityNormal)
		}
		time.Sleep(mq.PullFrequency)
	}
}
