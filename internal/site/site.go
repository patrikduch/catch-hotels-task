package site

import (
	"fmt"
	"sync"
	"time"
)

type Site struct {
	URL           string
	MinDuration   time.Duration
	AvgDuration   time.Duration
	MaxDuration   time.Duration
	MinSize       int64
	AvgSize       int64
	MaxSize       int64
	SuccessCount  int
	TotalCount    int
	lastRequest   time.Time
	mu            sync.Mutex
}

// NewSite creates a new Site instance
func NewSite(url string) *Site {
	return &Site{
		URL:         url,
		lastRequest: time.Now().Add(-10 * time.Second), // ensure ready immediately
	}
}

// IsReadyForRequest determines if enough time has passed since last request
func (s *Site) IsReadyForRequest() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return time.Since(s.lastRequest) >= 5*time.Second
}

// UpdateStats updates the stats of the site after a request
func (s *Site) UpdateStats(duration time.Duration, size int64, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TotalCount++
	if success {
		s.SuccessCount++

		if s.MinDuration == 0 || duration < s.MinDuration {
			s.MinDuration = duration
		}
		if duration > s.MaxDuration {
			s.MaxDuration = duration
		}
		s.AvgDuration = ((s.AvgDuration * time.Duration(s.SuccessCount-1)) + duration) / time.Duration(s.SuccessCount)

		if s.MinSize == 0 || size < s.MinSize {
			s.MinSize = size
		}
		if size > s.MaxSize {
			s.MaxSize = size
		}
		s.AvgSize = ((s.AvgSize * int64(s.SuccessCount-1)) + size) / int64(s.SuccessCount)
	}

	s.lastRequest = time.Now()
}

// GetStats returns a slice of string stats for table display
func (s *Site) GetStats() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	formatDuration := func(d time.Duration) string {
		if d == 0 {
			return "n/a"
		}
		return fmt.Sprintf("%.2f ms", float64(d.Milliseconds()))
	}

	formatSize := func(b int64) string {
		if b <= 0 {
			return "n/a"
		}
		return fmt.Sprintf("%d B", b)
	}

	return []string{
		s.URL,
		formatDuration(s.MinDuration),
		formatDuration(s.AvgDuration),
		formatDuration(s.MaxDuration),
		formatSize(s.MinSize),
		formatSize(s.AvgSize),
		formatSize(s.MaxSize),
		fmt.Sprintf("%d/%d", s.SuccessCount, s.TotalCount),
	}
}

// ForceReady sets lastRequest back in time to allow immediate request (used in tests)
func (s *Site) ForceReady() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastRequest = time.Now().Add(-10 * time.Second)
}
