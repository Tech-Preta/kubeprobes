package main

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "probes",
	Short: "Probes is a CLI tool for scanning Kubernetes probes",
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan Kubernetes workloads for probes",
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		kubeContext, _ := cmd.Flags().GetString("kubeContext")
		namespace, _ := cmd.Flags().GetString("namespace")
		probeType, _ := cmd.Flags().GetString("probe-type")
		recommendation, _ := cmd.Flags().GetBool("recommendation")

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

		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				if probeType == "liveness" || probeType == "" {
					if container.LivenessProbe == nil {
						fmt.Printf("Pod %s in namespace %s does not have a liveness probe\n", pod.Name, pod.Namespace)
						if recommendation {
							fmt.Println("Recommendation: Add a liveness probe to ensure the container is running correctly.")
						}
					}
				}
				if probeType == "readiness" || probeType == "" {
					if container.ReadinessProbe == nil {
						fmt.Printf("Pod %s in namespace %s does not have a readiness probe\n", pod.Name, pod.Namespace)
						if recommendation {
							fmt.Println("Recommendation: Add a readiness probe to ensure the container is ready to accept traffic.")
						}
					}
				}
				if probeType == "startup" || probeType == "" {
					if container.StartupProbe == nil {
						fmt.Printf("Pod %s in namespace %s does not have a startup probe\n", pod.Name, pod.Namespace)
						if recommendation {
							fmt.Println("Recommendation: Add a startup probe to ensure the container has started successfully.")
						}
					}
				}
			}
		}
	},
}

func main() {
	scanCmd.Flags().StringP("kubeconfig", "k", "", "path to the kubeconfig file")
	scanCmd.Flags().StringP("kubeContext", "c", "", "Kubernetes context")
	scanCmd.Flags().StringP("namespace", "n", "", "Kubernetes namespace")
	scanCmd.Flags().StringP("probe-type", "p", "", "type of probe to scan for (liveness, readiness, startup)")
	scanCmd.Flags().BoolP("recommendation", "r", false, "show recommendations for missing probes")

	rootCmd.AddCommand(scanCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %s", err.Error())
	}
}
