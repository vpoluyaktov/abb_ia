package main_test

import (
	"os"
	"testing"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/config"
  "github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

const (
  logFileName string = "/tmp/audiobook_creator_IA.test.log"
  logLevel = logger.DEBUG
)

func TestMain(m *testing.M) {
  logger.Init(logFileName, logLevel)
  config.Load()
	os.Exit(m.Run())
}