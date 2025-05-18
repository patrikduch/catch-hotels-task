package monitor

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"catch-hotels-task/internal/site"
	"catch-hotels-task/internal/ui"
)

type Monitor struct {
	urls     []string
	sites    map[string]*site.Site
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	interval time.Duration
	timeout  time.Duration
}

func New(urls []string) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	sites := make(map[string]*site.Site)
	for _, url := range urls {
		sites[url] = site.NewSite(url)
	}
	return &Monitor{
		urls:     urls,
		sites:    sites,
		ctx:      ctx,
		cancel:   cancel,
		interval: 5 * time.Second,
		timeout:  10 * time.Second,
	}
}

func (m *Monitor) Start() {
	for _, url := range m.urls {
		m.wg.Add(1)
		go m.worker(url)
	}
	go m.renderLoop()
}

func (m *Monitor) Stop() {
	m.cancel()
	m.wg.Wait()
	ui.DisplayTable(m.ToSummaries(), false)
}

func (m *Monitor) worker(url string) {
	defer m.wg.Done()
	client := http.Client{
		Timeout: m.timeout,
	}

	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			// Check if it's time for the next request
			if !m.sites[url].IsReadyForRequest() {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Send HTTP request
			start := time.Now()
			req, err := http.NewRequestWithContext(m.ctx, "GET", url, nil)
			if err != nil {
				m.sites[url].UpdateStats(time.Since(start), -1, false)
				continue
			}

			resp, err := client.Do(req)
			if err != nil {
				m.sites[url].UpdateStats(time.Since(start), -1, false)
				continue
			}

			// Process response
			var size int64 = -1
			if resp.Body != nil {
				body, _ := io.ReadAll(resp.Body)
				size = int64(len(body))
				resp.Body.Close()
			}

			duration := time.Since(start)
			success := resp.StatusCode >= 200 && resp.StatusCode < 400

			// Update statistics
			m.sites[url].UpdateStats(duration, size, success)
		}
	}
}

func (m *Monitor) renderLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			ui.DisplayTable(m.ToSummaries(), true)
		}
	}
}

func (m *Monitor) ToSummaries() [][]string {
	var summaries [][]string

	for _, url := range m.urls {
		stats := m.sites[url].GetStats()
		summaries = append(summaries, stats)
	}
	return summaries
}