package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is set during build time
	Version = "dev"
	// Commit is set during build time
	Commit = "unknown"
	// Date is set during build time
	Date = "unknown"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long: `Print the version information of kubeprobes.

This command displays the version, commit hash, build date, and Go version
used to build this binary. This information is useful for troubleshooting
and ensuring you're running the expected version.`,
		Example: `  # Show detailed version information  
  kubeprobes version

  # Show only version number (useful for scripts)
  kubeprobes version --output=short

  # POSIX syntax: short flag
  kubeprobes version -o short

  # Show version information as JSON
  kubeprobes version --output=json

  # POSIX syntax: short flag with equals
  kubeprobes version -o=json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				return fmt.Errorf("failed to read output flag: %w", err)
			}

			switch output {
			case "short":
				cmd.Print(Version)
				cmd.Print("\n")
			case "json":
				cmd.Printf(`{"version":"%s","commit":"%s","date":"%s","goVersion":"%s"}%s`,
					Version, Commit, Date, runtime.Version(), "\n")
			default:
				cmd.Printf("kubeprobes version %s\n", Version)
				cmd.Printf("Commit: %s\n", Commit)
				cmd.Printf("Date: %s\n", Date)
				cmd.Printf("Go version: %s\n", runtime.Version())
			}
			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "default", "Output format: default, short, or json")

	// Add custom completion for output flag
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"default", "short", "json"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		// Log error but don't fail the command creation
		fmt.Printf("Warning: failed to register completion for output flag: %v\n", err)
	}

	return cmd
}
