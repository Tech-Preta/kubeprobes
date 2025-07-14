package scanner

import (
	"strings"
	"testing"
)

func TestNewScanCommand(t *testing.T) {
	cmd := NewScanCommand()

	// Test that the command is properly configured
	if cmd.Use != "scan" {
		t.Errorf("Expected command use to be 'scan', got '%s'", cmd.Use)
	}

	if cmd.Short != "Scan Kubernetes workloads for probes" {
		t.Errorf("Expected short description to be 'Scan Kubernetes workloads for probes', got '%s'", cmd.Short)
	}

	// Test that all required flags are present
	expectedFlags := []string{
		"kubeconfig", "kubeContext", "namespace", "probe-type", 
		"recommendation", "all-namespaces", "output", "fail-on-warn",
	}

	for _, flagName := range expectedFlags {
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to be present", flagName)
		}
	}
}

func TestFlagDescriptions(t *testing.T) {
	cmd := NewScanCommand()

	testCases := []struct {
		flagName    string
		shouldContain string
	}{
		{"kubeconfig", "Path to the kubeconfig file"},
		{"kubeContext", "Kubernetes context to use"},
		{"namespace", "Kubernetes namespace to scan"},
		{"all-namespaces", "Scan all namespaces"},
		{"probe-type", "Type of probe to scan for"},
		{"recommendation", "Show detailed recommendations"},
		{"output", "Output format: text, json, or yaml"},
		{"fail-on-warn", "Exit with code 1 if warnings are found"},
	}

	for _, tc := range testCases {
		flag := cmd.Flags().Lookup(tc.flagName)
		if flag == nil {
			t.Errorf("Flag '%s' not found", tc.flagName)
			continue
		}

		if !strings.Contains(flag.Usage, tc.shouldContain) {
			t.Errorf("Flag '%s' description should contain '%s', but got '%s'", 
				tc.flagName, tc.shouldContain, flag.Usage)
		}
	}
}

func TestFlagDefaults(t *testing.T) {
	cmd := NewScanCommand()

	testCases := []struct {
		flagName     string
		expectedDefault string
	}{
		{"namespace", "default"},
		{"output", "text"},
		{"all-namespaces", "false"},
		{"fail-on-warn", "false"},
		{"recommendation", "false"},
	}

	for _, tc := range testCases {
		flag := cmd.Flags().Lookup(tc.flagName)
		if flag == nil {
			t.Errorf("Flag '%s' not found", tc.flagName)
			continue
		}

		if flag.DefValue != tc.expectedDefault {
			t.Errorf("Flag '%s' should have default value '%s', but got '%s'", 
				tc.flagName, tc.expectedDefault, flag.DefValue)
		}
	}
}

func TestOutputFormatValidation(t *testing.T) {
	cmd := NewScanCommand()
	
	// Set required flags to avoid other validation errors
	cmd.Flags().Set("kubeconfig", "/fake/path")
	
	// Test valid output formats - these should not error on the format validation
	validFormats := []string{"text", "json", "yaml"}
	for _, format := range validFormats {
		cmd.Flags().Set("output", format)
		// We expect this to fail on kubernetes client creation, not format validation
		err := cmd.RunE(cmd, []string{})
		if err != nil && strings.Contains(err.Error(), "invalid output format") {
			t.Errorf("Valid format '%s' was rejected", format)
		}
	}
	
	// Test invalid output format
	cmd.Flags().Set("output", "invalid")
	err := cmd.RunE(cmd, []string{})
	if err == nil || !strings.Contains(err.Error(), "invalid output format") {
		t.Errorf("Invalid output format should be rejected, got error: %v", err)
	}
}

func TestCompletionFunctions(t *testing.T) {
	cmd := NewScanCommand()

	// Test probe-type completion
	probeTypeFlag := cmd.Flags().Lookup("probe-type")
	if probeTypeFlag == nil {
		t.Fatal("probe-type flag not found")
	}

	// Test output completion
	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Fatal("output flag not found")
	}

	// The completion functions are registered, but we can't easily test them here
	// without more complex setup. The important thing is that the flags exist
	// and the command builds successfully.
}