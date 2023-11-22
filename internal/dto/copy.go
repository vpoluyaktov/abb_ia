package dto

import "fmt"

type CopyCommand struct {
	Audiobook *Audiobook
}

func (c *CopyCommand) String() string {
	return fmt.Sprintf("CopyCommand: %s", c.Audiobook.String())
}

type CopyFileProgress struct {
	FileId   int
	FileName string
	Percent  int
}

func (c *CopyFileProgress) String() string {
	return fmt.Sprintf("CopyFileProgress: %d, %s, %d", c.FileId, c.FileName, c.Percent)
}

type CopyProgress struct {
	Elapsed string // time since started
	Percent int
	Files   string // files encoded
	Bytes   string // total bytes copied
	Speed   string // encode speed bytes/s
	ETA     string // ETA in seconds
}

func (c *CopyProgress) String() string {
	return fmt.Sprintf("CopyProgress: %d", c.Percent)
}

type CopyComplete struct {
	Audiobook *Audiobook
}

func (c *CopyComplete) String() string {
	return fmt.Sprintf("CopyComplete: %s", c.Audiobook.String())
}
