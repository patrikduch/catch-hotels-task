package ui

import "fmt"

// PrintUsage prints the usage instructions
func PrintUsage() {
	fmt.Println("Usage: go run main.go [URL1] [URL2] ...")
}

// PrintError prints an error message
func PrintError(err error) {
	fmt.Println("Error:", err)
}

// PrintShutdownMessage prints a message during shutdown
func PrintShutdownMessage() {
	fmt.Println("\nShutting down... Waiting for running requests to complete.")
}

// PrintShutdownComplete prints a message when shutdown is complete
func PrintShutdownComplete() {
	fmt.Println("\nApplication has been successfully shut down.")
}
