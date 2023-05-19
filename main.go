package main

import (
	"flag"
	"os"

	"github.com/vpoluyaktov/audiobook_creator_IA/cmd"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/config"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

func main() {
	config.Load()

	// command line arguments
	logLevel := flag.String("log-level", "INFO", "Logging level")
	useMock := flag.Bool("mock-load", false, "Use mock data")
	saveMock := flag.Bool("mock-save", false, "Save mock data")
	help := flag.Bool("help", false, "Display usage information")
	flag.Parse()

	// get IA search condition from command line if specified
	searchCondition := flag.Arg(0)

	// display usage information
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// save runtime configuration
	config.SetSearchCondition(searchCondition)
	config.SetLogLevel(*logLevel)
	config.UseMock(*useMock)
	config.SaveMock(*saveMock)

	logger.Init(config.LogFileName(), *logLevel)
	cmd.Execute()
}
