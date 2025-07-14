# Auto-completion for kubeprobes

The `kubeprobes` CLI tool provides auto-completion support for Bash, Zsh, Fish, and PowerShell shells. This feature helps you quickly discover and complete commands, flags, and their values when using the tool interactively.

## Generating Completion Scripts

You can generate completion scripts for your preferred shell using the `completion` command:

### Bash

```bash
# Generate completion script
./kubeprobes completion bash

# To enable completion for the current session:
source <(./kubeprobes completion bash)

# To enable completion permanently, add to your shell profile:
echo 'source <(kubeprobes completion bash)' >> ~/.bashrc
```

### Zsh

```zsh
# Generate completion script
./kubeprobes completion zsh

# To enable completion for the current session:
source <(./kubeprobes completion zsh)

# To enable completion permanently:
echo 'source <(kubeprobes completion zsh)' >> ~/.zshrc

# For Oh My Zsh users, you can also create a completion file:
kubeprobes completion zsh > "${fpath[1]}/_kubeprobes"
```

### Fish

```fish
# Generate completion script
./kubeprobes completion fish

# To enable completion permanently:
kubeprobes completion fish | source

# Or save to Fish's completion directory:
kubeprobes completion fish > ~/.config/fish/completions/kubeprobes.fish
```

### PowerShell

```powershell
# Generate completion script
./kubeprobes completion powershell

# To enable completion for the current session:
kubeprobes completion powershell | Out-String | Invoke-Expression

# To enable completion permanently, add to your PowerShell profile:
echo 'kubeprobes completion powershell | Out-String | Invoke-Expression' >> $PROFILE
```

## Usage Examples

Once auto-completion is enabled, you can use the TAB key to complete commands and options:

```bash
# Complete commands
./kubeprobes <TAB>
# Shows: completion  help  scan

# Complete completion subcommands
./kubeprobes completion <TAB>
# Shows: bash  fish  powershell  zsh

# Complete scan command options
./kubeprobes scan <TAB>
# Shows available flags and options for the scan command
```

## Supported Features

The auto-completion supports:

- **Command completion**: Complete main commands (`scan`, `completion`, `help`)
- **Subcommand completion**: Complete subcommands under `completion`
- **Flag completion**: Complete long and short flags for commands
- **Context-aware suggestions**: Provides relevant options based on the current command context

## Verification

To verify that auto-completion is working correctly:

1. Install the completion script for your shell (see instructions above)
2. Open a new terminal session
3. Type `./kubeprobes <TAB>` and verify that commands are suggested
4. Type `./kubeprobes completion <TAB>` and verify that shell options are suggested

## Troubleshooting

### Bash
- Ensure you have `bash-completion` package installed
- Verify that your `.bashrc` sources the completion script
- Check that the completion script is valid: `bash -n <(kubeprobes completion bash)`

### Zsh
- Make sure the completion function is in your `fpath`
- Verify that `compinit` is called in your `.zshrc`
- Check for syntax errors: `zsh -n <(kubeprobes completion zsh)`

### Fish
- Verify that the completion file is in the correct directory
- Check Fish's completion path: `echo $fish_complete_path`
- Test the completion file: `fish -c "source <(kubeprobes completion fish)"`

### PowerShell
- Ensure your execution policy allows script execution
- Verify the completion script loads without errors
- Check your PowerShell profile exists: `Test-Path $PROFILE`

For more information about shell completion, see the documentation for your specific shell.