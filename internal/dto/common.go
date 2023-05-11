package dto

import (
	"fmt"
)

type StopCommand struct {
	Process string
}

func (c *StopCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Process)
}
