package dto

import "fmt"

type DownloadCommand struct {
	Item *IAItem
}

func (c *DownloadCommand) String() string {
	return fmt.Sprintf("DownloadCommand: %s", c.Item.String())
}

