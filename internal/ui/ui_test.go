package ui

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	fmt.Printf("Running tests from file: %s\n", filepath.Base(fileName))
}

func TestValidateURLs(t *testing.T) {
	fmt.Println("Executing test suite: TestValidateURLs")

	err := ValidateURLs([]string{"https://example.com", "http://localhost:8080"})
	assert.NoError(t, err)

	err = ValidateURLs([]string{":::invalid-url:::"})
	assert.Error(t, err, "Expected error for clearly invalid URL")
}

func TestDisplayTable(t *testing.T) {
	fmt.Println("Executing test suite: TestDisplayTable")
	t.Parallel()

	assert.NotPanics(t, func() {
		DisplayTable([][]string{}, true)
	})

	mockData := [][]string{
		{
			"https://example.com",
			"100ms", "150ms", "200ms",
			"1KB", "1.5KB", "2KB", "1/1",
		},
	}

	assert.NotPanics(t, func() {
		DisplayTable(mockData, false)
	})
}
