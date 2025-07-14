package scanner

import (
	"strings"
	"testing"
)

func TestNewProbeScannerErrorMessages(t *testing.T) {
	tests := []struct {
		name          string
		probeType     string
		expectedError string
		shouldError   bool
	}{
		{
			name:          "invalid probe type shows descriptive error",
			probeType:     "invalid",
			expectedError: "invalid probe type 'invalid'",
			shouldError:   true,
		},
		{
			name:          "invalid probe type includes valid options",
			probeType:     "wrong",
			expectedError: "liveness: Checks if the container is running",
			shouldError:   true,
		},
		{
			name:          "invalid probe type includes example",
			probeType:     "bad",
			expectedError: "Example: kubeprobes scan --probe-type liveness",
			shouldError:   true,
		},
		{
			name:        "valid liveness probe type",
			probeType:   "liveness",
			shouldError: false,
		},
		{
			name:        "valid readiness probe type",
			probeType:   "readiness",
			shouldError: false,
		},
		{
			name:        "valid startup probe type",
			probeType:   "startup",
			shouldError: false,
		},
		{
			name:        "empty probe type is valid",
			probeType:   "",
			shouldError: false,
		},
		{
			name:        "case insensitive probe type",
			probeType:   "LIVENESS",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProbeScanner("", "", "default", tt.probeType, false)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error to contain %q, got %q", tt.expectedError, err.Error())
				}
			} else {
				// We expect connection errors since we're not providing valid kubeconfig
				// but the probe type validation should pass
				if err != nil && strings.Contains(err.Error(), "invalid probe type") {
					t.Errorf("Unexpected probe type validation error: %v", err)
				}
			}
		})
	}
}

func TestProbeIssuesFoundError(t *testing.T) {
	err := &ProbeIssuesFoundError{Message: "test error"}
	
	if err.Error() != "test error" {
		t.Errorf("Expected error message 'test error', got %q", err.Error())
	}
}

func TestValidProbeTypes(t *testing.T) {
	expectedTypes := map[string]bool{
		"liveness":  true,
		"readiness": true,
		"startup":   true,
		"":          true,
	}

	for probeType, expected := range expectedTypes {
		if validProbeTypes[probeType] != expected {
			t.Errorf("Expected validProbeTypes[%q] to be %v, got %v", probeType, expected, validProbeTypes[probeType])
		}
	}

	// Test invalid types
	invalidTypes := []string{"invalid", "wrong", "bad", "test"}
	for _, probeType := range invalidTypes {
		if validProbeTypes[probeType] {
			t.Errorf("Expected validProbeTypes[%q] to be false, got true", probeType)
		}
	}
}