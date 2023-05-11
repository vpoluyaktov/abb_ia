package dto

import "fmt"

type UpdateStatus struct {
	Message string
}

func (c *UpdateStatus) String() string {
	return fmt.Sprintf("%T: %s", c, c.Message)
}

type SetBusyIndicator struct {
	Busy bool
}

func (c *SetBusyIndicator) String() string {
	busyStr := ""
	if c.Busy {
		busyStr = "Busy"
	} else {
		busyStr = "Unbusy"
	}
	return fmt.Sprintf("%T: %s", c, busyStr)
}

