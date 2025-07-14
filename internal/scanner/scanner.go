package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"kubeprobes/pkg/kubernetes"
	"gopkg.in/yaml.v3"
)

var validProbeTypes = map[string]bool{
	"liveness":  true,
	"readiness": true,
	"startup":   true,
	"":          true, // empty string means all types
}

// ProbeIssue represents a single probe issue found
type ProbeIssue struct {
	Namespace     string `json:"namespace" yaml:"namespace"`
	PodName       string `json:"podName" yaml:"podName"`
	ContainerName string `json:"containerName" yaml:"containerName"`
	ProbeType     string `json:"probeType" yaml:"probeType"`
	Message       string `json:"message" yaml:"message"`
	Recommendation string `json:"recommendation,omitempty" yaml:"recommendation,omitempty"`
}

// ScanResult represents the complete scan result
type ScanResult struct {
	Issues    []ProbeIssue `json:"issues" yaml:"issues"`
	Summary   string       `json:"summary" yaml:"summary"`
	Namespace string       `json:"namespace" yaml:"namespace"`
	ExitCode  int          `json:"exitCode" yaml:"exitCode"`
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
	output         string
	failOnWarn     bool
}

// NewProbeScanner creates a new probe scanner instance
func NewProbeScanner(kubeconfig, kubeContext, namespace, probeType string, recommendation bool, output string, failOnWarn bool) (*ProbeScanner, error) {
	if namespace == "" {
		// Empty namespace means all namespaces
	} else if namespace == "default" {
		// Keep default namespace as is
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
		output:         output,
		failOnWarn:     failOnWarn,
	}, nil
}

// Scan performs the probe scanning
func (ps *ProbeScanner) Scan(ctx context.Context) error {
	pods, err := ps.kubeClient.GetPods(ctx, ps.namespace)
	if err != nil {
		return fmt.Errorf("error listing pods: %w", err)
	}

	namespaceDisplay := ps.namespace
	if ps.namespace == "" {
		namespaceDisplay = "all namespaces"
	}

	if len(pods.Items) == 0 {
		summary := fmt.Sprintf("No pods found in %s", namespaceDisplay)
		result := ScanResult{
			Issues:    []ProbeIssue{},
			Summary:   summary,
			Namespace: namespaceDisplay,
			ExitCode:  0,
		}
		
		return ps.outputResult(result)
	}

	var issues []ProbeIssue
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			// Check liveness probe
			if ps.probeType == "liveness" || ps.probeType == "" {
				if container.LivenessProbe == nil {
					issue := ProbeIssue{
						Namespace:     pod.Namespace,
						PodName:       pod.Name,
						ContainerName: container.Name,
						ProbeType:     "liveness",
						Message:       "missing liveness probe",
					}
					if ps.recommendation {
						issue.Recommendation = "Add a liveness probe to ensure the container is running correctly."
					}
					issues = append(issues, issue)
				}
			}
			
			// Check readiness probe
			if ps.probeType == "readiness" || ps.probeType == "" {
				if container.ReadinessProbe == nil {
					issue := ProbeIssue{
						Namespace:     pod.Namespace,
						PodName:       pod.Name,
						ContainerName: container.Name,
						ProbeType:     "readiness",
						Message:       "missing readiness probe",
					}
					if ps.recommendation {
						issue.Recommendation = "Add a readiness probe to ensure the container is ready to accept traffic."
					}
					issues = append(issues, issue)
				}
			}
			
			// Check startup probe
			if ps.probeType == "startup" || ps.probeType == "" {
				if container.StartupProbe == nil {
					issue := ProbeIssue{
						Namespace:     pod.Namespace,
						PodName:       pod.Name,
						ContainerName: container.Name,
						ProbeType:     "startup",
						Message:       "missing startup probe",
					}
					if ps.recommendation {
						issue.Recommendation = "Add a startup probe to ensure the container has started successfully."
					}
					issues = append(issues, issue)
				}
			}
		}
	}

	var summary string
	var exitCode int
	
	if len(issues) == 0 {
		summary = fmt.Sprintf("No probe issues found in %s", namespaceDisplay)
		exitCode = 0
	} else {
		summary = fmt.Sprintf("Found %d probe issues in %s", len(issues), namespaceDisplay)
		if ps.failOnWarn {
			exitCode = 1
		} else {
			exitCode = 0
		}
	}

	result := ScanResult{
		Issues:    issues,
		Summary:   summary,
		Namespace: namespaceDisplay,
		ExitCode:  exitCode,
	}

	err = ps.outputResult(result)
	if err != nil {
		return err
	}

	// Handle exit code
	if exitCode == 1 {
		return &ProbeIssuesFoundError{Message: "probe issues found"}
	}

	return nil
}

// outputResult outputs the scan result in the requested format
func (ps *ProbeScanner) outputResult(result ScanResult) error {
	switch ps.output {
	case "json":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling to JSON: %w", err)
		}
		fmt.Println(string(data))
		
	case "yaml":
		data, err := yaml.Marshal(result)
		if err != nil {
			return fmt.Errorf("error marshaling to YAML: %w", err)
		}
		fmt.Print(string(data))
		
	default: // text format
		if len(result.Issues) == 0 {
			fmt.Println(result.Summary)
		} else {
			for _, issue := range result.Issues {
				fmt.Printf("[WARNING] Pod %s/%s (container: %s) is %s\n",
					issue.Namespace, issue.PodName, issue.ContainerName, issue.Message)
				if issue.Recommendation != "" {
					fmt.Printf("  Recommendation: %s\n", issue.Recommendation)
				}
			}
			if ps.failOnWarn {
				fmt.Println("Issues found. Exiting with status code 1.")
			}
		}
	}
	
	return nil
}
