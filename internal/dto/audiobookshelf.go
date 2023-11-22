package dto

import "fmt"

type AudiobookshelfScanCommand struct {
	Audiobook *Audiobook
}

func (c *AudiobookshelfScanCommand) String() string {
	return fmt.Sprintf("AudiobookshelfScanCommand: %s", c.Audiobook.String())
}

type ScanComplete struct {
	Audiobook *Audiobook
}

func (c *ScanComplete) String() string {
	return fmt.Sprintf("ScanComplete: %s", c.Audiobook.String())
}

type AudiobookshelfUploadCommand struct {
	Audiobook *Audiobook
}

func (c *AudiobookshelfUploadCommand) String() string {
	return fmt.Sprintf("AudiobookshelfUploadCommand: %s", c.Audiobook.String())
}

type UploadFileProgress struct {
	FileId   int
	FileName string
	Percent  int
}

func (c *UploadFileProgress) String() string {
	return fmt.Sprintf("UploadFileProgress: %d, %s, %d", c.FileId, c.FileName, c.Percent)
}

type UploadProgress struct {
	Elapsed string // time since started
	Percent int
	Files   string // files encoded
	Bytes   string // total bytes copied
	Speed   string // encode speed bytes/s
	ETA     string // ETA in seconds
}

func (c *UploadProgress) String() string {
	return fmt.Sprintf("UploadProgress: %d", c.Percent)
}

type UploadComplete struct {
	Audiobook *Audiobook
}

func (c *UploadComplete) String() string {
	return fmt.Sprintf("UploadComplete: %s", c.Audiobook.String())
}
