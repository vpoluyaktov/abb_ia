package cmd

import (
	"fmt"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/pkg/ia_client"
)

func Execute() {
  logger.Info("Application started") 

  ia := ia_client.New()

  res := ia.SearchByTitle("NASA", "audio")

  for i, doc := range res.Response.Docs {
    logger.Debug(fmt.Sprintf("%d - %s", i, doc.Title))
  } 

  res = ia.Search("https://archive.org/details/OTRR_Frank_Race_Singles", "audio") 

  for i, doc := range res.Response.Docs {
    logger.Debug(fmt.Sprintf("%d - %s", i, doc.Title))
  }

  logger.Info("Application finished") 

}