package ui

import (
	"fmt"
	"strings"
)

// DisplayTable renders a table with the current statistics
func DisplayTable(data [][]string, clearScreen bool) {
	if clearScreen {
		fmt.Print("\033[H\033[2J") // ANSI clear screen
	}

	headers := []string{
		"URL",
		"Min Duration",
		"Avg Duration",
		"Max Duration",
		"Min Size",
		"Avg Size",
		"Max Size",
		"OK",
	}

	// Determine column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	for _, row := range data {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print separator
	separator := "+"
	for _, width := range colWidths {
		separator += strings.Repeat("-", width+2) + "+"
	}

	// Print header
	fmt.Println(separator)
	fmt.Print("|")
	for i, header := range headers {
		fmt.Printf(" %-*s |", colWidths[i], header)
	}
	fmt.Println()
	fmt.Println(separator)

	// Print rows
	for _, row := range data {
		fmt.Print("|")
		for i, cell := range row {
			fmt.Printf(" %-*s |", colWidths[i], cell)
		}
		fmt.Println()
	}
	fmt.Println(separator)
}
