package cmd

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/actions"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/event"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/ui"
)

func Execute() {
	logger.Info("Application started")
	
  d := event.NewDispatcher()
  ui := ui.NewTUI(d)
	c := controller.NewController(d)

  c.Run()
	ui.Run()
	logger.Info("Application finished")
}
