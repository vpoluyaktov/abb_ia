package dto

import "fmt"

type AudiobookshelfScanCommand struct {
	Audiobook *Audiobook
}

func (c *AudiobookshelfScanCommand) String() string {
	return fmt.Sprintf("AudiobookshelfScanCommand: %s", c.Audiobook.String())
}