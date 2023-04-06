package main

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/cmd"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/config"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

const (
	logFileName string = "/tmp/audiobook_creator_IA.log"
	logLevel           = logger.DEBUG
)

func main() {
	logger.Init(logFileName, logLevel)
	config.Load()
	cmd.Execute()
}
