package scanner

import (
	"bytes"
	"strings"
	"testing"
)

func TestScanCommandHelpText(t *testing.T) {
	cmd := NewScanCommand()

	// Capture help output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Test that help contains improved descriptions
	expectedContents := []string{
		"Scan Kubernetes workloads for missing liveness, readiness, or startup probes",
		"This command connects to your Kubernetes cluster",
		"Exit codes:",
		"0: No probe issues found",
		"1: Probe issues detected",
		"Examples:",
		"kubeconfig file (defaults to $KUBECONFIG or ~/.kube/config)",
		"Show actionable recommendations for missing probes",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected help text to contain %q", expected)
		}
	}
}

func TestScanCommandExamples(t *testing.T) {
	cmd := NewScanCommand()

	// Capture help output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Test that examples are present and useful
	expectedExamples := []string{
		"kubeprobes scan",
		"kubeprobes scan --recommendation",
		"kubeprobes scan --namespace my-app --probe-type liveness",
		"kubeprobes scan --kubeconfig ~/.kube/config",
		"kubeprobes scan --kubeContext production",
	}

	for _, example := range expectedExamples {
		if !strings.Contains(output, example) {
			t.Errorf("Expected help text to contain example %q", example)
		}
	}
}

func TestScanCommandFlagDescriptions(t *testing.T) {
	cmd := NewScanCommand()

	// Test that flags have improved descriptions
	expectedFlags := map[string]string{
		"kubeconfig":     "Path to the kubeconfig file (defaults to $KUBECONFIG or ~/.kube/config)",
		"kubeContext":    "Kubernetes context to use (defaults to current context)",
		"namespace":      "Kubernetes namespace to scan (default: default)",
		"probe-type":     "Type of probe to scan for: liveness, readiness, or startup (default: all types)",
		"recommendation": "Show actionable recommendations for missing probes",
	}

	for flagName, expectedUsage := range expectedFlags {
		flag := cmd.Flag(flagName)
		if flag == nil {
			t.Errorf("Flag %q should exist", flagName)
			continue
		}

		if flag.Usage != expectedUsage {
			t.Errorf("Flag %q usage text should be %q, got %q", flagName, expectedUsage, flag.Usage)
		}
	}
}

func TestScanCommandErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorText   string
	}{
		{
			name:        "invalid probe type shows helpful error",
			args:        []string{"--probe-type", "invalid"},
			expectError: true,
			errorText:   "invalid probe type 'invalid'",
		},
		{
			name:        "nonexistent kubeconfig shows helpful error",
			args:        []string{"--kubeconfig", "/nonexistent/file"},
			expectError: true,
			errorText:   "failed to connect to Kubernetes cluster",
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

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorText, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}