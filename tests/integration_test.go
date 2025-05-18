package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"catch-hotels-task/internal/monitor"
	"catch-hotels-task/internal/ui"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	fmt.Printf("Running tests from file: %s\n", filepath.Base(fileName))
}

// Integration test - simulates running the entire application
func TestIntegration(t *testing.T) {
	fmt.Println("Executing test suite: TestIntegration")

	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register responders for test URLs
	httpmock.RegisterResponder("GET", "https://example.com",
		httpmock.NewStringResponder(200, "Test response"))

	httpmock.RegisterResponder("GET", "https://slow.com",
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(500 * time.Millisecond) // Simulate slow response
			return httpmock.NewStringResponse(200, "Slow response"), nil
		})

	// Simulate command line arguments
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"cmd", "https://example.com", "https://slow.com"}

	// Run test version of main logic
	urls := os.Args[1:]

	// Validate URLs
	err := ui.ValidateURLs(urls)
	assert.NoError(t, err, "URL validation should pass")

	// Start monitoring
	mon := monitor.New(urls)
	go mon.Start()

	// Let it collect data
	time.Sleep(2 * time.Second)
	mon.Stop()

	// Get summaries
	summaries := mon.ToSummaries()
	assert.Len(t, summaries, 2)

	// Print results
	fmt.Println("Final test results:")
	for i, summary := range summaries {
		fmt.Printf("Site %d (%s): %s OK, Avg Duration: %s\n",
			i+1, summary[0], summary[7], summary[2])
	}

	// Validate httpmock was used
	info := httpmock.GetCallCountInfo()
	assert.GreaterOrEqual(t, info["GET https://example.com"], 1)
	assert.GreaterOrEqual(t, info["GET https://slow.com"], 1)
}
