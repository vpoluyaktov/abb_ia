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

  item := ia.GetItemById("OTRR_Frank_Race_Singles")
  if item!= nil {
    fmt.Printf("Title: %s\n", item.Metadata.Title)
    fmt.Printf("Server: %s\n", item.Server)
    fmt.Printf("Directory: %s\n", item.Dir)
    fmt.Printf("Description: %s\n", item.Metadata.Description)
    fmt.Printf("Creator: %s\n", item.Metadata.Creator)
    fmt.Printf("Image: %s\n", item.Misc.Image)
    
    for file, meta := range item.Files {
			fmt.Printf("%s -> %s\n", file, meta.Format)
		} 

  }

  logger.Info("Application finished") 

}
