package main_test

import (
	"os"
	"testing"

	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
)

func TestMain(m *testing.M) {
	config.Load()
	logger.Init(config.LogFileName(), "DEBUG")
	os.Exit(m.Run())
}
