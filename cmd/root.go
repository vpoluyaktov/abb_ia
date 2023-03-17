package cmd

import (
	log "github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/pkg/ia_client"
)

func Execute() {
  log.Info("Application started") 

  ia := ia_client.New()

  ia.Search("NASA")

  log.Info("Application finished") 


}