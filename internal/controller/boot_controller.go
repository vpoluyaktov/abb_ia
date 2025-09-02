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
	"fmt"
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
	m, err := c.mq.GetMessage(mq.BootController)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get message for BootController: %v", err))
		return
	}
	if m == nil {
		return // No message available
	}
	c.dispatchMessage(m)
}

func (c *BootController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.BootstrapCommand:
		go c.bootStrap(dto)
	default:
		m.UnsupportedTypeError(mq.BootController)
	}
}

func (c *BootController) bootStrap(_ *dto.BootstrapCommand) {
	// Detect operation system
	os := runtime.GOOS
	logger.Debug(fmt.Sprintf("Operating system detected: %s", os))

	switch os {
	case "windows", "darwin", "linux":
		logger.Info(fmt.Sprintf("Starting application on %s", os))
	default:
		logger.Error(fmt.Sprintf("Unsupported operating system: %s", os))
		c.mq.SendMessage(mq.BootController, mq.SearchPage, &dto.Error{
			Message: fmt.Sprintf("Unsupported operating system: %s", os),
		}, mq.PriorityCritical)
		return
	}

	// Wait for all components to initialize
	time.Sleep(3 * time.Second)

	// Check dependencies and version
	if c.checkFFmpeg() {
		c.checkNewVersion()
	} else {
		logger.Error("Failed to verify FFmpeg installation")
	}
}

func (c *BootController) checkFFmpeg() bool {
	if !(utils.CommandExists("ffmpeg") && utils.CommandExists("ffprobe")) {
		logger.Fatal("Bootstrap: ffmpeg or ffprobe command not found")
		c.mq.SendMessage(mq.BootController, mq.SearchPage, &dto.FFMPEGNotFoundError{}, mq.PriorityNormal)
		return false
	}
	return true
}

func (c *BootController) checkNewVersion() {
	appVersion := config.Instance().AppVersion()
	if appVersion == "0.0.0" {
		logger.Debug("Local development version detected, skipping version check")
		return
	}

	git := github.NewClient(config.Instance().GetRepoOwner(), config.Instance().GetRepoName())
	latestVersion, err := git.GetLatestVersion()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to check for new version: %v", err))
		c.mq.SendMessage(mq.BootController, mq.SearchPage, &dto.Error{Message: "Failed to check for updates"}, mq.PriorityNormal)
		return
	}

	result, err := github.CompareVersions(latestVersion, appVersion)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to compare versions %s and %s: %v", latestVersion, appVersion, err))
		return
	}

	if result > 0 {
		logger.Info(fmt.Sprintf("New version available: %s (current: %s)", latestVersion, appVersion))
		c.mq.SendMessage(mq.BootController, mq.SearchPage, &dto.NewAppVersionFound{
			CurrentVersion: appVersion,
			NewVersion:    latestVersion,
		}, mq.PriorityNormal)
	} else {
		logger.Debug(fmt.Sprintf("Application is up to date (version %s)", appVersion))
	}
}
