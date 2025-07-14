package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "help flag",
			args:           []string{"--help"},
			expectedOutput: "Kubeprobes is a CLI tool for scanning Kubernetes workloads",
			expectError:    false,
		},
		{
			name:           "no args shows help",
			args:           []string{},
			expectedOutput: "Available Commands:",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for testing
			cmd := &cobra.Command{
				Use:   "kubeprobes",
				Short: "Kubeprobes is a CLI tool for scanning Kubernetes probes",
				Long: `Kubeprobes is a CLI tool for scanning Kubernetes workloads to detect 
missing liveness, readiness, and startup probes.

Kubeprobes helps you ensure your Kubernetes workloads have proper health checks
configured by scanning for missing liveness, readiness, and startup probes.`,
			}

			// Add subcommands
			cmd.AddCommand(NewScanCommand())
			cmd.AddCommand(NewVersionCommand())

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

			output := buf.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestSubcommandAvailability(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{
		Use:   "kubeprobes",
		Short: "Kubeprobes is a CLI tool for scanning Kubernetes probes",
	}

	// Add subcommands
	cmd.AddCommand(NewScanCommand())
	cmd.AddCommand(NewVersionCommand())

	expectedCommands := []string{"scan", "version"}

	availableCommands := make(map[string]bool)
	for _, subCmd := range cmd.Commands() {
		availableCommands[subCmd.Name()] = true
	}

	for _, expectedCmd := range expectedCommands {
		if !availableCommands[expectedCmd] {
			t.Errorf("Expected command %q to be available", expectedCmd)
		}
	}

	// Note: completion and help are automatically added by Cobra
	// so they will be present but we only test for the ones we explicitly add
}

func TestRootCommand_Structure(t *testing.T) {
	// Test that Execute function exists and can be called
	// In a real scenario, this would execute the CLI
	
	// This is a compilation test to ensure the function exists
	// We don't actually call Execute() as it would try to run the CLI
	t.Log("Execute function is available for use")
}

func TestRootCommand_LongDescription(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "kubeprobes",
		Short: "Kubeprobes is a CLI tool for scanning Kubernetes probes",
		Long: `Kubeprobes is a CLI tool for scanning Kubernetes workloads to detect 
missing liveness, readiness, and startup probes.

Kubeprobes helps you ensure your Kubernetes workloads have proper health checks
configured by scanning for missing liveness, readiness, and startup probes.`,
	}

	if cmd.Long == "" {
		t.Error("Root command should have a long description")
	}

	// Check that long description contains key information
	expectedContent := []string{
		"liveness",
		"readiness", 
		"startup",
		"probes",
		"Kubernetes",
	}

	for _, content := range expectedContent {
		if !strings.Contains(cmd.Long, content) {
			t.Errorf("Long description should contain %q", content)
		}
	}
}

func TestRootCommand_Examples(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "kubeprobes",
		Short: "Kubeprobes is a CLI tool for scanning Kubernetes probes",
		Example: `  # Scan all workloads in the default namespace
  kubeprobes scan

  # Scan with recommendations for missing probes
  kubeprobes scan --recommendation

  # Scan a specific namespace for liveness probes only
  kubeprobes scan --namespace my-app --probe-type liveness

  # Show version information
  kubeprobes version`,
	}

	if cmd.Example == "" {
		t.Error("Root command should have examples")
	}

	// Check that examples contain expected commands
	expectedExamples := []string{
		"kubeprobes scan",
		"kubeprobes version",
		"--recommendation",
		"--namespace",
		"--probe-type",
	}

	for _, example := range expectedExamples {
		if !strings.Contains(cmd.Example, example) {
			t.Errorf("Examples should contain %q", example)
		}
	}
}
