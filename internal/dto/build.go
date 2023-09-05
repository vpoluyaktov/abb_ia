package dto

import "fmt"

type BuildCommand struct {
	Audiobook *Audiobook
}

func (c *BuildCommand) String() string {
	return fmt.Sprintf("BuildCommand: %s", c.Audiobook.String())
}


type BuildProgress struct {
	FileId   int
	FileName string
	Percent  int
}

func (c *BuildProgress) String() string {
	return fmt.Sprintf("BuildProgress: %d, %s, %d", c.FileId, c.FileName, c.Percent)
}

type BuildComplete struct {
	Audiobook *Audiobook
}

func (c *BuildComplete) String() string {
	return fmt.Sprintf("BuildComplete: %s", c.Audiobook.String())
}
