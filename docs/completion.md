# Auto-completion

O kubeprobes oferece suporte completo a auto-completion para bash, zsh, fish e PowerShell utilizando o framework Cobra.

## Como funciona

O auto-completion permite que você use a tecla Tab para:
- Completar comandos e subcomandos
- Completar nomes de flags
- Obter sugestões de valores para algumas flags

## Configuração por Shell

### Bash

#### Pré-requisitos
- Pacote `bash-completion` instalado no sistema

#### Instalação temporária (sessão atual)
```bash
source <(kubeprobes completion bash)
```

#### Instalação permanente

**Linux:**
```bash
kubeprobes completion bash > /etc/bash_completion.d/kubeprobes
```

**macOS (usando Homebrew):**
```bash
kubeprobes completion bash > $(brew --prefix)/etc/bash_completion.d/kubeprobes
```

### Zsh

#### Pré-requisitos
Certifique-se de que o auto-completion esteja habilitado no zsh:
```bash
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

#### Instalação temporária (sessão atual)
```bash
source <(kubeprobes completion zsh)
```

#### Instalação permanente

**Linux:**
```bash
kubeprobes completion zsh > "${fpath[1]}/_kubeprobes"
```

**macOS (usando Homebrew):**
```bash
kubeprobes completion zsh > $(brew --prefix)/share/zsh/site-functions/_kubeprobes
```

### Fish

#### Instalação temporária (sessão atual)
```bash
kubeprobes completion fish | source
```

#### Instalação permanente
```bash
kubeprobes completion fish > ~/.config/fish/completions/kubeprobes.fish
```

### PowerShell

#### Instalação temporária (sessão atual)
```powershell
kubeprobes completion powershell | Out-String | Invoke-Expression
```

#### Instalação permanente
1. Execute o comando acima para obter o script de completion
2. Adicione a saída ao seu perfil do PowerShell
3. Para encontrar o local do perfil: `$PROFILE`

## Recursos de Auto-completion

### Comandos
- `kubeprobes` + Tab: mostra subcomandos disponíveis (scan, completion, help)
- `kubeprobes scan` + Tab: mostra flags disponíveis

### Flags
- `--` + Tab: lista todas as flags disponíveis
- `--probe-type` + Tab: sugere valores válidos (liveness, readiness, startup)
- `--namespace` + Tab: sugere namespaces do cluster (quando configurado)

### Exemplos de uso

```bash
# Completar comando
kubeprobes [Tab]
# Resultado: completion help scan

# Completar flags
kubeprobes scan --[Tab]
# Resultado: --help --kubeconfig --kubeContext --namespace --probe-type --recommendation

# Completar valores de probe-type
kubeprobes scan --probe-type [Tab]
# Resultado: liveness readiness startup
```

## Solução de problemas

### Auto-completion não funciona

1. **Verifique se o kubeprobes está no PATH:**
   ```bash
   which kubeprobes
   ```

2. **Verifique se o shell suporta completion:**
   - Bash: verifique se `bash-completion` está instalado
   - Zsh: verifique se `compinit` está carregado
   - Fish: deve funcionar nativamente
   - PowerShell: deve funcionar nativamente

3. **Recarregue o shell:**
   ```bash
   # Bash/Zsh
   source ~/.bashrc  # ou ~/.zshrc
   
   # Fish
   source ~/.config/fish/config.fish
   ```

4. **Verifique se o script foi instalado corretamente:**
   ```bash
   # Bash
   ls /etc/bash_completion.d/kubeprobes
   
   # Zsh
   ls "${fpath[1]}/_kubeprobes"
   
   # Fish
   ls ~/.config/fish/completions/kubeprobes.fish
   ```

### Auto-completion parcial

Se apenas alguns recursos de completion funcionam, pode ser necessário:
- Atualizar para uma versão mais recente do shell
- Reinstalar o script de completion
- Verificar se não há conflitos com outros scripts de completion

## Desabilitando descrições

Por padrão, o auto-completion inclui descrições dos comandos e flags. Para desabilitá-las:

```bash
# Bash
kubeprobes completion bash --no-descriptions > /etc/bash_completion.d/kubeprobes

# Zsh
kubeprobes completion zsh --no-descriptions > "${fpath[1]}/_kubeprobes"

# Fish
kubeprobes completion fish --no-descriptions > ~/.config/fish/completions/kubeprobes.fish

# PowerShell
kubeprobes completion powershell --no-descriptions | Out-String | Invoke-Expression
```