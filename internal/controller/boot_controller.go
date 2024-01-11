package controller

import (
	"runtime"
	"time"

	"abb_ia/internal/config"
	"abb_ia/internal/dto"
	"abb_ia/internal/github"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"
)

type BootController struct {
	mq *mq.Dispatcher
}

func NewBootController(dispatcher *mq.Dispatcher) *BootController {
	c := &BootController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.BootController, c.dispatchMessage)
	go c.bootStrap(&dto.BootstrapCommand{})
	return c
}

func (c *BootController) checkMQ() {
	m := c.mq.GetMessage(mq.BootController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *BootController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.BootstrapCommand:
		go c.bootStrap(dto)
	default:
		m.UnsupportedTypeError(mq.BootController)
	}
}

func (c *BootController) bootStrap(cmd *dto.BootstrapCommand) {

	// detect operation system
	os := runtime.GOOS
	logger.Debug("Operation system detected: " + os)
	switch os {
	case "windows":
	case "darwin":
	case "linux":
	default:
		logger.Error("Unknown operation system detected: " + os)
	}

	// wait for all components to initialize
	time.Sleep(3 * time.Second)
	if c.checkFFmpeg() {
		c.checkNewVersion()
	}
}

func (c *BootController) checkFFmpeg() bool {
	if !(utils.CommandExists("ffmpeg") && utils.CommandExists("ffprobe")) {
		logger.Fatal("Bootstrap: ffmpeg or ffprobe command not found")
		c.mq.SendMessage(mq.BootController, mq.SearchPage, &dto.FFMPEGNotFoundError{}, true)
		return false
	}
	return true
}

func (c *BootController) checkNewVersion() {

	latestVersion, err := github.GetLatestVersion(config.Instance().GetRepoOwner(), config.Instance().GetRepoName())
	if err != nil {
		logger.Error("Can't check new version: " + err.Error())
		return
	}

	result, err := github.CompareVersions(latestVersion, config.Instance().AppVersion())
	if err != nil {
		logger.Error("Can not compare versions: " + err.Error())
		return
	}

	if result > 0 {
		c.mq.SendMessage(mq.BootController, mq.SearchPage, &dto.NewAppVersionFound{CurrentVersion: config.Instance().AppVersion(), NewVersion: latestVersion}, true)
	}
}
