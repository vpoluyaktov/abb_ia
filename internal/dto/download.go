package dto

import "fmt"

type DownloadCommand struct {
	item IAItem
}

func (c *DownloadCommand) String() string {
	return fmt.Sprintf("DownloadCommand: %s", c.item.String())
}

