package dto

import "fmt"

type CopyCommand struct {
	Audiobook *Audiobook
}

func (c *CopyCommand) String() string {
	return fmt.Sprintf("CopyCommand: %s", c.Audiobook.String())
}


type CopyProgress struct {
	FileId   int
	FileName string
	Percent  int
}

func (c *CopyProgress) String() string {
	return fmt.Sprintf("CopyProgress: %d, %s, %d", c.FileId, c.FileName, c.Percent)
}

type CopyComplete struct {
	Audiobook *Audiobook
}

func (c *CopyComplete) String() string {
	return fmt.Sprintf("CopyComplete: %s", c.Audiobook.String())
}