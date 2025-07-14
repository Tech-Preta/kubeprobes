# Docker

Este documento descreve como usar o Docker para executar o kubeprobes.

## Construir a imagem Docker

Para construir a imagem Docker, execute o seguinte comando na raiz do projeto:

```sh
docker build -t kubeprobes:1.2.0 .
```

# Como executar o kubeprobes via Docker

## Comandos Básicos

### Exibir ajuda geral
```sh
docker run --rm kubeprobes:1.2.0 --help
```

### Exibir informações de versão
```sh
docker run --rm kubeprobes:1.2.0 version
```

### Exibir ajuda do comando scan
```sh
docker run --rm kubeprobes:1.2.0 scan --help
```

## Comando Scan - Análise de Probes

> **Atenção:** Se estiver usando um cluster local (ex: kind, minikube), adicione a flag `--network=host` para o container conseguir acessar o cluster.

### Scan básico

Monte seu kubeconfig local como volume e execute o scan:

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig
```

### Usando variável de ambiente KUBECONFIG

Alternativamente, você pode usar a variável de ambiente KUBECONFIG:

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  -e KUBECONFIG=/kubeconfig \
  kubeprobes:1.2.0 scan
```

### Analisar um namespace específico

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --namespace NOME_DO_NAMESPACE
```

### Analisar todos os namespaces

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --namespace ""
```

### Scan com recomendações

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --recommendation
```

### Scan de um tipo específico de probe

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --probe-type liveness
```

## Comando Version

### Versão detalhada
```sh
docker run --rm kubeprobes:1.2.0 version
```

### Versão resumida (para scripts)
```sh
docker run --rm kubeprobes:1.2.0 version --output=short
```

### Versão em formato JSON
```sh
docker run --rm kubeprobes:1.2.0 version --output=json
```

### Exemplos de saída

```
[WARNING] Pod default/vote-595d458c7c-fmj9g (container: vote) is missing a liveness probe
[WARNING] Pod default/vote-595d458c7c-fmj9g (container: vote) is missing a readiness probe  
[WARNING] Pod default/vote-595d458c7c-fmj9g (container: vote) is missing a startup probe
```

Se não houver pods no namespace:
```
No pods found in namespace NOME_DO_NAMESPACE
```

## Comando Completion

Gerar scripts de autocompletion:

```sh
# Para bash
docker run --rm kubeprobes:1.2.0 completion bash

# Para zsh
docker run --rm kubeprobes:1.2.0 completion zsh

# Para fish
docker run --rm kubeprobes:1.2.0 completion fish

# Para PowerShell
docker run --rm kubeprobes:1.2.0 completion powershell
```

## Dicas Avançadas

Você pode combinar várias opções para análises mais específicas:

```sh
# Exemplo completo: scan de probes de liveness com recomendações em namespace específico
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan \
    --kubeconfig /kubeconfig \
    --namespace production \
    --probe-type liveness \
    --recommendation
```

Para ver todas as opções disponíveis de qualquer subcomando:

```sh
docker run --rm kubeprobes:1.2.0 <subcomando> --help
```