# Projeto

Este projeto verifica se todos os pods em um cluster Kubernetes têm sondas de vida, prontidão e inicialização configuradas.

## Como executar

Defina a variável de ambiente KUBECONFIG para o caminho do seu arquivo kubeconfig, então execute:

```bash
go run cmd/projeto/main.go
