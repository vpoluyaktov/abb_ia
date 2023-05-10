package dto

import (
	"fmt"
)

const StopCommandType = "dto.StopCommand"

type StopCommand struct {
	Process string
}

func (c *StopCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Process)
}
