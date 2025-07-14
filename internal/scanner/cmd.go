package scanner

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// NewScanCommand creates the scan command
func NewScanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan Kubernetes workloads for probes",
		Long: `Scan Kubernetes workloads for missing liveness, readiness, or startup probes.

This command connects to your Kubernetes cluster and examines pods to identify
containers that are missing health check probes. Proper probe configuration is
essential for reliable application deployments and effective health monitoring.

Exit codes:
  0: No probe issues found
  1: Probe issues detected`,
		Example: `  # Scan all workloads in the default namespace
  kubeprobes scan

  # Scan with recommendations for missing probes
  kubeprobes scan --recommendation

  # Scan a specific namespace for liveness probes only
  kubeprobes scan --namespace my-app --probe-type liveness

  # Scan using a specific kubeconfig file
  kubeprobes scan --kubeconfig ~/.kube/config

  # Scan with a specific Kubernetes context
  kubeprobes scan --kubeContext production`,
		RunE: func(cmd *cobra.Command, args []string) error {
			kubeconfig, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return fmt.Errorf("failed to read kubeconfig flag: %w", err)
			}

			kubeContext, err := cmd.Flags().GetString("kubeContext")
			if err != nil {
				return fmt.Errorf("failed to read kubeContext flag: %w", err)
			}

			namespace, err := cmd.Flags().GetString("namespace")
			if err != nil {
				return fmt.Errorf("failed to read namespace flag: %w", err)
			}

			probeType, err := cmd.Flags().GetString("probe-type")
			if err != nil {
				return fmt.Errorf("failed to read probe-type flag: %w", err)
			}

			recommendation, err := cmd.Flags().GetBool("recommendation")
			if err != nil {
				return fmt.Errorf("failed to read recommendation flag: %w", err)
			}

			// Use context with timeout instead of context.TODO()
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			scanner, err := NewProbeScanner(kubeconfig, kubeContext, namespace, probeType, recommendation)
			if err != nil {
				return err
			}

			return scanner.Scan(ctx)
		},
	}

	cmd.Flags().StringP("kubeconfig", "k", "", "Path to the kubeconfig file (defaults to $KUBECONFIG or ~/.kube/config)")
	cmd.Flags().StringP("kubeContext", "c", "", "Kubernetes context to use (defaults to current context)")
	cmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace to scan (default: default)")
	cmd.Flags().StringP("probe-type", "p", "", "Type of probe to scan for: liveness, readiness, or startup (default: all types)")
	cmd.Flags().BoolP("recommendation", "r", false, "Show actionable recommendations for missing probes")

	// Add custom completion for probe-type flag
	err := cmd.RegisterFlagCompletionFunc("probe-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"liveness", "readiness", "startup"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		// Log error but don't fail the command creation
		fmt.Printf("Warning: failed to register completion for probe-type flag: %v\n", err)
	}

	return cmd
}
