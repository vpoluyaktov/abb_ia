package controller

import (
	"strconv"

	"github.com/vpoluyaktov/abb_ia/internal/audiobookshelf"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type AudiobookshelfController struct {
	mq *mq.Dispatcher
}

func NewAudiobookshelfController(dispatcher *mq.Dispatcher) *AudiobookshelfController {
	c := &AudiobookshelfController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.AudiobookshelfController, c.dispatchMessage)
	return c
}

func (c *AudiobookshelfController) checkMQ() {
	m := c.mq.GetMessage(mq.AudiobookshelfController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *AudiobookshelfController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.AudiobookshelfScanCommand:
		go c.audiobookshelfScan(dto)
	case *dto.AudiobookshelfUploadCommand:
		go c.uploadAudiobook(dto)
	default:
		m.UnsupportedTypeError(mq.AudiobookshelfController)
	}
}

func (c *AudiobookshelfController) audiobookshelfScan(cmd *dto.AudiobookshelfScanCommand) {
	logger.Info(mq.AudiobookshelfController + " received " + cmd.String())
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
}

func (c *AudiobookshelfController) uploadAudiobook(cmd *dto.AudiobookshelfUploadCommand) {
	logger.Info(mq.AudiobookshelfController + " received " + cmd.String())
	ab := cmd.Audiobook
	url := ab.Config.GetAudiobookshelfUrl()
	username := ab.Config.GetAudiobookshelfUser()
	password := ab.Config.GetAudiobookshelfPassword()
	libraryName := ab.Config.GetAudiobookshelfLibrary()

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
		err = absClient.UploadBook(ab, libraryID, folders[0].ID, c.updateFileUplodProgress)
		if err != nil {
			logger.Error("Can't upload the audiobook to audiobookshelf server: " + err.Error())
			return
		}
	}
}

func (c *AudiobookshelfController) updateFileUplodProgress(bytesUploaded int, totalBytes int) {

	if totalBytes != 0 {
		percent := bytesUploaded / totalBytes * 100
		logger.Debug("Upload percent: " + strconv.Itoa(percent))
	}
}
