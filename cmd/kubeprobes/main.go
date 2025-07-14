package main

import (
	"log"
	"os"

	"kubeprobes/internal/cli"
	"kubeprobes/internal/scanner"
)

func main() {
	if err := cli.Execute(); err != nil {
		// Check if it's our custom error indicating probe issues found
		if _, ok := err.(*scanner.ProbeIssuesFoundError); ok {
			os.Exit(1)
		}
		log.Fatalf("Error executing command: %s", err.Error())
	}
}
