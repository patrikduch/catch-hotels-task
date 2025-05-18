package main

import (
	"os"
	"os/signal"
	"syscall"

	"catch-hotels-task/internal/monitor"
	"catch-hotels-task/internal/ui"
)

func main() {
	// Process command-line arguments
	if len(os.Args) < 2 {
		ui.PrintUsage()
		os.Exit(1)
	}

	urls := os.Args[1:]

	// Validate URLs
	if err := ui.ValidateURLs(urls); err != nil {
		ui.PrintError(err)
		return
	}

	// Create and start monitor
	mon := monitor.New(urls)
	mon.Start()

	// Capture Ctrl+C signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Graceful shutdown
	ui.PrintShutdownMessage()
	mon.Stop()
	ui.PrintShutdownComplete()
}
