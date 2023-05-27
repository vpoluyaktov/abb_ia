package dto

import "fmt"

type DownloadCommand struct {
	Audiobook *Audiobook
}

func (c *DownloadCommand) String() string {
	return fmt.Sprintf("DownloadCommand: %s", c.Audiobook.String())
}

type DisplayBookInfoCommand struct {
	Audiobook *Audiobook
}

func (c *DisplayBookInfoCommand) String() string {
	return fmt.Sprintf("DisplayBookInfoCommand: %s", c.Audiobook.String())
}

type DownloadProgress struct {
	FileId   int
	FileName string
	Percent  int
}

func (c *DownloadProgress) String() string {
	return fmt.Sprintf("DownloadProgress: %d, %s, %d", c.FileId, c.FileName, c.Percent)
}
