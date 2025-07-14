package scanner

import (
	"context"
	"fmt"
	"strings"

	"kubeprobes/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

var validProbeTypes = map[string]bool{
	"liveness":  true,
	"readiness": true,
	"startup":   true,
	"":          true, // empty string means all types
}

// KubernetesClient interface for testing
type KubernetesClient interface {
	GetPods(ctx context.Context, namespace string) (*corev1.PodList, error)
}

// ProbeIssuesFoundError indicates that probe issues were found during scanning
type ProbeIssuesFoundError struct {
	Message string
}

func (e *ProbeIssuesFoundError) Error() string {
	return e.Message
}

// ProbeScanner handles the scanning logic
type ProbeScanner struct {
	kubeClient     KubernetesClient
	namespace      string
	probeType      string
	recommendation bool
}

// NewProbeScanner creates a new probe scanner instance
func NewProbeScanner(kubeconfig, kubeContext, namespace, probeType string, recommendation bool) (*ProbeScanner, error) {
	if namespace == "" {
		namespace = "default"
	}

	probeType = strings.ToLower(probeType)
	if !validProbeTypes[probeType] {
		return nil, fmt.Errorf("invalid probe type '%s'\n\nValid probe types are:\n  - liveness: Checks if the container is running\n  - readiness: Checks if the container is ready to accept traffic\n  - startup: Checks if the container has started successfully\n\nExample: kubeprobes scan --probe-type liveness", probeType)
	}

	kubeClient, err := kubernetes.NewClient(kubeconfig, kubeContext)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kubernetes cluster: %w\n\nTroubleshooting tips:\n  - Ensure kubectl is configured and working: kubectl cluster-info\n  - Check kubeconfig file exists and is readable\n  - Verify the specified context exists: kubectl config get-contexts\n  - Try without specifying kubeconfig to use default: kubeprobes scan", err)
	}

	return NewProbeScannerWithClient(kubeClient, namespace, probeType, recommendation), nil
}

// NewProbeScannerWithClient creates a new probe scanner with a given client (useful for testing)
func NewProbeScannerWithClient(kubeClient KubernetesClient, namespace, probeType string, recommendation bool) *ProbeScanner {
	if namespace == "" {
		namespace = "default"
	}

	return &ProbeScanner{
		kubeClient:     kubeClient,
		namespace:      namespace,
		probeType:      strings.ToLower(probeType),
		recommendation: recommendation,
	}
}

// Scan performs the probe scanning
func (ps *ProbeScanner) Scan(ctx context.Context) error {
	pods, err := ps.kubeClient.GetPods(ctx, ps.namespace)
	if err != nil {
		return fmt.Errorf("failed to retrieve pods from namespace '%s': %w\n\nTroubleshooting tips:\n  - Verify the namespace exists: kubectl get namespaces\n  - Check if you have permissions: kubectl auth can-i get pods --namespace %s\n  - Ensure cluster connection is working: kubectl cluster-info", ps.namespace, err, ps.namespace)
	}

	if len(pods.Items) == 0 {
		fmt.Printf("ℹ️  No pods found in namespace '%s'\n\nSuggestions:\n  - Check if pods exist: kubectl get pods --namespace %s\n  - Try scanning a different namespace: kubeprobes scan --namespace <namespace>\n  - List all namespaces: kubectl get namespaces\n", ps.namespace, ps.namespace)
		return nil
	}

	issuesFound := false
	scannedContainers := 0
	
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			scannedContainers++
			if ps.probeType == "liveness" || ps.probeType == "" {
				if container.LivenessProbe == nil {
					issuesFound = true
					fmt.Printf("⚠️  [MISSING LIVENESS PROBE] Pod %s/%s (container: %s)\n",
						pod.Namespace, pod.Name, container.Name)
					if ps.recommendation {
						fmt.Printf("   💡 Recommendation: Add a liveness probe to detect if the container becomes unresponsive.\n")
						fmt.Printf("      Example: HTTP check on /health endpoint or exec command like 'ps aux | grep myapp'\n")
					}
				}
			}
			if ps.probeType == "readiness" || ps.probeType == "" {
				if container.ReadinessProbe == nil {
					issuesFound = true
					fmt.Printf("⚠️  [MISSING READINESS PROBE] Pod %s/%s (container: %s)\n",
						pod.Namespace, pod.Name, container.Name)
					if ps.recommendation {
						fmt.Printf("   💡 Recommendation: Add a readiness probe to ensure container is ready before receiving traffic.\n")
						fmt.Printf("      Example: HTTP check on /ready endpoint or TCP socket check on application port\n")
					}
				}
			}
			if ps.probeType == "startup" || ps.probeType == "" {
				if container.StartupProbe == nil {
					issuesFound = true
					fmt.Printf("⚠️  [MISSING STARTUP PROBE] Pod %s/%s (container: %s)\n",
						pod.Namespace, pod.Name, container.Name)
					if ps.recommendation {
						fmt.Printf("   💡 Recommendation: Add a startup probe for slow-starting containers to avoid premature kills.\n")
						fmt.Printf("      Example: HTTP check with longer initial delay and period for application startup\n")
					}
				}
			}
		}
	}

	if !issuesFound {
		fmt.Printf("✅ No probe issues found in namespace '%s' (scanned %d containers)\n", ps.namespace, scannedContainers)
		if !ps.recommendation {
			fmt.Printf("💡 Tip: Use --recommendation flag to see best practices for probe configuration\n")
		}
		return nil
	}

	fmt.Printf("\n📊 Summary: Found probe issues in namespace '%s' (scanned %d containers)\n", ps.namespace, scannedContainers)
	if !ps.recommendation {
		fmt.Printf("💡 Tip: Run with --recommendation flag to get actionable suggestions\n")
	}
	fmt.Printf("📚 Learn more: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/\n")
	
	return &ProbeIssuesFoundError{Message: "probe issues detected"}
}
