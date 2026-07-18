package cli

import (
	"kubeprobes/internal/scanner"

	"github.com/spf13/cobra"
)

// NewScanCommand creates the scan command
func NewScanCommand() *cobra.Command {
	return scanner.NewScanCommand()
}
