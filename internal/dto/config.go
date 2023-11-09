package dto

import (
	"fmt"

	"github.com/vpoluyaktov/abb_ia/internal/config"
)

type DisplayConfigCommand struct {
	Config config.Config
}

func (c *DisplayConfigCommand) String() string {
	return fmt.Sprintf("DisplayConfigCommand: %s", c.Config.GetLogLevel())
}

type SaveConfigCommand struct {
	Config config.Config
}

func (c *SaveConfigCommand) String() string {
	return fmt.Sprintf("SaveConfigCommand: %s", c.Config.GetLogLevel())
}
