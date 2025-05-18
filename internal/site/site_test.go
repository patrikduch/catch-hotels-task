package site

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	fmt.Printf("Running tests from file: %s\n", filepath.Base(fileName))
}

// Tests for Site structure
func TestSite_UpdateStats(t *testing.T) {
	fmt.Println("Executing test suite: TestSite_UpdateStats")
	t.Parallel()

	site := NewSite("https://example.com")

	// First successful request
	site.UpdateStats(100*time.Millisecond, 1000, true)
	assert.Equal(t, 100*time.Millisecond, site.MinDuration)
	assert.Equal(t, 100*time.Millisecond, site.AvgDuration)
	assert.Equal(t, 100*time.Millisecond, site.MaxDuration)
	assert.Equal(t, int64(1000), site.MinSize)
	assert.Equal(t, int64(1000), site.AvgSize)
	assert.Equal(t, int64(1000), site.MaxSize)
	assert.Equal(t, 1, site.SuccessCount)
	assert.Equal(t, 1, site.TotalCount)

	// Second successful request (slower and smaller)
	site.UpdateStats(200*time.Millisecond, 500, true)
	assert.Equal(t, 100*time.Millisecond, site.MinDuration)
	assert.Equal(t, 150*time.Millisecond, site.AvgDuration)
	assert.Equal(t, 200*time.Millisecond, site.MaxDuration)
	assert.Equal(t, int64(500), site.MinSize)
	assert.Equal(t, int64(750), site.AvgSize)
	assert.Equal(t, int64(1000), site.MaxSize)
	assert.Equal(t, 2, site.SuccessCount)
	assert.Equal(t, 2, site.TotalCount)

	// Third request (failed) — should NOT affect duration/size stats
	site.UpdateStats(50*time.Millisecond, 0, false)
	assert.Equal(t, 100*time.Millisecond, site.MinDuration)
	assert.Equal(t, 150*time.Millisecond, site.AvgDuration)
	assert.Equal(t, 200*time.Millisecond, site.MaxDuration)
	assert.Equal(t, int64(500), site.MinSize)
	assert.Equal(t, int64(750), site.AvgSize)
	assert.Equal(t, int64(1000), site.MaxSize)
	assert.Equal(t, 2, site.SuccessCount)
	assert.Equal(t, 3, site.TotalCount)
}

func TestSite_IsReadyForRequest(t *testing.T) {
	fmt.Println("Executing test suite: TestSite_IsReadyForRequest")
	t.Parallel()

	site := NewSite("https://example.com")

	// First request is always ready
	assert.True(t, site.IsReadyForRequest())

	// Simulate a request
	site.UpdateStats(100*time.Millisecond, 1000, true)

	// Should not be ready (less than 5s)
	assert.False(t, site.IsReadyForRequest())

	// Manually set lastRequest 6s ago
	site.mu.Lock()
	site.lastRequest = time.Now().Add(-6 * time.Second)
	site.mu.Unlock()

	assert.True(t, site.IsReadyForRequest())
}

func TestSite_GetStats(t *testing.T) {
	fmt.Println("Executing test suite: TestSite_GetStats")
	t.Parallel()

	site := NewSite("https://example.com")

	// Before any requests
	stats := site.GetStats()
	assert.Equal(t, "https://example.com", stats[0])
	assert.Equal(t, "n/a", stats[1]) // MinDuration
	assert.Equal(t, "n/a", stats[2]) // AvgDuration
	assert.Equal(t, "n/a", stats[3]) // MaxDuration
	assert.Equal(t, "n/a", stats[4]) // MinSize
	assert.Equal(t, "n/a", stats[5]) // AvgSize
	assert.Equal(t, "n/a", stats[6]) // MaxSize
	assert.Equal(t, "0/0", stats[7]) // OK

	// After one success
	site.UpdateStats(100*time.Millisecond, 1000, true)
	stats = site.GetStats()
	assert.Equal(t, "https://example.com", stats[0])
	assert.Equal(t, "100.00 ms", stats[1])
	assert.Equal(t, "100.00 ms", stats[2])
	assert.Equal(t, "100.00 ms", stats[3])
	assert.Equal(t, "1000 B", stats[4])
	assert.Equal(t, "1000 B", stats[5])
	assert.Equal(t, "1000 B", stats[6])
	assert.Equal(t, "1/1", stats[7])
}
