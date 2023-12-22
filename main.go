package main

import (
	"flag"
	"os"
	"strings"

	"abb_ia/cmd"
	"abb_ia/internal/config"
	"abb_ia/internal/logger"
	"abb_ia/internal/utils"
)

// Min screen size for comfortable layout 45x125 characters
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
	if searchCondition != "" {
		condition := strings.Split(searchCondition, " - ")
		if len(condition) == 2 {
			config.Instance().SetDefaultAuthor(condition[0])
			config.Instance().SetDefaultTitle(condition[1])
		}
	}
	if utils.IsFlagPassed("log-level") {
		config.Instance().SetLogLevel(*logLevel)
	}
	if utils.IsFlagPassed("mock-load") {
		config.Instance().SetUseMock(*useMock)
	}
	if utils.IsFlagPassed("mock-save") {
		config.Instance().SetSaveMock(*saveMock)
	}

	logger.Init(config.Instance().GetLogFileName(), config.Instance().GetLogLevel())
	cmd.Execute()
}
