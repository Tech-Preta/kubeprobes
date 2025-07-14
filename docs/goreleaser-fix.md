# Soluções para problemas do GoReleaser

## Problema 1: Token GitHub não reconhecido

**Problema**: O GoReleaser estava procurando por `GITHUB_TOKEN` mas você estava usando `USER_TOKEN`.

**Solução**: O workflow do GitHub Actions já está configurado corretamente:

```yaml
env:
  GITHUB_TOKEN: ${{ secrets.USER_TOKEN }}
```

Isso mapeia sua secret `USER_TOKEN` para a variável de ambiente `GITHUB_TOKEN` que o GoReleaser espera.

## Problema 2: "does not contain a main function"

**Problema**: O GoReleaser não estava encontrando a função main porque não sabia onde procurar.

**Solução**: Adicionado o caminho correto no `.goreleaser.yml`:

```yaml
builds:
  - id: kubeprobes
    main: ./cmd/kubeprobes  # <- Esta linha foi adicionada
    binary: kubeprobes
```

## Melhorias adicionais implementadas

1. **Configuração completa do projeto**:
   - Adicionado `project_name: kubeprobes`
   - Configurado `id` para o build

2. **Configuração de release**:
   - Configurado owner e repositório corretos
   - Configurado template de nome para releases

3. **Otimizações do workflow**:
   - Adicionado `cache: true` para melhor performance
   - Adicionado permissão `packages: write`

4. **Build flags**:
   - Adicionado ldflags para otimização binária
   - Configurado para incluir versão, commit e data

## Testando localmente

Para testar se a configuração está funcionando:

```bash
# Instalar goreleaser (se não tiver)
go install github.com/goreleaser/goreleaser@latest

# Testar build sem release
goreleaser build --snapshot --clean

# Testar release completo (sem publicar)
goreleaser release --snapshot --clean
```

## Configuração necessária no GitHub

1. Certifique-se de que a secret `USER_TOKEN` está configurada no repositório
2. O token deve ter as seguintes permissões:
   - `contents:write` (para criar releases)
   - `packages:write` (se for publicar packages)

## Próximos passos

Após fazer essas mudanças, o próximo release com uma tag (ex: `v2.0.1`) deve funcionar corretamente.
