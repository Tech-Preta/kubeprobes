package scanner

import (
	"context"
	"fmt"
	"strings"

	"kubeprobes/pkg/kubernetes"
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

// ProbeScanner handles the scanning logic
type ProbeScanner struct {
	kubeClient     *kubernetes.Client
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
		return nil, fmt.Errorf("invalid probe type: %s. Valid types are: liveness, readiness, startup", probeType)
	}

	kubeClient, err := kubernetes.NewClient(kubeconfig, kubeContext)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes client: %w", err)
	}

	return &ProbeScanner{
		kubeClient:     kubeClient,
		namespace:      namespace,
		probeType:      probeType,
		recommendation: recommendation,
	}, nil
}

// Scan performs the probe scanning
func (ps *ProbeScanner) Scan(ctx context.Context) error {
	pods, err := ps.kubeClient.GetPods(ctx, ps.namespace)
	if err != nil {
		return fmt.Errorf("error listing pods: %w", err)
	}

	if len(pods.Items) == 0 {
		fmt.Printf("No pods found in namespace %s\n", ps.namespace)
		return nil
	}

	issuesFound := false
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if ps.probeType == "liveness" || ps.probeType == "" {
				if container.LivenessProbe == nil {
					issuesFound = true
					fmt.Printf("[WARNING] Pod %s/%s (container: %s) is missing a liveness probe\n",
						pod.Namespace, pod.Name, container.Name)
					if ps.recommendation {
						fmt.Println("  Recommendation: Add a liveness probe to ensure the container is running correctly.")
					}
				}
			}
			if ps.probeType == "readiness" || ps.probeType == "" {
				if container.ReadinessProbe == nil {
					issuesFound = true
					fmt.Printf("[WARNING] Pod %s/%s (container: %s) is missing a readiness probe\n",
						pod.Namespace, pod.Name, container.Name)
					if ps.recommendation {
						fmt.Println("  Recommendation: Add a readiness probe to ensure the container is ready to accept traffic.")
					}
				}
			}
			if ps.probeType == "startup" || ps.probeType == "" {
				if container.StartupProbe == nil {
					issuesFound = true
					fmt.Printf("[WARNING] Pod %s/%s (container: %s) is missing a startup probe\n",
						pod.Namespace, pod.Name, container.Name)
					if ps.recommendation {
						fmt.Println("  Recommendation: Add a startup probe to ensure the container has started successfully.")
					}
				}
			}
		}
	}

	if !issuesFound {
		fmt.Printf("No probe issues found in namespace %s\n", ps.namespace)
		return nil
	}

	fmt.Println("Issues found. Exiting with status code 1.")
	return &ProbeIssuesFoundError{Message: "probe issues found"}
}
