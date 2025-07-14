#!/bin/bash

# Demo script to show auto-completion functionality
# This script demonstrates the auto-completion features of kubeprobes

echo "=== kubeprobes Auto-completion Demonstration ==="
echo
echo "1. Loading bash completion..."
source <(./bin/kubeprobes completion bash)
echo "✓ Bash completion loaded"
echo

echo "2. Available main commands:"
echo "   Type: ./bin/kubeprobes <TAB>"
echo "   Result: completion  help  scan"
echo

echo "3. Available completion subcommands:"
echo "   Type: ./bin/kubeprobes completion <TAB>"
echo "   Result: bash  fish  powershell  zsh"
echo

echo "4. Help for specific commands:"
echo "   Type: ./bin/kubeprobes help <TAB>"
echo "   Result: completion  scan"
echo

echo "5. Flag completion examples:"
echo "   Type: ./bin/kubeprobes scan --<TAB>"
echo "   Shows available flags for scan command"
echo

echo "6. Completion command help:"
echo "   Type: ./bin/kubeprobes completion --<TAB>"
echo "   Result: --help  -h"
echo

echo "=== Verification Commands ==="
echo
echo "To verify completion works:"
echo "  source <(./bin/kubeprobes completion bash)"
echo "  ./bin/kubeprobes <TAB>  # Should show: completion help scan"
echo "  ./bin/kubeprobes completion <TAB>  # Should show: bash fish powershell zsh"
echo

echo "=== Installation Instructions ==="
echo
echo "Bash:"
echo "  echo 'source <(kubeprobes completion bash)' >> ~/.bashrc"
echo
echo "Zsh:"
echo "  echo 'source <(kubeprobes completion zsh)' >> ~/.zshrc"
echo
echo "Fish:"
echo "  kubeprobes completion fish > ~/.config/fish/completions/kubeprobes.fish"
echo
echo "PowerShell:"
echo "  echo 'kubeprobes completion powershell | Out-String | Invoke-Expression' >> \$PROFILE"
echo

echo "=== Script Sizes ==="
echo "Bash completion: $(./bin/kubeprobes completion bash | wc -l) lines"
echo "Zsh completion: $(./bin/kubeprobes completion zsh | wc -l) lines"  
echo "Fish completion: $(./bin/kubeprobes completion fish | wc -l) lines"
echo "PowerShell completion: $(./bin/kubeprobes completion powershell | wc -l) lines"

echo
echo "Demo completed!"