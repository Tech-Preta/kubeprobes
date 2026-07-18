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
configured by scanning for missing liveness, readiness, and startup probes.
Proper probe configuration is essential for reliable deployments, effective load
balancing, and early detection of application issues.

Health check probes are critical for:
  • Liveness probes: Detect when to restart containers
  • Readiness probes: Control traffic routing to healthy containers  
  • Startup probes: Handle slow-starting containers gracefully

POSIX Compliance:
This CLI follows POSIX command line conventions including grouped short flags,
flexible flag ordering, and both --flag=value and --flag value syntax.`,
	Example: `  # Quick scan of default namespace
  kubeprobes scan

  # Scan with detailed recommendations
  kubeprobes scan --recommendation

  # POSIX syntax: grouped short flags
  kubeprobes scan -rp liveness

  # POSIX syntax: flexible flag order
  kubeprobes scan -r --namespace test -p readiness

  # Scan specific namespace for liveness probes only
  kubeprobes scan --namespace my-app --probe-type liveness

  # Scan using specific kubeconfig
  kubeprobes scan --kubeconfig ~/.kube/prod-config

  # Check tool version
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
