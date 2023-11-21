package dto

import "fmt"

type AudiobookshelfScanCommand struct {
	Audiobook *Audiobook
}

func (c *AudiobookshelfScanCommand) String() string {
	return fmt.Sprintf("AudiobookshelfScanCommand: %s", c.Audiobook.String())
}

type AudiobookshelfUploadCommand struct {
	Audiobook *Audiobook
}

func (c *AudiobookshelfUploadCommand) String() string {
	return fmt.Sprintf("AudiobookshelfUploadCommand: %s", c.Audiobook.String())
}

type UploadComplete struct {
	Audiobook *Audiobook
}

func (c *UploadComplete) String() string {
	return fmt.Sprintf("UploadComplete: %s", c.Audiobook.String())
}