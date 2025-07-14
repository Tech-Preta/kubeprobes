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

			allNamespaces, err := cmd.Flags().GetBool("all-namespaces")
			if err != nil {
				return fmt.Errorf("error getting all-namespaces flag: %w", err)
			}

			probeType, err := cmd.Flags().GetString("probe-type")
			if err != nil {
				return fmt.Errorf("error getting probe-type flag: %w", err)
			}

			recommendation, err := cmd.Flags().GetBool("recommendation")
			if err != nil {
				return fmt.Errorf("error getting recommendation flag: %w", err)
			}

			output, err := cmd.Flags().GetString("output")
			if err != nil {
				return fmt.Errorf("error getting output flag: %w", err)
			}

			failOnWarn, err := cmd.Flags().GetBool("fail-on-warn")
			if err != nil {
				return fmt.Errorf("error getting fail-on-warn flag: %w", err)
			}

			// Validate output format
			if output != "text" && output != "json" && output != "yaml" {
				return fmt.Errorf("invalid output format: %s. Valid formats are: text, json, yaml", output)
			}

			// Handle all-namespaces flag
			if allNamespaces {
				namespace = ""
			}

			// Use context with timeout instead of context.TODO()
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			scanner, err := NewProbeScanner(kubeconfig, kubeContext, namespace, probeType, recommendation, output, failOnWarn)
			if err != nil {
				return fmt.Errorf("error creating scanner: %w", err)
			}

			return scanner.Scan(ctx)
		},
	}

	cmd.Flags().StringP("kubeconfig", "k", "", "Path to the kubeconfig file. If not provided, uses default kubeconfig location")
	cmd.Flags().StringP("kubeContext", "c", "", "Kubernetes context to use. If not provided, uses current context")
	cmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace to scan. Use --all-namespaces to scan all namespaces")
	cmd.Flags().StringP("probe-type", "p", "", "Type of probe to scan for: liveness, readiness, startup. If not provided, scans all types")
	cmd.Flags().BoolP("recommendation", "r", false, "Show detailed recommendations for missing probes")
	cmd.Flags().BoolP("all-namespaces", "A", false, "Scan all namespaces instead of a specific namespace")
	cmd.Flags().StringP("output", "o", "text", "Output format: text, json, or yaml")
	cmd.Flags().BoolP("fail-on-warn", "f", false, "Exit with code 1 if warnings are found (treats warnings as failures)")

	// Add custom completion for probe-type flag
	err := cmd.RegisterFlagCompletionFunc("probe-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"liveness", "readiness", "startup"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		// Log error but don't fail the command creation
		fmt.Printf("Warning: failed to register completion for probe-type flag: %v\n", err)
	}

	// Add custom completion for output flag
	err = cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"text", "json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		// Log error but don't fail the command creation
		fmt.Printf("Warning: failed to register completion for output flag: %v\n", err)
	}

	return cmd
}
