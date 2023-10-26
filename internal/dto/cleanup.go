package dto

import "fmt"

type CleanupCommand struct {
	Audiobook *Audiobook
}

func (c *CleanupCommand) String() string {
	return fmt.Sprintf("CleanupCommand: %s", c.Audiobook.String())
}