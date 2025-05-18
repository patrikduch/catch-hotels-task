// internal/ui/utils.go

package ui

import (
	"fmt"
	"net/http"
)

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
