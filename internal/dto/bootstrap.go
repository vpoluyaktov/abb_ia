package dto

import "fmt"

type BootstrapCommand struct {
}

func (c *BootstrapCommand) String() string {
	return fmt.Sprintf("BootstrapCommand: bootstrap")
}

type FFMPEGNotFoundError struct {
}

func (c *FFMPEGNotFoundError) String() string {
	return fmt.Sprintf("FFMPEGNotFoundError: FFMPEG not found")
}

type NewAppVersionFound struct {
	CurrentVersion string
	NewVersion string
}

func (c *NewAppVersionFound) String() string {
	return fmt.Sprintf("NewAppVersionFound: %s", c.NewVersion)
}