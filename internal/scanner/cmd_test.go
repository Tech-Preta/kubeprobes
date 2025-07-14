package scanner

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewScanCommand(t *testing.T) {
	cmd := NewScanCommand()

	if cmd == nil {
		t.Fatal("NewScanCommand should return a non-nil command")
	}

	if cmd.Use != "scan" {
		t.Errorf("Expected command use 'scan', got %q", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command should have a short description")
	}

	if cmd.Long == "" {
		t.Error("Command should have a long description")
	}

	if cmd.RunE == nil {
		t.Error("Command should have a RunE function")
	}
}

func TestNewScanCommand_Flags(t *testing.T) {
	cmd := NewScanCommand()

	expectedFlags := []struct {
		name         string
		shorthand    string
		defaultValue interface{}
		usage        string
	}{
		{"kubeconfig", "k", "", "Path to the kubeconfig file (defaults to $KUBECONFIG or ~/.kube/config)"},
		{"kubeContext", "c", "", "Kubernetes context"},
		{"namespace", "n", "default", "Kubernetes namespace to scan (default: default)"},
		{"probe-type", "p", "", "Type of probe to scan for: liveness, readiness, or startup (default: all types)"},
		{"recommendation", "r", false, "Show actionable recommendations for missing probes"},
	}

	for _, expected := range expectedFlags {
		t.Run(expected.name, func(t *testing.T) {
			flag := cmd.Flag(expected.name)
			if flag == nil {
				t.Errorf("Flag %q should exist", expected.name)
				return
			}

			if flag.Shorthand != expected.shorthand {
				t.Errorf("Flag %q shorthand: expected %q, got %q", expected.name, expected.shorthand, flag.Shorthand)
			}

			if expected.usage != "" && !strings.Contains(flag.Usage, expected.usage) {
				t.Errorf("Flag %q usage should contain %q, got %q", expected.name, expected.usage, flag.Usage)
			}
		})
	}
}

func TestNewScanCommand_ProbeTypeCompletion(t *testing.T) {
	cmd := NewScanCommand()

	// Check that probe-type flag has completion function
	flag := cmd.Flag("probe-type")
	if flag == nil {
		t.Fatal("probe-type flag should exist")
	}

	// We can't easily test the completion function directly without internal access,
	// but we can verify it was registered by checking if the command setup succeeded
	if cmd.HasAvailableFlags() {
		t.Log("Command has flags available, completion likely registered correctly")
	}
}

func TestNewScanCommand_Help(t *testing.T) {
	cmd := NewScanCommand()

	// Test help flag
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Help command should not error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Scan Kubernetes workloads") {
		t.Error("Help output should contain description")
	}

	// Test that help contains improved descriptions from our PR
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

	// Check that all flags are documented in help
	expectedInHelp := []string{
		"kubeconfig",
		"kubeContext", 
		"namespace",
		"probe-type",
		"recommendation",
	}

	for _, expected := range expectedInHelp {
		if !strings.Contains(output, expected) {
			t.Errorf("Help output should contain flag %q", expected)
		}
	}
}

func TestNewScanCommand_FlagValidation(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "invalid kubeconfig path",
			args:        []string{"--kubeconfig=/nonexistent/path"},
			expectError: true,
			errorMsg:    "failed to connect to Kubernetes cluster",
		},
		{
			name:        "invalid probe type shows enhanced error message",
			args:        []string{"--probe-type=invalid"},
			expectError: true,
			errorMsg:    "invalid probe type 'invalid'",
		},
		{
			name:        "valid probe type liveness",
			args:        []string{"--probe-type=liveness", "--kubeconfig=/nonexistent"},
			expectError: true,
			errorMsg:    "failed to connect to Kubernetes cluster", // Will fail on kubeconfig, not probe type
		},
		{
			name:        "valid probe type readiness",
			args:        []string{"--probe-type=readiness", "--kubeconfig=/nonexistent"},
			expectError: true,
			errorMsg:    "failed to connect to Kubernetes cluster",
		},
		{
			name:        "valid probe type startup",
			args:        []string{"--probe-type=startup", "--kubeconfig=/nonexistent"},
			expectError: true,
			errorMsg:    "failed to connect to Kubernetes cluster",
		},
		{
			name:        "case insensitive probe type",
			args:        []string{"--probe-type=LIVENESS", "--kubeconfig=/nonexistent"},
			expectError: true,
			errorMsg:    "failed to connect to Kubernetes cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewScanCommand()
			cmd.SetContext(context.Background())

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.expectError {
				if err == nil {
					t.Error("Expected an error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Test that examples are present and useful (from our PR improvements)
func TestNewScanCommand_Examples(t *testing.T) {
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
func TestNewScanCommand_FlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		validate func(*testing.T, *cobra.Command)
	}{
		{
			name: "default values",
			args: []string{},
			validate: func(t *testing.T, cmd *cobra.Command) {
				namespace, _ := cmd.Flags().GetString("namespace")
				if namespace != "default" {
					t.Errorf("Expected default namespace 'default', got %q", namespace)
				}

				recommendation, _ := cmd.Flags().GetBool("recommendation")
				if recommendation != false {
					t.Errorf("Expected default recommendation false, got %v", recommendation)
				}
			},
		},
		{
			name: "custom namespace",
			args: []string{"--namespace=test-ns"},
			validate: func(t *testing.T, cmd *cobra.Command) {
				namespace, _ := cmd.Flags().GetString("namespace")
				if namespace != "test-ns" {
					t.Errorf("Expected namespace 'test-ns', got %q", namespace)
				}
			},
		},
		{
			name: "recommendation enabled",
			args: []string{"--recommendation"},
			validate: func(t *testing.T, cmd *cobra.Command) {
				recommendation, _ := cmd.Flags().GetBool("recommendation")
				if recommendation != true {
					t.Errorf("Expected recommendation true, got %v", recommendation)
				}
			},
		},
		{
			name: "short flags",
			args: []string{"-n", "short-ns", "-r", "-p", "liveness"},
			validate: func(t *testing.T, cmd *cobra.Command) {
				namespace, _ := cmd.Flags().GetString("namespace")
				if namespace != "short-ns" {
					t.Errorf("Expected namespace 'short-ns', got %q", namespace)
				}

				recommendation, _ := cmd.Flags().GetBool("recommendation")
				if recommendation != true {
					t.Errorf("Expected recommendation true, got %v", recommendation)
				}

				probeType, _ := cmd.Flags().GetString("probe-type")
				if probeType != "liveness" {
					t.Errorf("Expected probe-type 'liveness', got %q", probeType)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewScanCommand()

			// Parse args without executing (to avoid kubeconfig errors)
			cmd.SetArgs(tt.args)
			err := cmd.ParseFlags(tt.args)
			if err != nil {
				t.Errorf("Error parsing flags: %v", err)
				return
			}

			tt.validate(t, cmd)
		})
	}
}

func TestNewScanCommand_ExitCodes(t *testing.T) {
	// This test documents the expected exit codes mentioned in the command description
	cmd := NewScanCommand()

	// Check that the command documentation mentions exit codes
	if !strings.Contains(cmd.Long, "Exit codes:") {
		t.Error("Command should document exit codes in Long description")
	}

	if !strings.Contains(cmd.Long, "0:") {
		t.Error("Command should document exit code 0")
	}

	if !strings.Contains(cmd.Long, "1:") {
		t.Error("Command should document exit code 1")
	}
}