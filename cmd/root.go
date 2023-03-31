package cmd

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
  "github.com/vpoluyaktov/audiobook_creator_IA/internal/ui"
)

func Execute() {
  logger.Info("Application started") 
  ui := ui.NewTUI()
  ui.Run()
  logger.Info("Application finished") 
}
