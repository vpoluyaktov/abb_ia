package controller

import (
	"github.com/vpoluyaktov/abb_ia/internal/audiobookshelf"
	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type AudiobookshelfController struct {
	mq *mq.Dispatcher
}

func NewAudiobookshelfController(dispatcher *mq.Dispatcher) *AudiobookshelfController {
	dc := &AudiobookshelfController{}
	dc.mq = dispatcher
	dc.mq.RegisterListener(mq.AudiobookshelfController, dc.dispatchMessage)
	return dc
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
	default:
		m.UnsupportedTypeError(mq.AudiobookshelfController)
	}
}

func (c *AudiobookshelfController) audiobookshelfScan(cmd *dto.AudiobookshelfScanCommand) {
	logger.Info(mq.AudiobookshelfController + " received " + cmd.String())
	url := config.Instance().GetAudiobookshelfUrl()
	username := config.Instance().GetAudiobookshelfUser()
	password := config.Instance().GetAudiobookshelfPassword()
	libraryName := config.Instance().GetAudiobookshelfLibrary()

	if url != "" && username != "" && password != "" && libraryName != "" {
		loginResp, err := audiobookshelf.Login(url+"/login", username, password)
		if err != nil {
			logger.Error("Can't login to audiobookshlf server: " + err.Error())
			return
		}
		libraryResponse, err := audiobookshelf.Libraries(url, loginResp.User.Token)
		if err != nil {
			logger.Error("Can't get a list of libraries from audiobookshlf server: " + err.Error())
			return
		}
		libraryID, err := audiobookshelf.GetLibraryByName(libraryResponse.Libraries, libraryName)
		if err != nil {
			logger.Error("Can't find audiobookshlf library by name: " + err.Error())
			return
		}
		err = audiobookshelf.ScanLibrary(url, loginResp.User.Token, libraryID)
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
