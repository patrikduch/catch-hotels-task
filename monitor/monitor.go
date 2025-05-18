package monitor

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"

    "catch-hotels-task/internal/site"
    "catch-hotels-task/internal/ui"
)

// ValidateURLs checks if URLs are valid
func ValidateURLs(urls []string) error {
    if len(urls) == 0 {
        return fmt.Errorf("you must provide at least one URL")
    }

    for _, url := range urls {
        _, err := http.NewRequest("GET", url, nil)
        if err != nil {
            return fmt.Errorf("invalid URL address: %s - %v", url, err)
        }
    }

    return nil
}

// MonitorSite monitors a given URL until it receives a signal to stop
func MonitorSite(ctx context.Context, site *site.Site, wg *sync.WaitGroup) {
    defer wg.Done()

    // Create an HTTP client with timeout
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Track if a request is in progress
    var inProgress bool
    var requestWg sync.WaitGroup

    for {
        // Check if monitoring should stop
        select {
        case <-ctx.Done():
            // Wait for any ongoing request to complete
            if inProgress {
                requestWg.Wait()
            }
            return
        default:
            // Check if it's time for the next request
            if !site.IsReadyForRequest() || inProgress {
                // Check for termination while waiting
                select {
                case <-ctx.Done():
                    if inProgress {
                        requestWg.Wait()
                    }
                    return
                case <-time.After(100 * time.Millisecond):
                    continue
                }
            }

            // Mark that a request is starting
            inProgress = true
            requestWg.Add(1)

            // Run request in a goroutine
            go func() {
                defer requestWg.Done()
                defer func() { inProgress = false }()

                // Create a new context for this request
                reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
                defer cancel()

                // Send HTTP GET request
                startTime := time.Now()
                req, err := http.NewRequestWithContext(reqCtx, "GET", site.URL, nil)
                if err != nil {
                    site.UpdateStats(time.Since(startTime), -1, false)
                    return
                }

                resp, err := client.Do(req)
                if err != nil {
                    site.UpdateStats(time.Since(startTime), -1, false)
                    return
                }

                // Read and measure response size
                body, err := io.ReadAll(resp.Body)
                resp.Body.Close() 
                
                duration := time.Since(startTime)
                size := int64(len(body))
                
                // Check if request was successful (HTTP 2xx or 3xx)
                success := resp.StatusCode >= 200 && resp.StatusCode < 400
                
                // Update statistics
                site.UpdateStats(duration, size, success)
            }()

            // Small delay to let the request start
            time.Sleep(10 * time.Millisecond)
        }
    }
}

// Start monitoring all URLs and display statistics
func Start(urls []string) (context.Context, context.CancelFunc, *sync.WaitGroup) {
    // Create sites for each URL
    sites := make([]*site.Site, len(urls))
    for i, url := range urls {
        sites[i] = site.NewSite(url)
    }

    // Create context for controlling goroutines
    ctx, cancel := context.WithCancel(context.Background())
    
    // Create WaitGroup for monitoring goroutines
    var wg sync.WaitGroup
    
    // Start monitoring for each URL
    for _, site := range sites {
        wg.Add(1)
        go MonitorSite(ctx, site, &wg)
    }

    // Start periodic display updates
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                // Collect data for display
                data := make([][]string, len(sites))
                for i, site := range sites {
                    data[i] = site.GetStats()
                }
                
                // Display table
                ui.DisplayTable(data, true)
            case <-ctx.Done():
                return
            }
        }
    }()

    return ctx, cancel, &wg
}