package dto

import (
	"fmt"

	"github.com/rivo/tview"
)

const AddPanelCommandType = "dto.AddPanelCommand"

type AddPanelCommand struct {
	Name string
	Grid *tview.Grid
}

func (c *AddPanelCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}

const RemovePanelCommandType = "dto.RemovePanelCommand"

type RemovePanelCommand struct {
	Name string
}

func (c *RemovePanelCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}

const ShowPanelCommandType = "dto.ShowPanelCommand"

type ShowPanelCommand struct {
	Name string
}

func (c *ShowPanelCommand) String() string {
	return fmt.Sprintf("%T: %s", c, c.Name)
}
