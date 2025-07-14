package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubeprobes",
	Short: "Kubeprobes is a CLI tool for scanning Kubernetes probes",
	Long: `Kubeprobes is a CLI tool for scanning Kubernetes workloads to detect 
missing liveness, readiness, and startup probes.

Kubeprobes helps you ensure your Kubernetes workloads have proper health checks
configured by scanning for missing liveness, readiness, and startup probes.`,
	Example: `  # Scan all workloads in the default namespace
  kubeprobes scan

  # Scan with recommendations for missing probes
  kubeprobes scan --recommendation

  # Scan a specific namespace for liveness probes only
  kubeprobes scan --namespace my-app --probe-type liveness

  # Show version information
  kubeprobes version`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(NewScanCommand())
	rootCmd.AddCommand(NewVersionCommand())
}
