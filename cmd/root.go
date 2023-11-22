package cmd

import (
	"abb_ia/internal/controller"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
	"abb_ia/internal/ui"
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
