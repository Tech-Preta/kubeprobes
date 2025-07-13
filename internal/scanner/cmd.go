package scanner

import (
	"log"

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
		Run: func(cmd *cobra.Command, args []string) {
			kubeconfig, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				log.Fatalf("Error getting kubeconfig flag: %s", err.Error())
			}

			kubeContext, err := cmd.Flags().GetString("kubeContext")
			if err != nil {
				log.Fatalf("Error getting kubeContext flag: %s", err.Error())
			}

			namespace, err := cmd.Flags().GetString("namespace")
			if err != nil {
				log.Fatalf("Error getting namespace flag: %s", err.Error())
			}

			probeType, err := cmd.Flags().GetString("probe-type")
			if err != nil {
				log.Fatalf("Error getting probe-type flag: %s", err.Error())
			}

			recommendation, err := cmd.Flags().GetBool("recommendation")
			if err != nil {
				log.Fatalf("Error getting recommendation flag: %s", err.Error())
			}

			scanner, err := NewProbeScanner(kubeconfig, kubeContext, namespace, probeType, recommendation)
			if err != nil {
				log.Fatalf("Error creating scanner: %s", err.Error())
			}

			if err := scanner.Scan(); err != nil {
				log.Fatalf("Error during scan: %s", err.Error())
			}
		},
	}

	cmd.Flags().StringP("kubeconfig", "k", "", "path to the kubeconfig file")
	cmd.Flags().StringP("kubeContext", "c", "", "Kubernetes context")
	cmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace (default: default)")
	cmd.Flags().StringP("probe-type", "p", "", "type of probe to scan for (liveness, readiness, startup)")
	cmd.Flags().BoolP("recommendation", "r", false, "show recommendations for missing probes")

	return cmd
}
