package controller

import (
	"abb_ia/internal/config"
	"abb_ia/internal/dto"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"
)

type ConfigController struct {
	mq *mq.Dispatcher
}

func NewConfigController(dispatcher *mq.Dispatcher) *ConfigController {
	c := &ConfigController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.ConfigController, c.dispatchMessage)
	return c
}

func (c *ConfigController) checkMQ() {
	m := c.mq.GetMessage(mq.ConfigController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *ConfigController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.SaveConfigCommand:
		go c.updateConfig(dto)
	default:
		m.UnsupportedTypeError(mq.ConfigController)
	}
}

func (c *ConfigController) updateConfig(cmd *dto.SaveConfigCommand) {
	logger.Debug(mq.ConfigController + ": Received UpdateConfigCommand")
	c.mq.SendMessage(mq.ConfigController, mq.Footer, &dto.UpdateStatus{Message: "Saving new default configuration"}, false)
	c.mq.SendMessage(mq.ConfigController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)

	config.SaveConfig(&cmd.Config)

	logger.SetLogLevel(logger.LogLevelType(utils.GetIndex(logger.LogLeves(), cmd.Config.GetLogLevel()) + 1))
	c.mq.SendMessage(mq.ConfigController, mq.SearchPage, &dto.UpdateSearchConfigCommand{Config: cmd.Config}, true)

	c.mq.SendMessage(mq.ConfigController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.ConfigController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
}
