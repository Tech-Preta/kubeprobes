# Terminal Session Demonstration

This file shows a realistic terminal session demonstrating the auto-completion functionality.

## Session Setup

```bash
$ cd kubeprobes
$ make build
Building kubeprobes...

$ source <(./bin/kubeprobes completion bash)
# Completion is now active for this session
```

## Command Completion Demonstrations

### 1. Main Command Completion
```bash
$ ./bin/kubeprobes <TAB>
completion  help        scan

$ ./bin/kubeprobes h<TAB>
help

$ ./bin/kubeprobes help <TAB>
completion  scan
```

### 2. Completion Subcommand Completion
```bash
$ ./bin/kubeprobes completion <TAB>
bash        fish        powershell  zsh

$ ./bin/kubeprobes completion b<TAB>
bash

$ ./bin/kubeprobes completion z<TAB>
zsh
```

### 3. Flag Completion
```bash
$ ./bin/kubeprobes --<TAB>
--help  -h

$ ./bin/kubeprobes completion --<TAB>
--help  -h

$ ./bin/kubeprobes scan --<TAB>
--help  -h
```

### 4. Completion Command Help Verification
```bash
$ ./bin/kubeprobes completion --help
Generate the autocompletion script for kubeprobes for the specified shell.
See each sub-command's help for details on how to use the generated script.

Usage:
  kubeprobes completion [command]

Available Commands:
  bash        Generate the autocompletion script for bash
  fish        Generate the autocompletion script for fish
  powershell  Generate the autocompletion script for powershell
  zsh         Generate the autocompletion script for zsh

Flags:
  -h, --help   help for completion

Use "kubeprobes completion [command] --help" for more information about a command.
```

### 5. Individual Shell Completion Help
```bash
$ ./bin/kubeprobes completion bash --help
Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(kubeprobes completion bash)

To load completions for every new session, execute once:

#### Linux:

	kubeprobes completion bash > /etc/bash_completion.d/kubeprobes

#### macOS:

	kubeprobes completion bash > $(brew --prefix)/etc/bash_completion.d/kubeprobes

You will need to start a new shell for this setup to take effect.

Usage:
  kubeprobes completion bash

Flags:
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

## Real-world Usage Examples

### Installing completion permanently
```bash
# For Bash
$ echo 'source <(kubeprobes completion bash)' >> ~/.bashrc
$ source ~/.bashrc

# Test that it works in new session
$ kubeprobes <TAB>
completion  help        scan
```

### Testing completion with different shells
```bash
# Generate and test zsh completion
$ ./bin/kubeprobes completion zsh > /tmp/kubeprobes_zsh_completion
$ zsh -c "source /tmp/kubeprobes_zsh_completion && echo 'Zsh completion loaded successfully'"
Zsh completion loaded successfully

# Generate and test fish completion  
$ ./bin/kubeprobes completion fish > /tmp/kubeprobes_fish_completion
$ fish -c "source /tmp/kubeprobes_fish_completion && echo 'Fish completion loaded successfully'"
Fish completion loaded successfully
```

### Completion script validation
```bash
$ bash -n <(./bin/kubeprobes completion bash) && echo "Bash completion syntax OK"
Bash completion syntax OK

$ zsh -n <(./bin/kubeprobes completion zsh) && echo "Zsh completion syntax OK"  
Zsh completion syntax OK

$ fish -n <(./bin/kubeprobes completion fish) && echo "Fish completion syntax OK"
Fish completion syntax OK
```

## Performance and Size Information

```bash
$ time ./bin/kubeprobes completion bash > /dev/null
real    0m0.015s
user    0m0.008s
sys     0m0.004s

$ wc -l examples/completion/*.sh examples/completion/*.zsh examples/completion/*.fish examples/completion/*.ps1
    338 examples/completion/kubeprobes_bash_completion.sh
    212 examples/completion/kubeprobes_zsh_completion.zsh
    235 examples/completion/kubeprobes_fish_completion.fish
    245 examples/completion/kubeprobes_powershell_completion.ps1
   1030 total
```

All completion scripts are generated quickly and provide comprehensive auto-completion functionality for the kubeprobes CLI tool.