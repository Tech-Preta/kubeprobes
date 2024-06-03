# KubeProbes - CLI Tool for Scanning Kubernetes Probes

Probes é uma ferramenta de linha de comando (CLI) desenvolvida em Go para escanear workloads do Kubernetes em busca de probes (sondas) definidas.

## Requisitos

- Go 1.13 ou superior
- kubectl
- Um cluster Kubernetes acessível

## Instalação

1. Clone o repositório:

```bash
https://github.com/Tech-Preta/kubeprobes.git
```

2. Entre no diretório do projeto:

```bash
cd kubeprobes
```

3. Compile o código fonte:

```bash
cd src
go build -o kubeprobes
```

4. Mova o binário para o diretório /usr/local/bin:

```bash
sudo mv kubeprobes /usr/local/bin
```

5. Verifique se a instalação foi bem sucedida:

```bash
kubeprobes --help
```

## Uso

### Comandos Disponíveis

- `scan`: Escaneia workloads do Kubernetes em busca de probes.
  
  Exemplo de uso:

```bash

kubeprobes scan -k <caminho-para-o-kubeconfig> -c <contexto-kubeconfig> -n <namespace> -p <tipo-de-probe> -r
```

### Flags

- `-k, --kubeconfig`: Caminho para o arquivo kubeconfig.
- `-c, --kubeContext`: Contexto do Kubernetes.
- `-n, --namespace`: Namespace do Kubernetes.
- `-p, --probe-type`: Tipo de probe para escanear (liveness, readiness, startup).
- `-r, --recommendation`: Mostrar recomendações para sondas ausentes.

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir um issue ou enviar um pull request.
