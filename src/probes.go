package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"
)

var validProbeTypes = map[string]bool{
	"liveness":  true,
	"readiness": true,
	"startup":   true,
	"":          true, // empty string means all types
}

// ProbeIssuesFoundError indicates that probe issues were found during scanning
type ProbeIssuesFoundError struct {
	Message string
}

func (e *ProbeIssuesFoundError) Error() string {
	return e.Message
}

// validateProbeType explicitly validates the probe type flag
func validateProbeType(probeType string) error {
	if !validProbeTypes[strings.ToLower(probeType)] {
		return fmt.Errorf("invalid probe type: %s. Valid types are: liveness, readiness, startup", probeType)
	}
	return nil
}

// scanProbes scans pods for missing probes and returns true if issues are found
func scanProbes(ctx context.Context, clientset *kubernetes.Clientset, namespace, probeType string, showRecommendations bool) (bool, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("error listing pods: %w", err)
	}

	if len(pods.Items) == 0 {
		fmt.Printf("No pods found in namespace %s\n", namespace)
		return false, nil
	}

	issuesFound := false
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if probeType == "liveness" || probeType == "" {
				if container.LivenessProbe == nil {
					issuesFound = true
					fmt.Printf("[WARNING] Pod %s/%s (container: %s) is missing a liveness probe\n",
						pod.Namespace, pod.Name, container.Name)
					if showRecommendations {
						fmt.Println("  Recommendation: Add a liveness probe to ensure the container is running correctly.")
					}
				}
			}
			if probeType == "readiness" || probeType == "" {
				if container.ReadinessProbe == nil {
					issuesFound = true
					fmt.Printf("[WARNING] Pod %s/%s (container: %s) is missing a readiness probe\n",
						pod.Namespace, pod.Name, container.Name)
					if showRecommendations {
						fmt.Println("  Recommendation: Add a readiness probe to ensure the container is ready to accept traffic.")
					}
				}
			}
			if probeType == "startup" || probeType == "" {
				if container.StartupProbe == nil {
					issuesFound = true
					fmt.Printf("[WARNING] Pod %s/%s (container: %s) is missing a startup probe\n",
						pod.Namespace, pod.Name, container.Name)
					if showRecommendations {
						fmt.Println("  Recommendation: Add a startup probe to ensure the container has started successfully.")
					}
				}
			}
		}
	}

	if !issuesFound {
		fmt.Printf("No probe issues found in namespace %s\n", namespace)
	}

	return issuesFound, nil
}

var rootCmd = &cobra.Command{
	Use:   "probes",
	Short: "Probes is a CLI tool for scanning Kubernetes probes",
}

var scanCmd = &cobra.Command{
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
		if namespace == "" {
			namespace = "default"
		}
		if err != nil {
			return fmt.Errorf("error getting namespace flag: %w", err)
		}

		probeType, err := cmd.Flags().GetString("probe-type")
		if err != nil {
			return fmt.Errorf("error getting probe-type flag: %w", err)
		}

		// Explicit validation for probe-type flag
		if err := validateProbeType(probeType); err != nil {
			return err
		}
		probeType = strings.ToLower(probeType)

		recommendation, err := cmd.Flags().GetBool("recommendation")
		if err != nil {
			return fmt.Errorf("error getting recommendation flag: %w", err)
		}

		config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			&clientcmd.ConfigOverrides{CurrentContext: kubeContext}).ClientConfig()
		if err != nil {
			return fmt.Errorf("error building kubeconfig: %w", err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("error building kubernetes clientset: %w", err)
		}

		// Use context with timeout instead of context.TODO()
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		issuesFound, err := scanProbes(ctx, clientset, namespace, probeType, recommendation)
		if err != nil {
			return err
		}

		if issuesFound {
			fmt.Println("Issues found. Exiting with status code 1.")
			return &ProbeIssuesFoundError{Message: "probe issues found"}
		}

		return nil
	},
}

func main() {
	scanCmd.Flags().StringP("kubeconfig", "k", "", "path to the kubeconfig file")
	scanCmd.Flags().StringP("kubeContext", "c", "", "Kubernetes context")
	scanCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace (default: default)")
	scanCmd.Flags().StringP("probe-type", "p", "", "type of probe to scan for (liveness, readiness, startup)")
	scanCmd.Flags().BoolP("recommendation", "r", false, "show recommendations for missing probes")

	rootCmd.AddCommand(scanCmd)
	if err := rootCmd.Execute(); err != nil {
		// Check if it's our custom error indicating probe issues found
		if _, ok := err.(*ProbeIssuesFoundError); ok {
			os.Exit(1)
		}
		log.Fatalf("Error executing command: %s", err.Error())
	}
}
