package controller

import (
	"fmt"
	"time"

	"abb_ia/internal/audiobookshelf"
	"abb_ia/internal/dto"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"
)

type UploadController struct {
	mq        *mq.Dispatcher
	ab        *dto.Audiobook
	startTime time.Time
	stopFlag  bool

	// progress tracking arrays
	filesUpload []fileUpload
}

type fileUpload struct {
	fileId      int
	fileSize    int64
	bytesCopied int64
	progress    int
}

func NewUploadController(dispatcher *mq.Dispatcher) *UploadController {
	c := &UploadController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.UploadController, c.dispatchMessage)
	return c
}

func (c *UploadController) checkMQ() {
	m := c.mq.GetMessage(mq.UploadController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *UploadController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.AbsScanCommand:
		go c.absScan(dto)
	case *dto.AbsUploadCommand:
		go c.absUpload(dto)
	default:
		m.UnsupportedTypeError(mq.UploadController)
	}
}

func (c *UploadController) absScan(cmd *dto.AbsScanCommand) {
	logger.Info(mq.UploadController + " received " + cmd.String())
	ab := cmd.Audiobook
	url := ab.Config.GetAudiobookshelfUrl()
	username := ab.Config.GetAudiobookshelfUser()
	password := ab.Config.GetAudiobookshelfPassword()
	libraryName := ab.Config.GetAudiobookshelfLibrary()

	if url != "" && username != "" && password != "" && libraryName != "" {
		absClient := audiobookshelf.NewClient(url)
		err := absClient.Login(username, password)
		if err != nil {
			logger.Error("Can't login to audiobookshlf server: " + err.Error())
			return
		}
		libraries, err := absClient.GetLibraries()
		if err != nil {
			logger.Error("Can't get a list of libraries from audiobookshlf server: " + err.Error())
			return
		}
		libraryID, err := absClient.GetLibraryId(libraries, libraryName)
		if err != nil {
			logger.Error("Can't find audiobookshlf library by name: " + err.Error())
			return
		}
		err = absClient.ScanLibrary(libraryID)
		if err != nil {
			logger.Error("Can't launch library scan on audiobookshlf server: " + err.Error())
			return
		}
		if err != nil {
			logger.Error("Can't launch library scan on audiobookshlf server: " + err.Error())
			return
		}
		logger.Info("A scan launched for library " + libraryName + " on audiobookshlf server")
	}
	c.mq.SendMessage(mq.UploadController, mq.BuildPage, &dto.ScanComplete{Audiobook: cmd.Audiobook}, true)
}

func (c *UploadController) absUpload(cmd *dto.AbsUploadCommand) {
	c.startTime = time.Now()
	logger.Info(mq.UploadController + " received " + cmd.String())
	c.ab = cmd.Audiobook
	url := c.ab.Config.GetAudiobookshelfUrl()
	username := c.ab.Config.GetAudiobookshelfUser()
	password := c.ab.Config.GetAudiobookshelfPassword()
	libraryName := c.ab.Config.GetAudiobookshelfLibrary()

	if url != "" && username != "" && password != "" && libraryName != "" {
		absClient := audiobookshelf.NewClient(url)
		err := absClient.Login(username, password)
		if err != nil {
			logger.Error("Can't login to audiobookshelf server: " + err.Error())
			return
		}
		libraries, err := absClient.GetLibraries()
		if err != nil {
			logger.Error("Can't get a list of libraries from audiobookshelf server: " + err.Error())
			return
		}
		libraryID, err := absClient.GetLibraryId(libraries, libraryName)
		if err != nil {
			logger.Error("Can't find audiobookshelf library by name: " + err.Error())
			return
		}
		folders, err := absClient.GetFolders(libraries, libraryName)
		if err != nil || len(folders) == 0 {
			logger.Error("Can't get a folder for library: " + err.Error())
			return
		}
		// TODO: Check if a folder selector is needed here. Let's use first folder in a library for upload
		folderID := folders[0].ID

		c.stopFlag = false
		c.filesUpload = make([]fileUpload, len(c.ab.Parts))
		go c.updateTotalUploadProgress()
		err = absClient.UploadBook(c.ab, libraryID, folderID, c.updateFileUplodProgress)
		if err != nil {
			logger.Error("Can't upload the audiobook to audiobookshelf server: " + err.Error())
		}
		c.stopFlag = true
	}
	c.mq.SendMessage(mq.UploadController, mq.BuildPage, &dto.UploadComplete{Audiobook: cmd.Audiobook}, true)
}

func (c *UploadController) updateFileUplodProgress(fileId int, fileName string, size int64, pos int64, percent int) {

	if c.filesUpload[fileId].progress != percent {
		// wrong calculation protection
		if percent < 0 {
			percent = 0
		} else if percent > 100 {
			percent = 100
		}

		// sent a message only if progress changed
		c.mq.SendMessage(mq.UploadController, mq.BuildPage, &dto.UploadFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, false)
	}
	c.filesUpload[fileId].fileId = fileId
	c.filesUpload[fileId].fileSize = size
	c.filesUpload[fileId].bytesCopied = pos
	c.filesUpload[fileId].progress = percent
}

func (c *UploadController) updateTotalUploadProgress() {
	var percent int = -1

	for !c.stopFlag && percent <= 100 {
		var totalSize = c.ab.TotalSize
		var totalBytesCopied int64 = 0
		filesCopied := 0
		for _, f := range c.filesUpload {
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
		if filesCopied == len(c.filesUpload) {
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

			c.mq.SendMessage(mq.UploadController, mq.BuildPage, &dto.UploadProgress{Elapsed: elapsedH, Percent: percent, Files: filesH, Bytes: bytesH, Speed: speedH, ETA: etaH}, false)
		}
		time.Sleep(mq.PullFrequency)
	}
}
