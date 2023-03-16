package main

import (
  "github.com/vpoluyaktov/audiobook_creator_IA/cmd"
  "github.com/vpoluyaktov/audiobook_creator_IA/internal/config"
  "github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

var (
  logFileName string = "/tmp/audiobook_creator_IA.log"
  logLevel = "DEBUG"
)

func main() {
  logger.Init(logFileName, logLevel)
  config.Load()
  cmd.Execute()
}
