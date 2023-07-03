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

type DownloadFileProgress struct {
	FileId   int
	FileName string
	Percent  int
}

func (c *DownloadFileProgress) String() string {
	return fmt.Sprintf("DownloadFileProgress: %d, %s, %d", c.FileId, c.FileName, c.Percent)
}

type DownloadProgress struct {
	Elapsed string // time since started
	Percent int
	Files   string // files downloaded
	Bytes   string // total bytes downloaded
	Speed   string // download speed bytes/s
	ETA     string // ETA in seconds
}

func (c *DownloadProgress) String() string {
	return fmt.Sprintf("DownloadProgress: %d", c.Percent)
}

type DownloadComplete struct {
	Audiobook *Audiobook
}

func (c *DownloadComplete) String() string {
	return fmt.Sprintf("DownloadComplete: %s", c.Audiobook.String())
}
