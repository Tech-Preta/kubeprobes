package cli

import (
	"kubeprobes/internal/scanner"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubeprobes",
	Short: "Kubeprobes is a CLI tool for scanning Kubernetes probes",
	Long: `Kubeprobes is a CLI tool for scanning Kubernetes workloads to detect 
missing liveness, readiness, and startup probes.`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(scanner.NewScanCommand())
}
