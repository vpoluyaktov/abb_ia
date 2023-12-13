package dto

import "fmt"

type CleanupCommand struct {
	Audiobook *Audiobook
}

func (c *CleanupCommand) String() string {
	return fmt.Sprintf("CleanupCommand: %s", c.Audiobook.String())
}

type CleanupComplete struct {
	Audiobook *Audiobook
}

func (c *CleanupComplete) String() string {
	return fmt.Sprintf("CleanupComplete: %s", c.Audiobook.String())
}
