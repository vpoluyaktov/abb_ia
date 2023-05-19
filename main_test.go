package main_test

import (
	"os"
	"testing"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/config"
  "github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

func TestMain(m *testing.M) {
  config.Load()
  logger.Init(config.LogFileName(), "DEBUG")
	os.Exit(m.Run())
}