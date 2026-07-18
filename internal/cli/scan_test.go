package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestScanCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "help flag",
			args:           []string{"--help"},
			expectedOutput: "Scan Kubernetes workloads for missing liveness, readiness, or startup probes",
			expectError:    false,
		},
		{
			name:           "invalid kubeconfig should error",
			args:           []string{"--kubeconfig=/nonexistent/path"},
			expectedOutput: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewScanCommand()

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectedOutput != "" {
				output := buf.String()
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
				}
			}
		})
	}
}

func TestScanCommandFlags(t *testing.T) {
	cmd := NewScanCommand()

	expectedFlags := []string{
		"kubeconfig",
		"kubeContext",
		"namespace",
		"probe-type",
		"recommendation",
	}

	for _, flagName := range expectedFlags {
		flag := cmd.Flag(flagName)
		if flag == nil {
			t.Errorf("Expected flag %q to be available", flagName)
		}
	}
}

func TestScanCommandCompletions(t *testing.T) {
	cmd := NewScanCommand()

	// Test probe-type flag completion
	flag := cmd.Flag("probe-type")
	if flag == nil {
		t.Fatal("probe-type flag should exist")
	}

	// The completion function should be registered for probe-type
	// We can't easily test the actual completion here without mocking,
	// but we can verify the flag exists and has proper setup
	if flag.Usage == "" {
		t.Error("probe-type flag should have usage text")
	}
}
