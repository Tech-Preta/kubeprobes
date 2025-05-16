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

var rootCmd = &cobra.Command{
	Use:   "probes",
	Short: "Probes is a CLI tool for scanning Kubernetes probes",
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan Kubernetes workloads for probes",
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
		if namespace == "" {
			namespace = "default"
		}

		probeType, err := cmd.Flags().GetString("probe-type")
		if err != nil {
			log.Fatalf("Error getting probe-type flag: %s", err.Error())
		}
		if !validProbeTypes[strings.ToLower(probeType)] {
			log.Fatalf("Invalid probe type: %s. Valid types are: liveness, readiness, startup", probeType)
		}

		recommendation, err := cmd.Flags().GetBool("recommendation")
		if err != nil {
			log.Fatalf("Error getting recommendation flag: %s", err.Error())
		}

		config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			&clientcmd.ConfigOverrides{CurrentContext: kubeContext}).ClientConfig()
		if err != nil {
			log.Fatalf("Error building kubeconfig: %s", err.Error())
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error building kubernetes clientset: %s", err.Error())
		}

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing pods: %s", err.Error())
		}

		if len(pods.Items) == 0 {
			fmt.Printf("No pods found in namespace %s\n", namespace)
			os.Exit(0)
		}

		issuesFound := false
		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				if probeType == "liveness" || probeType == "" {
					if container.LivenessProbe == nil {
						issuesFound = true
						fmt.Printf("[WARNING] Pod %s/%s (container: %s) is missing a liveness probe\n", 
							pod.Namespace, pod.Name, container.Name)
						if recommendation {
							fmt.Println("  Recommendation: Add a liveness probe to ensure the container is running correctly.")
						}
					}
				}
				if probeType == "readiness" || probeType == "" {
					if container.ReadinessProbe == nil {
						issuesFound = true
						fmt.Printf("[WARNING] Pod %s/%s (container: %s) is missing a readiness probe\n", 
							pod.Namespace, pod.Name, container.Name)
						if recommendation {
							fmt.Println("  Recommendation: Add a readiness probe to ensure the container is ready to accept traffic.")
						}
					}
				}
				if probeType == "startup" || probeType == "" {
					if container.StartupProbe == nil {
						issuesFound = true
						fmt.Printf("[WARNING] Pod %s/%s (container: %s) is missing a startup probe\n", 
							pod.Namespace, pod.Name, container.Name)
						if recommendation {
							fmt.Println("  Recommendation: Add a startup probe to ensure the container has started successfully.")
						}
					}
				}
			}
		}

		if !issuesFound {
			fmt.Printf("No probe issues found in namespace %s\n", namespace)
		} else {
			os.Exit(1)
		}
	},
}

// main initializes CLI flags, sets up commands, and starts the probes CLI tool.
func main() {
	scanCmd.Flags().StringP("kubeconfig", "k", "", "path to the kubeconfig file")
	scanCmd.Flags().StringP("kubeContext", "c", "", "Kubernetes context")
	scanCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace (default: default)")
	scanCmd.Flags().StringP("probe-type", "p", "", "type of probe to scan for (liveness, readiness, startup)")
	scanCmd.Flags().BoolP("recommendation", "r", false, "show recommendations for missing probes")

	rootCmd.AddCommand(scanCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %s", err.Error())
	}
}
