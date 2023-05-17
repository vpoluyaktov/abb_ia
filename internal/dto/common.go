package dto

import (
	"fmt"
)

type StopCommand struct {
	Process string
	Reason string
}

func (c *StopCommand) String() string {
	return fmt.Sprintf("%T: Process: %s, Reason: %s", c, c.Process, c.Reason)
}
