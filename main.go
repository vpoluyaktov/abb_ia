package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"abb_ia/cmd"
	"abb_ia/internal/config"
	"abb_ia/internal/logger"
	"abb_ia/internal/monitoring"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"
)

// Min screen size for comfortable layout 45x125 characters
func main() {
	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Shutdown signal received")
		cancel()
	}()

	// Initialize configuration
	config.Load()

	// command line arguments
	logLevel := flag.String("log-level", "INFO", "Logging level")
	useMock := flag.Bool("mock-load", false, "Use mock data")
	saveMock := flag.Bool("mock-save", false, "Save mock data")
	help := flag.Bool("help", false, "Display usage information")
	enableMetrics := flag.Bool("enable-metrics", false, "Enable metrics collection (disabled by default)")
	flag.Parse()

	// Initialize logger
	logger.Init(config.Instance().GetLogFileName(), config.Instance().GetLogLevel())

	// Initialize metrics collector
	if *enableMetrics {
		monitoring.EnableMetrics()
		logger.Info("Metrics collection enabled")
	} else {
		logger.Info("Metrics collection disabled (default)")
	}
	metricsCollector := monitoring.GetMetricsCollector()
	if *enableMetrics {
		metricsCollector.StartMetricsReporter(1 * time.Minute)
	}

	// get IA search condition from command line if specified
	searchCondition := flag.Arg(0)

	// display usage information
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Create message dispatcher
	dispatcher := mq.NewDispatcher()
	defer func() {
		if err := dispatcher.Shutdown(5 * time.Second); err != nil {
			logger.Error(fmt.Sprintf("Error shutting down dispatcher: %v", err))
		}
	}()

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

	// Create circuit breakers for external services if needed
	// Example:
	// archiveOrgCB := utils.NewCircuitBreaker("archive.org",
	//	utils.WithFailureThreshold(3),
	//	utils.WithResetTimeout(30*time.Second),
	// )

	// Start the application
	startTime := time.Now()
	metricsCollector.SetGauge("app_uptime_seconds", 0, monitoring.Labels{
		"version": config.Instance().AppVersion(),
	})

	// Update uptime metric
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metricsCollector.SetGauge("app_uptime_seconds",
					time.Since(startTime).Seconds(),
					monitoring.Labels{"version": config.Instance().AppVersion()})
			}
		}
	}()

	cmd.Execute()
}
