package scanner

import (
	"context"
	"errors"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// mockKubernetesClient implements a mock for testing
type mockKubernetesClient struct {
	pods *corev1.PodList
	err  error
}

func (m *mockKubernetesClient) GetPods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.pods, nil
}

func TestNewProbeScanner_Validation(t *testing.T) {
	tests := []struct {
		name           string
		probeType      string
		expectedError  string
	}{
		{
			name:           "valid liveness probe type",
			probeType:      "liveness",
			expectedError:  "",
		},
		{
			name:           "valid readiness probe type",
			probeType:      "readiness",
			expectedError:  "",
		},
		{
			name:           "valid startup probe type",
			probeType:      "startup",
			expectedError:  "",
		},
		{
			name:           "valid empty probe type",
			probeType:      "",
			expectedError:  "",
		},
		{
			name:           "case insensitive probe type",
			probeType:      "LIVENESS",
			expectedError:  "",
		},
		{
			name:           "invalid probe type",
			probeType:      "invalid",
			expectedError:  "invalid probe type: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use a temporary kubeconfig for testing validation only
			tempConfig := ""
			
			_, err := NewProbeScanner(tempConfig, "", "default", tt.probeType, false)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error containing %q, got none", tt.expectedError)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing %q, got %q", tt.expectedError, err.Error())
				}
				return
			}

			// For valid probe types, we expect a kubernetes client error since we don't have a real cluster
			if err != nil && !strings.Contains(err.Error(), "error creating kubernetes client") {
				t.Errorf("Unexpected error (expected kubernetes client error): %v", err)
			}
		})
	}
}

func TestNewProbeScannerWithClient(t *testing.T) {
	mockClient := &mockKubernetesClient{
		pods: &corev1.PodList{Items: []corev1.Pod{}},
		err:  nil,
	}

	tests := []struct {
		name           string
		namespace      string
		probeType      string
		recommendation bool
	}{
		{
			name:           "valid probe scanner with default namespace",
			namespace:      "",
			probeType:      "liveness",
			recommendation: false,
		},
		{
			name:           "valid probe scanner with specific namespace",
			namespace:      "test-namespace",
			probeType:      "readiness",
			recommendation: true,
		},
		{
			name:           "valid probe scanner with startup probe",
			namespace:      "default",
			probeType:      "startup",
			recommendation: false,
		},
		{
			name:           "valid probe scanner with empty probe type (all)",
			namespace:      "default",
			probeType:      "",
			recommendation: false,
		},
		{
			name:           "case insensitive probe type",
			namespace:      "default",
			probeType:      "LIVENESS",
			recommendation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewProbeScannerWithClient(mockClient, tt.namespace, tt.probeType, tt.recommendation)

			if scanner == nil {
				t.Error("Expected scanner to be non-nil")
				return
			}

			// Check that namespace defaults to "default" when empty
			expectedNamespace := tt.namespace
			if expectedNamespace == "" {
				expectedNamespace = "default"
			}
			if scanner.namespace != expectedNamespace {
				t.Errorf("Expected namespace %q, got %q", expectedNamespace, scanner.namespace)
			}

			// Check probe type is converted to lowercase
			expectedProbeType := strings.ToLower(tt.probeType)
			if scanner.probeType != expectedProbeType {
				t.Errorf("Expected probe type %q, got %q", expectedProbeType, scanner.probeType)
			}

			if scanner.recommendation != tt.recommendation {
				t.Errorf("Expected recommendation %v, got %v", tt.recommendation, scanner.recommendation)
			}
		})
	}
}

func TestProbeIssuesFoundError(t *testing.T) {
	err := &ProbeIssuesFoundError{Message: "test error message"}
	expected := "test error message"
	
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

func TestProbeScanner_Scan_NoPods(t *testing.T) {
	// Create a mock client that returns no pods
	mockClient := &mockKubernetesClient{
		pods: &corev1.PodList{Items: []corev1.Pod{}},
		err:  nil,
	}

	scanner := &ProbeScanner{
		kubeClient:     mockClient,
		namespace:      "test-namespace",
		probeType:      "",
		recommendation: false,
	}

	err := scanner.Scan(context.Background())
	if err != nil {
		t.Errorf("Unexpected error when no pods found: %v", err)
	}
}

func TestProbeScanner_Scan_GetPodsError(t *testing.T) {
	// Create a mock client that returns an error
	mockClient := &mockKubernetesClient{
		pods: nil,
		err:  errors.New("kubernetes api error"),
	}

	scanner := &ProbeScanner{
		kubeClient:     mockClient,
		namespace:      "test-namespace",
		probeType:      "",
		recommendation: false,
	}

	err := scanner.Scan(context.Background())
	if err == nil {
		t.Error("Expected error when GetPods fails")
		return
	}
	
	if !strings.Contains(err.Error(), "error listing pods") {
		t.Errorf("Expected error to contain 'error listing pods', got %q", err.Error())
	}
}

func TestProbeScanner_Scan_WithProbes(t *testing.T) {
	// Create pods with all probes configured
	pod := createTestPod("test-pod", "test-namespace", true, true, true)
	mockClient := &mockKubernetesClient{
		pods: &corev1.PodList{Items: []corev1.Pod{pod}},
		err:  nil,
	}

	tests := []struct {
		name      string
		probeType string
	}{
		{"all probe types", ""},
		{"liveness only", "liveness"},
		{"readiness only", "readiness"},
		{"startup only", "startup"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := &ProbeScanner{
				kubeClient:     mockClient,
				namespace:      "test-namespace",
				probeType:      tt.probeType,
				recommendation: false,
			}

			err := scanner.Scan(context.Background())
			if err != nil {
				t.Errorf("Unexpected error when all probes are configured: %v", err)
			}
		})
	}
}

func TestProbeScanner_Scan_MissingProbes(t *testing.T) {
	tests := []struct {
		name           string
		probeType      string
		hasLiveness    bool
		hasReadiness   bool
		hasStartup     bool
		expectIssues   bool
		recommendation bool
	}{
		{
			name:         "missing liveness probe - scan all",
			probeType:    "",
			hasLiveness:  false,
			hasReadiness: true,
			hasStartup:   true,
			expectIssues: true,
		},
		{
			name:         "missing readiness probe - scan all",
			probeType:    "",
			hasLiveness:  true,
			hasReadiness: false,
			hasStartup:   true,
			expectIssues: true,
		},
		{
			name:         "missing startup probe - scan all",
			probeType:    "",
			hasLiveness:  true,
			hasReadiness: true,
			hasStartup:   false,
			expectIssues: true,
		},
		{
			name:         "missing liveness probe - scan liveness only",
			probeType:    "liveness",
			hasLiveness:  false,
			hasReadiness: true,
			hasStartup:   true,
			expectIssues: true,
		},
		{
			name:         "missing readiness probe - scan liveness only",
			probeType:    "liveness",
			hasLiveness:  true,
			hasReadiness: false,
			hasStartup:   true,
			expectIssues: false,
		},
		{
			name:           "missing all probes with recommendations",
			probeType:      "",
			hasLiveness:    false,
			hasReadiness:   false,
			hasStartup:     false,
			expectIssues:   true,
			recommendation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := createTestPod("test-pod", "test-namespace", tt.hasLiveness, tt.hasReadiness, tt.hasStartup)
			mockClient := &mockKubernetesClient{
				pods: &corev1.PodList{Items: []corev1.Pod{pod}},
				err:  nil,
			}

			scanner := &ProbeScanner{
				kubeClient:     mockClient,
				namespace:      "test-namespace",
				probeType:      tt.probeType,
				recommendation: tt.recommendation,
			}

			err := scanner.Scan(context.Background())

			if tt.expectIssues {
				if err == nil {
					t.Error("Expected ProbeIssuesFoundError but got none")
					return
				}
				
				var probeErr *ProbeIssuesFoundError
				if !errors.As(err, &probeErr) {
					t.Errorf("Expected ProbeIssuesFoundError, got %T: %v", err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestProbeScanner_Scan_MultipleContainers(t *testing.T) {
	// Create a pod with multiple containers, some missing probes
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "multi-container-pod",
			Namespace: "test-namespace",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "container1",
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{Path: "/healthz"},
						},
					},
					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{Path: "/ready"},
						},
					},
					// Missing startup probe
				},
				{
					Name: "container2",
					// Missing all probes
				},
			},
		},
	}

	mockClient := &mockKubernetesClient{
		pods: &corev1.PodList{Items: []corev1.Pod{*pod}},
		err:  nil,
	}

	scanner := &ProbeScanner{
		kubeClient:     mockClient,
		namespace:      "test-namespace",
		probeType:      "",
		recommendation: false,
	}

	err := scanner.Scan(context.Background())
	if err == nil {
		t.Error("Expected ProbeIssuesFoundError for containers missing probes")
		return
	}

	var probeErr *ProbeIssuesFoundError
	if !errors.As(err, &probeErr) {
		t.Errorf("Expected ProbeIssuesFoundError, got %T: %v", err, err)
	}
}

// Helper function to create test pods
func createTestPod(name, namespace string, hasLiveness, hasReadiness, hasStartup bool) corev1.Pod {
	container := corev1.Container{
		Name: "test-container",
	}

	if hasLiveness {
		container.LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{Path: "/healthz"},
			},
		}
	}

	if hasReadiness {
		container.ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{Path: "/ready"},
			},
		}
	}

	if hasStartup {
		container.StartupProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{Path: "/startup"},
			},
		}
	}

	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{container},
		},
	}
}