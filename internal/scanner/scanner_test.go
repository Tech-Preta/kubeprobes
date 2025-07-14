package scanner

import (
	"testing"
)

func TestValidProbeTypes(t *testing.T) {
	testCases := []struct {
		probeType string
		isValid   bool
	}{
		{"liveness", true},
		{"readiness", true},
		{"startup", true},
		{"", true}, // empty means all types
		{"invalid", false},
		{"LIVENESS", false}, // case sensitive before normalization
	}

	for _, tc := range testCases {
		_, exists := validProbeTypes[tc.probeType]
		if exists != tc.isValid {
			t.Errorf("Probe type '%s' validity: expected %v, got %v", tc.probeType, tc.isValid, exists)
		}
	}
}

func TestNewProbeScannerValidation(t *testing.T) {
	testCases := []struct {
		name          string
		probeType     string
		output        string
		shouldError   bool
		expectedError string
	}{
		{
			name:        "valid liveness probe",
			probeType:   "liveness",
			output:      "text",
			shouldError: false,
		},
		{
			name:        "valid empty probe type",
			probeType:   "",
			output:      "json",
			shouldError: false,
		},
		{
			name:          "invalid probe type",
			probeType:     "invalid",
			output:        "text",
			shouldError:   true,
			expectedError: "invalid probe type",
		},
		{
			name:        "case insensitive probe type",
			probeType:   "LIVENESS",
			output:      "yaml",
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scanner, err := NewProbeScanner("", "", "test", tc.probeType, false, tc.output, false)
			
			if tc.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tc.expectedError != "" && !contains(err.Error(), tc.expectedError) {
					t.Errorf("Expected error to contain '%s', got '%s'", tc.expectedError, err.Error())
				}
			} else {
				if err != nil {
					// We expect kubernetes client creation to fail, but not probe type validation
					if contains(err.Error(), "invalid probe type") {
						t.Errorf("Unexpected probe type validation error: %v", err)
					}
				} else if scanner == nil {
					t.Error("Expected scanner to be created")
				}
			}
		})
	}
}

func TestProbeIssueStruct(t *testing.T) {
	issue := ProbeIssue{
		Namespace:     "test-ns",
		PodName:       "test-pod",
		ContainerName: "test-container",
		ProbeType:     "liveness",
		Message:       "missing liveness probe",
		Recommendation: "Add a liveness probe",
	}

	if issue.Namespace != "test-ns" {
		t.Errorf("Expected namespace 'test-ns', got '%s'", issue.Namespace)
	}
	if issue.ProbeType != "liveness" {
		t.Errorf("Expected probe type 'liveness', got '%s'", issue.ProbeType)
	}
}

func TestScanResultStruct(t *testing.T) {
	issues := []ProbeIssue{
		{
			Namespace: "test",
			PodName:   "pod1",
			ProbeType: "liveness",
		},
	}

	result := ScanResult{
		Issues:    issues,
		Summary:   "Found 1 issue",
		Namespace: "test",
		ExitCode:  0,
	}

	if len(result.Issues) != 1 {
		t.Errorf("Expected 1 issue, got %d", len(result.Issues))
	}
	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}