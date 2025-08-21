package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	matches, err := filepath.Glob("**/Dockerfile")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Matches: %v\n", matches)
}