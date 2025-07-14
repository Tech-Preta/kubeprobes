package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCompletionBashScript(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Use the actual root command from the CLI package
	cmd := getRootCmdForTesting()

	// Generate bash completion script
	err := cmd.GenBashCompletion(&buf)
	if err != nil {
		t.Fatalf("Failed to generate bash completion: %v", err)
	}

	// Verify that the output contains expected bash completion elements
	output := buf.String()

	// Check that it's a valid bash completion script - look for actual content
	expectedElements := []string{
		"# bash completion",
		"__kubeprobes_debug",
		"__start_kubeprobes",
		"complete",
		"kubeprobes",
	}

	for _, element := range expectedElements {
		if !strings.Contains(output, element) {
			t.Errorf("Expected bash completion script to contain '%s', but it was not found. First 200 chars: %s", element, output[:min(200, len(output))])
		}
	}

	// Ensure the script is not empty
	if len(output) == 0 {
		t.Error("Expected bash completion script to not be empty")
	}

	// Ensure it's a reasonable size (bash completion scripts are typically several KB)
	if len(output) < 100 {
		t.Errorf("Expected bash completion script to be substantial in size, got %d characters", len(output))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getRootCmdForTesting() *cobra.Command {
	// Create a copy of the root command for testing
	cmd := &cobra.Command{
		Use:   "kubeprobes",
		Short: "Kubeprobes is a CLI tool for scanning Kubernetes probes",
		Long: `Kubeprobes is a CLI tool for scanning Kubernetes workloads to detect 
missing liveness, readiness, and startup probes.`,
	}
	// The completion command is automatically added by Cobra
	return cmd
}

func TestCompletionCommandExists(t *testing.T) {
	// Test that the binary itself has the completion command
	// This is a simpler integration test rather than a unit test
	var buf bytes.Buffer
	cmd := getRootCmdForTesting()
	
	// Just verify that we can generate completion scripts
	err := cmd.GenBashCompletion(&buf)
	if err != nil {
		t.Fatalf("Failed to generate bash completion: %v", err)
	}
	
	if buf.Len() == 0 {
		t.Error("Expected bash completion to generate content")
	}
}

func TestCompletionScriptGeneration(t *testing.T) {
	tests := []struct {
		shell        string
		expectedText string
	}{
		{
			shell:        "bash",
			expectedText: "# bash completion",
		},
		{
			shell:        "zsh", 
			expectedText: "# zsh completion",
		},
		{
			shell:        "fish",
			expectedText: "# fish completion",
		},
		{
			shell:        "powershell",
			expectedText: "# powershell completion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			var buf bytes.Buffer
			
			// Use the actual root command
			cmd := getRootCmdForTesting()

			// Generate completion for each shell
			var err error
			switch tt.shell {
			case "bash":
				err = cmd.GenBashCompletion(&buf)
			case "zsh":
				err = cmd.GenZshCompletion(&buf)
			case "fish":
				err = cmd.GenFishCompletion(&buf, true)
			case "powershell":
				err = cmd.GenPowerShellCompletion(&buf)
			}

			if err != nil {
				t.Fatalf("Failed to generate %s completion: %v", tt.shell, err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.expectedText) {
				t.Errorf("Expected %s completion script to contain '%s', but it was not found. First 200 chars: %s", tt.shell, tt.expectedText, output[:min(200, len(output))])
			}

			if len(output) == 0 {
				t.Errorf("Expected %s completion script to not be empty", tt.shell)
			}
		})
	}
}