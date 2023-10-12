package controller

import (
	"time"

	"github.com/vpoluyaktov/abb_ia/internal/dto"

	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type BuildController struct {
	mq        *mq.Dispatcher
	item      *dto.IAItem
	startTime time.Time
	files     []fileBuild
	stopFlag  bool
}

type fileBuild struct {
	fileId          int
	fileSize        int64
	bytesDownloaded int64
	progress        int
}

func NewBuildController(dispatcher *mq.Dispatcher) *BuildController {
	dc := &BuildController{}
	dc.mq = dispatcher
	dc.mq.RegisterListener(mq.BuildController, dc.dispatchMessage)
	return dc
}

func (c *BuildController) checkMQ() {
	m := c.mq.GetMessage(mq.BuildController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *BuildController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.BuildCommand:
		go c.startBuild(dto)
	case *dto.CopyCommand:
		go c.startCopy(dto)
	case *dto.StopCommand:
		go c.stopBuild(dto)
	default:
		m.UnsupportedTypeError(mq.BuildController)
	}
}

func (c *BuildController) stopBuild(cmd *dto.StopCommand) {
	c.stopFlag = true
	logger.Debug(mq.BuildController + ": Received StopBuild command")
}

func (c *BuildController) startBuild(cmd *dto.BuildCommand) {
	c.startTime = time.Now()
	logger.Info(mq.BuildController + " received " + cmd.String())
	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.UpdateStatus{Message: "Building mp3 files..."}, false)
	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)
	c.item = cmd.Audiobook.IAItem
	// outputDir := filepath.Join("output", c.item.ID)

	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.BuildComplete{Audiobook: cmd.Audiobook}, true)
}

func (c *BuildController) startCopy(cmd *dto.CopyCommand) {

}

func (c *BuildController) updateBuildProgress(fileId int, fileName string, size int64, pos int64, percent int) {
	if c.files[fileId].progress != percent {
		// sent a message only if progress changed
		c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, false)
	}
	c.files[fileId].fileId = fileId
	c.files[fileId].fileSize = size
	c.files[fileId].bytesDownloaded = pos
	c.files[fileId].progress = percent
}

func (c *BuildController) updateCopyProgress(fileId int, fileName string, size int64, pos int64, percent int) {
	if c.files[fileId].progress != percent {
		// sent a message only if progress changed
		c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, false)
	}
	c.files[fileId].fileId = fileId
	c.files[fileId].fileSize = size
	c.files[fileId].bytesDownloaded = pos
	c.files[fileId].progress = percent
}
