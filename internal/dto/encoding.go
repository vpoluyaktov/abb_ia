package dto

import "fmt"

type EncodeCommand struct {
	Audiobook *Audiobook
}

func (c *EncodeCommand) String() string {
	return fmt.Sprintf("EncodeCommand: %s", c.Audiobook.String())
}

type EncodingFileProgress struct {
	FileId   int
	FileName string
	Percent  int
}

func (c *EncodingFileProgress) String() string {
	return fmt.Sprintf("EncodingFileProgress: %d, %s, %d", c.FileId, c.FileName, c.Percent)
}

type EncodingProgress struct {
	Duration string // time since started
	Percent  int
	Files    string // files downloaded
	Bytes    string // total bytes downloaded
	Speed    string // download speed bytes/s
	ETA      string // ETA in seconds
}

func (c *EncodingProgress) String() string {
	return fmt.Sprintf("EncodingProgress: %d", c.Percent)
}
