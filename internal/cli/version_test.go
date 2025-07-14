package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "default output",
			args:           []string{},
			expectedOutput: "kubeprobes version dev",
			expectError:    false,
		},
		{
			name:           "short output",
			args:           []string{"--output=short"},
			expectedOutput: "dev",
			expectError:    false,
		},
		{
			name:           "json output",
			args:           []string{"--output=json"},
			expectedOutput: `"version":"dev"`,
			expectError:    false,
		},
		{
			name:           "help flag",
			args:           []string{"--help"},
			expectedOutput: "Print the version information",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewVersionCommand()

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

func TestVersionCommandJsonOutput(t *testing.T) {
	cmd := NewVersionCommand()

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--output=json"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	
	// Verify it's valid JSON
	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v\nOutput was: %q", err, output)
	}

	// Check required fields
	requiredFields := []string{"version", "commit", "date", "goVersion"}
	for _, field := range requiredFields {
		if _, exists := result[field]; !exists {
			t.Errorf("JSON output missing required field: %s", field)
		}
	}
}

func TestVersionCommandInvalidOutput(t *testing.T) {
	cmd := NewVersionCommand()

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--output=invalid"})

	// This should not error, but should default to normal output
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "kubeprobes version") {
		t.Errorf("Expected default output format for invalid output flag")
	}
}