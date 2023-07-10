package cmd

import (
	"github.com/vpoluyaktov/abb_ia/internal/controller"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
	"github.com/vpoluyaktov/abb_ia/internal/ui"
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
