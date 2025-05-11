package main

import (
	"fmt"
	"os"
)

func main() {
	storage := NewStorage("")

	rootCmd := setupCommands(storage)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
