package cmd

import (
	"fmt"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/pkg/ia_client"
)

func Execute() {
  logger.Info("Application started") 

  ia := ia_client.New()

  res := ia.Search("NASA", "audio")

  for i, doc := range res.Response.Docs {
    fmt.Printf("%d - %s (%s)\n", i + 1, doc.Title, doc.Identifier)
  } 

 

  logger.Info("Application finished") 

}
