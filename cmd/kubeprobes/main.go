package main

import (
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
		// The error has already been displayed by cobra, so just exit
		os.Exit(1)
	}
}
