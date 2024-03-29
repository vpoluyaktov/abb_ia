package dto

import (
	"fmt"

	"github.com/vpoluyaktov/tview"
)

type AddPageCommand struct {
	Name string
	Grid *tview.Grid
}

func (c *AddPageCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}

type RemovePageCommand struct {
	Name string
}

func (c *RemovePageCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}

type ShowPageCommand struct {
	Name string
}

func (c *ShowPageCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}

type SwitchToPageCommand struct {
	Name string
}

func (c *SwitchToPageCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}
