package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	filePath := "result.txt"
	if len(os.Args) > 1 {
		fmt.Println("Using file:", os.Args[1])
		filePath = os.Args[1]
	}
	iterations := 10000
	expectedWrites := makeExpectedValues("Write", iterations)
	expectedReads := makeExpectedValues("Read", iterations)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matchesExpected(line, &expectedWrites)
		matchesExpected(line, &expectedReads)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Check if all expected writes and reads were found
	if !allFound(expectedWrites) {
		log.Println("Not all Write values found")
	} else {
		fmt.Println("All Write values found")
	}
	if !allFound(expectedReads) {
		log.Println("Not all Read values found")
	} else {
		fmt.Println("All Read values found")
	}
}

// makeExpectedValues creates a map of expected values for a given prefix and count.
func makeExpectedValues(prefix string, count int) map[string]bool {
	expected := make(map[string]bool)
	for i := 0; i < count; i++ {
		expected[fmt.Sprintf("%s %d", prefix, i)] = false
	}
	return expected
}

// matchesExpected checks if the line matches any expected value in the map.
func matchesExpected(line string, expected *map[string]bool) bool {
	if line == "" {
		return true
	}
	for expectedVal := range *expected {
		if strings.HasSuffix(line, expectedVal) {
			(*expected)[expectedVal] = true
			return true
		}
	}
	return false
}

// allFound checks if all expected values in the map have been marked as found.
func allFound(expected map[string]bool) bool {
	for val, found := range expected {
		if !found {
			fmt.Println("Expected value not found:", val)
			return false
		}
	}
	return true
}
