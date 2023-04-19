package dto

import (
	"fmt"

	"github.com/rivo/tview"
)

const AddPageCommandType = "dto.AddPageCommand"

type AddPageCommand struct {
	Name string
	Grid *tview.Grid
}

func (c *AddPageCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}

const RemovePageCommandType = "dto.RemovePageCommand"

type RemovePageCommand struct {
	Name string
}

func (c *RemovePageCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}

const ShowPageCommandType = "dto.ShowPageCommand"

type ShowPageCommand struct {
	Name string
}

func (c *ShowPageCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}

const SwitchToPageCommandType = "dto.SwitchToPageCommand"

type SwitchToPageCommand struct {
	Name string
}

func (c *SwitchToPageCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}
