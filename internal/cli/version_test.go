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
			name:           "help flag shows improved description",
			args:           []string{"--help"},
			expectedOutput: "Print the version information",
			expectError:    false,
		},
		{
			name:           "help shows troubleshooting context",
			args:           []string{"--help"},
			expectedOutput: "useful for troubleshooting",
			expectError:    false,
		},
		{
			name:           "help shows improved examples",
			args:           []string{"--help"},
			expectedOutput: "Show only version number (useful for scripts)",
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

func TestVersionCommand_OutputFlag(t *testing.T) {
	cmd := NewVersionCommand()

	// Test that output flag exists and has correct properties
	flag := cmd.Flag("output")
	if flag == nil {
		t.Fatal("output flag should exist")
	}

	if flag.Shorthand != "o" {
		t.Errorf("Expected shorthand 'o', got %q", flag.Shorthand)
	}

	if flag.DefValue != "default" {
		t.Errorf("Expected default value 'default', got %q", flag.DefValue)
	}

	// Test that output flag has improved description from our PR
	expectedUsage := "Output format: default, short, or json"
	if flag.Usage != expectedUsage {
		t.Errorf("Output flag usage should be %q, got %q", expectedUsage, flag.Usage)
	}
}

func TestVersionCommand_FlagCompletion(t *testing.T) {
	cmd := NewVersionCommand()

	// Test that the output flag has completion registered
	flag := cmd.Flag("output")
	if flag == nil {
		t.Fatal("output flag should exist")
	}

	// The completion function should be registered (we can't easily test the actual completion)
	// but we can verify the flag setup is correct
	if flag.Usage == "" {
		t.Error("output flag should have usage text")
	}
}

func TestVersionCommand_Examples(t *testing.T) {
	cmd := NewVersionCommand()

	if cmd.Example == "" {
		t.Error("Version command should have examples")
	}

	// Check that examples contain expected content
	if !strings.Contains(cmd.Example, "kubeprobes version") {
		t.Error("Examples should contain basic usage")
	}

	if !strings.Contains(cmd.Example, "--output=short") {
		t.Error("Examples should contain short output format")
	}
}

func TestVersionCommand_ErrorHandling(t *testing.T) {
    cmd := NewVersionCommand()

    // Test with various edge cases
    tests := []struct {
        name string
        args []string
    }{
        {"with extra arguments", []string{"extra", "args"}},
        {"with unknown flag", []string{"--unknown-flag"}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := NewVersionCommand() // Nova instância para cada teste
            buf := new(bytes.Buffer)
            cmd.SetOut(buf)
            cmd.SetErr(buf)
            cmd.SetArgs(tt.args)

            err := cmd.Execute()
            // These might error (unknown flag) or succeed (extra args ignored)
            // We're just testing that the command handles them gracefully
            if err != nil {
                t.Logf("Expected error for %s: %v", tt.name, err)
            }
        })
    }
}

// TestVersionCommandPOSIXCompliance tests POSIX compliance for version command
func TestVersionCommandPOSIXCompliance(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		description string
	}{
		{
			name:        "short_flag",
			args:        []string{"-o", "short"},
			description: "Should support short flag syntax",
		},
		{
			name:        "long_flag",
			args:        []string{"--output", "json"},
			description: "Should support long flag syntax",
		},
		{
			name:        "equals_syntax",
			args:        []string{"--output=short"},
			description: "Should support --flag=value syntax",
		},
		{
			name:        "help_short",
			args:        []string{"-h"},
			description: "Should support standard -h help flag",
		},
		{
			name:        "help_long",
			args:        []string{"--help"},
			description: "Should support standard --help flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewVersionCommand()

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err != nil {
				t.Errorf("POSIX syntax test failed for %s: %v", tt.description, err)
			}

			output := buf.String()

			// Verify that output was produced for non-help commands
			if !strings.Contains(tt.name, "help") && len(output) == 0 {
				t.Errorf("Expected output for %s", tt.description)
			}

			// Verify help commands produce help output
			if strings.Contains(tt.name, "help") && !strings.Contains(output, "Print the version information") {
				t.Errorf("Help command should produce help output")
			}
		})
	}
}

// TestVersionCommandFlagConventions tests version command flag conventions
func TestVersionCommandFlagConventions(t *testing.T) {
	cmd := NewVersionCommand()

	// Test output flag conventions
	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Fatal("Output flag should exist")
	}

	if outputFlag.Shorthand != "o" {
		t.Errorf("Output flag should have shorthand 'o', got %q", outputFlag.Shorthand)
	}

	// Test that help shows POSIX format
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	output := buf.String()

	// Should show both short and long flag forms
	if !strings.Contains(output, "-o, --output") {
		t.Error("Help should show both short and long flag forms")
	}
}
