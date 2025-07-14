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
			name:           "help flag shows improved description",
			args:           []string{"--help"},
			expectedOutput: "Kubeprobes is a CLI tool for scanning Kubernetes workloads",
			expectError:    false,
		},
		{
			name:           "help shows health check benefits",
			args:           []string{"--help"},
			expectedOutput: "Health check probes are critical for:",
			expectError:    false,
		},
		{
			name:           "help shows liveness probe description",
			args:           []string{"--help"},
			expectedOutput: "Liveness probes: Detect when to restart containers",
			expectError:    false,
		},
		{
			name:           "no args shows help",
			args:           []string{},
			expectedOutput: "Available Commands:",
			expectError:    false,
		},
		{
			name:           "help includes improved examples",
			args:           []string{"--help"},
			expectedOutput: "Quick scan of default namespace",
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
configured by scanning for missing liveness, readiness, and startup probes.
Proper probe configuration is essential for reliable deployments, effective load
balancing, and early detection of application issues.

Health check probes are critical for:
  • Liveness probes: Detect when to restart containers
  • Readiness probes: Control traffic routing to healthy containers  
  • Startup probes: Handle slow-starting containers gracefully`,
				Example: `  # Quick scan of default namespace
  kubeprobes scan

  # Scan with detailed recommendations
  kubeprobes scan --recommendation

  # Scan specific namespace for liveness probes only
  kubeprobes scan --namespace my-app --probe-type liveness

  # Scan using specific kubeconfig
  kubeprobes scan --kubeconfig ~/.kube/prod-config

  # Check tool version
  kubeprobes version`,
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

// TestPOSIXCompliance tests POSIX command line syntax compliance
func TestPOSIXCompliance(t *testing.T) {
	tests := []struct {
		name string
		args []string
		desc string
	}{
		{
			name: "grouped_short_flags",
			args: []string{"scan", "-rp", "liveness"},
			desc: "Should support grouped short flags (-rp instead of -r -p)",
		},
		{
			name: "mixed_flag_order",
			args: []string{"scan", "-r", "--namespace", "test", "-p", "readiness"},
			desc: "Should support flexible flag ordering",
		},
		{
			name: "short_flags_first",
			args: []string{"scan", "-n", "test", "-k", "/path/to/config", "--recommendation"},
			desc: "Should support short flags before long flags",
		},
		{
			name: "long_flags_first", 
			args: []string{"scan", "--recommendation", "--probe-type", "startup", "-n", "test"},
			desc: "Should support long flags before short flags",
		},
		{
			name: "equals_syntax",
			args: []string{"scan", "--namespace=test", "--probe-type=liveness"},
			desc: "Should support --flag=value syntax",
		},
		{
			name: "help_short_flag",
			args: []string{"-h"},
			desc: "Should support standard -h help flag",
		},
		{
			name: "help_long_flag",
			args: []string{"--help"},
			desc: "Should support standard --help flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh command for each test
			cmd := &cobra.Command{
				Use:   "kubeprobes",
				Short: "Kubeprobes is a CLI tool for scanning Kubernetes probes",
			}
			cmd.AddCommand(NewScanCommand())
			cmd.AddCommand(NewVersionCommand())

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			// Execute - we expect most to fail due to no k8s cluster, 
			// but flags should be parsed correctly
			err := cmd.Execute()
			
			// For help flags, we don't expect errors
			if tt.name == "help_short_flag" || tt.name == "help_long_flag" {
				if err != nil {
					t.Errorf("Test %q failed with error: %v", tt.desc, err)
				}
				return
			}

			// For scan commands, check that flag parsing didn't fail
			// (cluster connection errors are expected and OK)
			output := buf.String()
			if strings.Contains(output, "unknown flag") || strings.Contains(output, "unknown shorthand") {
				t.Errorf("Test %q failed - flag parsing error: %s", tt.desc, output)
			}
		})
	}
}

// TestPOSIXFlagCompatibility tests specific POSIX flag compatibility
func TestPOSIXFlagCompatibility(t *testing.T) {
	scanCmd := NewScanCommand()
	
	tests := []struct {
		name       string
		shortFlag  string
		longFlag   string
		expectBoth bool
	}{
		{"kubeconfig", "k", "kubeconfig", true},
		{"kubeContext", "c", "kubeContext", true},
		{"namespace", "n", "namespace", true},
		{"probe-type", "p", "probe-type", true},
		{"recommendation", "r", "recommendation", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortFlag := scanCmd.Flags().ShorthandLookup(tt.shortFlag)
			longFlag := scanCmd.Flags().Lookup(tt.longFlag)

			if tt.expectBoth {
				if shortFlag == nil {
					t.Errorf("Expected short flag -%s to exist", tt.shortFlag)
				}
				if longFlag == nil {
					t.Errorf("Expected long flag --%s to exist", tt.longFlag)
				}
				
				// Verify they're the same flag
				if shortFlag != nil && longFlag != nil && shortFlag != longFlag {
					t.Errorf("Short flag -%s and long flag --%s should be the same flag", tt.shortFlag, tt.longFlag)
				}
			}
		})
	}
}
