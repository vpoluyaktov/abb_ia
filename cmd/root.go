package cmd

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/controller"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/ui"
)

func Execute() {
	logger.Info("Application started")

	d := mq.NewDispatcher()
	c := controller.NewConductor(d)
	ui := ui.NewTUI(d)

	c.Run()
	ui.Run()
	logger.Info("Application finished")
}
