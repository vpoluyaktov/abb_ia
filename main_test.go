package main_test

import (
	"os"
	"testing"

	"abb_ia/internal/config"
	"abb_ia/internal/logger"
)

func TestMain(m *testing.M) {
	config.Load()
	logger.Init(config.Instance().GetLogFileName(), "DEBUG")
	os.Exit(m.Run())
}
