package dto

import "fmt"

type ChaptersEditCommand struct {
	Audiobook *Audiobook
}

func (c *ChaptersEditCommand) String() string {
	return fmt.Sprintf("ChaptersEditCommand: %s", c.Audiobook.String())
}

type ChaptersEditComplete struct {
	Audiobook *Audiobook
}

func (c *ChaptersEditComplete) String() string {
	return fmt.Sprintf("ChaptersEditComplete: %s", c.Audiobook.String())
}
