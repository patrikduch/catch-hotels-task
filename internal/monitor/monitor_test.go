package monitor

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
	"time"
	"sync"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	fmt.Printf("Running tests from file: %s\n", filepath.Base(fileName))
}

func TestMonitor_SuccessfulRequest(t *testing.T) {
	fmt.Println("Executing test: TestMonitor_SuccessfulRequest")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://example.com",
		httpmock.NewStringResponder(200, "Hello, World!"))

	mon := New([]string{"https://example.com"})
	go mon.Start()
	time.Sleep(1 * time.Second)
	mon.Stop()

	summaries := mon.ToSummaries()
	assert.Len(t, summaries, 1)
	assert.Equal(t, "https://example.com", summaries[0][0])
	assert.Equal(t, "1/1", summaries[0][7])
}

func TestMonitor_FailedRequest(t *testing.T) {
	fmt.Println("Executing test: TestMonitor_FailedRequest")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://error.com",
		httpmock.NewStringResponder(500, "Internal Server Error"))

	mon := New([]string{"https://error.com"})
	go mon.Start()
	time.Sleep(1 * time.Second)
	mon.Stop()

	summaries := mon.ToSummaries()
	assert.Len(t, summaries, 1)
	assert.Equal(t, "https://error.com", summaries[0][0])
	assert.Equal(t, "0/1", summaries[0][7])
}

func TestMonitor_MultipleRequests(t *testing.T) {
	fmt.Println("Executing test: TestMonitor_MultipleRequests")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://example.com",
		httpmock.NewStringResponder(200, "Hello again!"))

	mon := New([]string{"https://example.com"})
	go mon.Start()

	time.Sleep(12 * time.Second)
	mon.Stop()

	summaries := mon.ToSummaries()
	assert.Len(t, summaries, 1)
	assert.Equal(t, "https://example.com", summaries[0][0])

	var success, total int
	fmt.Sscanf(summaries[0][7], "%d/%d", &success, &total)

	assert.GreaterOrEqual(t, total, 2)
	assert.GreaterOrEqual(t, success, 2)
}

func TestMonitor_GracefulShutdown(t *testing.T) {
	fmt.Println("Executing test: TestMonitor_GracefulShutdown")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var mu sync.Mutex
	requestStarted := false
	requestFinished := make(chan struct{})

	// Responder that simulates a slow request but is NOT affected by context cancellation
	httpmock.RegisterResponder("GET", "https://slow.example.com",
		func(req *http.Request) (*http.Response, error) {
			mu.Lock()
			requestStarted = true
			mu.Unlock()

			// Sleep to simulate long-running request
			time.Sleep(2 * time.Second)

			// Signal request finished
			close(requestFinished)

			return httpmock.NewStringResponse(200, "Slow response"), nil
		})

	mon := New([]string{"https://slow.example.com"})
	for _, s := range mon.sites {
		s.ForceReady()
	}

	go mon.Start()

	// Wait until request has started
	ok := waitUntilTrue(1*time.Second, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return requestStarted
	})
	if !ok {
		t.Fatal("Request did not start in time")
	}

	// Wait for the request to finish
	select {
	case <-requestFinished:
		// Proceed to stop only after request finished
	case <-time.After(3 * time.Second):
		t.Fatal("Request did not finish in time")
	}

	// Now call Stop (ctx gets cancelled *after* request has finished)
	start := time.Now()
	mon.Stop()
	elapsed := time.Since(start)

	summaries := mon.ToSummaries()
	assert.Equal(t, "https://slow.example.com", summaries[0][0])
	assert.Equal(t, "1/1", summaries[0][7])
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(100), "Shutdown was not graceful")
}

func waitUntilTrue(timeout time.Duration, condition func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}