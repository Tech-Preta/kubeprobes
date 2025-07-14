# Reestruturação do Projeto Kubeprobes

Este documento descreve a reestruturação do projeto Kubeprobes para seguir as convenções recomendadas do Go.

## Estrutura Anterior vs Nova Estrutura

### Estrutura Anterior
```
kubeprobes/
├── src/
│   ├── go.mod
│   ├── go.sum
│   └── probes.go
├── samples/
│   ├── ns.yaml
│   └── probes.yaml
└── [outros arquivos]
```

### Nova Estrutura (Go Layout Padrão)
```
kubeprobes/
├── cmd/
│   └── kubeprobes/
│       └── main.go
├── internal/
│   ├── cli/
│   │   └── root.go
│   └── scanner/
│       ├── cmd.go
│       └── scanner.go
├── pkg/
│   └── kubernetes/
│       └── client.go
├── examples/
│   ├── ns.yaml
│   └── probes.yaml
├── .vscode/
│   ├── launch.json
│   ├── settings.json
│   └── tasks.json
├── go.mod
├── go.sum
├── Makefile
└── [outros arquivos]
```

## Principais Mudanças

### 1. Diretório `/cmd`
- **Antes**: Código principal em `src/probes.go`
- **Agora**: Aplicação principal em `cmd/kubeprobes/main.go`
- **Benefício**: Segue a convenção Go para aplicações executáveis

### 2. Diretório `/internal`
- **Novo**: Código privado da aplicação em `internal/`
- **Conteúdo**: 
  - `internal/cli/`: Comandos CLI
  - `internal/scanner/`: Lógica de escaneamento
- **Benefício**: Código não pode ser importado por outros projetos

### 3. Diretório `/pkg`
- **Novo**: Código reutilizável em `pkg/`
- **Conteúdo**: `pkg/kubernetes/`: Cliente Kubernetes
- **Benefício**: Código que pode ser importado por outros projetos

### 4. Diretório `/examples`
- **Antes**: `samples/`
- **Agora**: `examples/`
- **Benefício**: Nome mais descritivo e padrão

### 5. Configurações de Desenvolvimento
- **Novo**: `.vscode/` com configurações para VS Code
- **Novo**: `Makefile` para automação de build
- **Benefício**: Melhor experiência de desenvolvimento

## Estrutura Modular

### main.go
```go
package main

import (
    "log"
    "kubeprobes/internal/cli"
)

func main() {
    if err := cli.Execute(); err != nil {
        log.Fatalf("Error executing command: %s", err.Error())
    }
}
```

### Separação de Responsabilidades
- **CLI**: Gerenciamento de comandos e flags
- **Scanner**: Lógica de escaneamento de probes
- **Kubernetes**: Cliente e operações do Kubernetes

## DevContainer Melhorado

### Extensões Adicionadas
- `golang.go`: Suporte completo ao Go
- `ms-vscode.makefile-tools`: Suporte a Makefiles

### Configurações Go
- IntelliSense completo
- Debug integrado
- Formatação automática
- Linting com golangci-lint

## Comandos Disponíveis

### Makefile
```bash
make build      # Compilar o projeto
make build-all  # Compilar para múltiplas plataformas
make test       # Executar testes
make clean      # Limpar artefatos
make fmt        # Formatar código
make lint       # Linting do código
make install    # Instalar no sistema
make run        # Executar a aplicação
```

### Tarefas VS Code
- Build (Ctrl+Shift+P -> Tasks: Run Task -> build)
- Test
- Run scan
- Clean
- Format
- Lint

## Benefícios da Reestruturação

1. **Organização Padrão**: Segue as convenções da comunidade Go
2. **Modularidade**: Código separado por responsabilidade
3. **Reutilização**: Pacotes em `/pkg` podem ser importados
4. **Privacidade**: Código em `/internal` é privado
5. **Desenvolvimento**: Configurações otimizadas para VS Code
6. **Build**: Makefile para automação
7. **Escalabilidade**: Estrutura preparada para crescimento

## Próximos Passos

1. ✅ Reestruturação de diretórios
2. ✅ Separação modular do código
3. ✅ Configuração do DevContainer
4. ✅ Criação do Makefile
5. ✅ Configurações VS Code
6. ✅ Atualização para Go 1.24.5
7. ✅ Resolução de dependências Go
8. ✅ Build e testes funcionais
9. ⏳ Testes unitários
10. ⏳ CI/CD atualizado
11. ⏳ Documentação atualizada

## Notas Técnicas

### Compatibilidade Go
- ✅ Projeto atualizado para Go 1.24.5
- ✅ DevContainer configurado para instalar Go 1.24.5
- ✅ Dependências do Kubernetes atualizadas para v0.31.0
- ✅ Build funcionando corretamente

### Build
- Artefatos gerados em `bin/`
- Suporte a múltiplas plataformas
- Build otimizado para containers
