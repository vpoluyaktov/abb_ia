package dto

import "fmt"

type BuildCommand struct {
	Audiobook *Audiobook
}

func (c *BuildCommand) String() string {
	return fmt.Sprintf("BuildCommand: %s", c.Audiobook.String())
}

type BuildFileProgress struct {
	FileId   int
	FileName string
	Percent  int
}

func (c *BuildFileProgress) String() string {
	return fmt.Sprintf("BuildFileProgress: %d, %s, %d", c.FileId, c.FileName, c.Percent)
}

type BuildProgress struct {
	Elapsed string // time since started
	Percent int
	Files   string // files encoded
	Speed   string // encode speed bytes/s
	ETA     string // ETA in seconds
}

func (c *BuildProgress) String() string {
	return fmt.Sprintf("BuildProgress: %d", c.Percent)
}


type BuildComplete struct {
	Audiobook *Audiobook
}

func (c *BuildComplete) String() string {
	return fmt.Sprintf("BuildComplete: %s", c.Audiobook.String())
}
