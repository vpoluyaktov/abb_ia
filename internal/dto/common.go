package dto

import (
	"fmt"

	"github.com/rivo/tview"
)

type Dto interface {
	String() string
}

const DrawCommandType = "dto.DrawCommand"

type DrawCommand struct {
	Primitive tview.Primitive
}

func (c *DrawCommand) String() string {
	if c.Primitive == nil {
		return fmt.Sprintf("%T", c)
	} else {
		return fmt.Sprintf("%T: %T", c, c.Primitive)
	}
}

const SetFocusCommandType = "dto.SetFocusCommand"

type SetFocusCommand struct {
	Primitive tview.Primitive
}

func (c *SetFocusCommand) String() string {
	return fmt.Sprintf("%T: %T", c, c.Primitive)
}
