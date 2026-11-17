# Guia de Contribuição

Obrigado por considerar contribuir com o projeto! Este documento fornece diretrizes e instruções para contribuir.

## Como Contribuir

1. Faça um fork do repositório
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Faça commit das suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Faça push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request

## Padrões de Código

- Siga as [boas práticas de Go](https://golang.org/doc/effective_go)
- Use `gofmt` para formatar seu código
- Execute `go vet` para análise estática
- Execute os linters com `make lint`
- **Testes são obrigatórios** para novas funcionalidades e correções de bugs
- Mantenha a cobertura de testes acima de 75%

## Requisitos de Testes

### Testes Obrigatórios

Todas as contribuições devem incluir testes adequados:

1. **Novas funcionalidades**: Devem ter testes unitários cobrindo cenários de sucesso, erro e casos extremos
2. **Correções de bugs**: Devem incluir um teste que reproduza o bug e verifique a correção
3. **Refatorações**: Devem manter ou melhorar a cobertura de testes existente

### Executando Testes

```bash
# Executar todos os testes
make test

# Executar testes com cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Executar testes de um pacote específico
go test -v ./internal/scanner/
```

### Padrões de Testes

- Use **testes baseados em tabela** para múltiplos cenários
- **Nomeie testes** de forma descritiva (ex: `TestScanner_InvalidProbeType`)
- **Mock dependências externas** (ex: API do Kubernetes)
- **Teste cenários de erro** além de casos de sucesso
- **Mantenha testes rápidos** (< 1 segundo por teste)

## Processo de Pull Request

1. **Execute todos os testes** e certifique-se de que passam (`make test`)
2. **Verifique a cobertura de testes** (`go test -coverprofile=coverage.out ./...`)
3. Atualize a documentação se necessário
4. **Adicione testes** para novas funcionalidades ou correções de bugs
5. Execute o linter (`make lint`) e corrija quaisquer problemas
6. Atualize o CHANGELOG.md
7. Descreva suas mudanças no PR de forma clara e detalhada

### Checklist do Pull Request

- [ ] Testes adicionados/atualizados e passando
- [ ] Cobertura de testes mantida/melhorada
- [ ] Código formatado (`gofmt`)
- [ ] Linter sem erros (`make lint`)
- [ ] Documentação atualizada (se aplicável)
- [ ] CHANGELOG.md atualizado

## Relatando Bugs

Use o template de issue para bugs e inclua:
- Descrição clara do problema
- Passos para reproduzir
- Comportamento esperado
- Comportamento atual
- Screenshots (se aplicável)
- Ambiente (sistema operacional, versão do Go, etc.)

## Sugerindo Features

Use o template de issue para features e inclua:
- Descrição clara da feature
- Casos de uso
- Benefícios
- Possíveis implementações

## Código de Conduta

Por favor, leia e siga nosso [Código de Conduta](CODE_OF_CONDUCT.md).

## Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a mesma licença do projeto.