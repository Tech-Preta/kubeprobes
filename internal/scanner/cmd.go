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

Exit codes:
  0: Nenhum problema encontrado
  1: Problemas de probe encontrados
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			kubeconfig, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return fmt.Errorf("error getting kubeconfig flag: %w", err)
			}

			kubeContext, err := cmd.Flags().GetString("kubeContext")
			if err != nil {
				return fmt.Errorf("error getting kubeContext flag: %w", err)
			}

			namespace, err := cmd.Flags().GetString("namespace")
			if err != nil {
				return fmt.Errorf("error getting namespace flag: %w", err)
			}

			probeType, err := cmd.Flags().GetString("probe-type")
			if err != nil {
				return fmt.Errorf("error getting probe-type flag: %w", err)
			}

			recommendation, err := cmd.Flags().GetBool("recommendation")
			if err != nil {
				return fmt.Errorf("error getting recommendation flag: %w", err)
			}

			// Use context with timeout instead of context.TODO()
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			scanner, err := NewProbeScanner(kubeconfig, kubeContext, namespace, probeType, recommendation)
			if err != nil {
				return fmt.Errorf("error creating scanner: %w", err)
			}

			return scanner.Scan(ctx)
		},
	}

	cmd.Flags().StringP("kubeconfig", "k", "", "path to the kubeconfig file")
	cmd.Flags().StringP("kubeContext", "c", "", "Kubernetes context")
	cmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace (default: default)")
	cmd.Flags().StringP("probe-type", "p", "", "type of probe to scan for (liveness, readiness, startup)")
	cmd.Flags().BoolP("recommendation", "r", false, "show recommendations for missing probes")

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
