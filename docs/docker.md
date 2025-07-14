# Docker

Este documento descreve como usar o Docker para executar o kubeprobes.

## Construir a imagem Docker

Para construir a imagem Docker, execute o seguinte comando na raiz do projeto:

```sh
docker build -t kubeprobes:1.2.0 .
```

# Como executar o kubeprobes via Docker

## Exibir ajuda

```sh
docker run --rm kubeprobes:1.2.0 --help
docker run --rm kubeprobes:1.2.0 scan --help
```

## Exemplos Básicos

> **Atenção:** Se estiver usando um cluster local (ex: kind, minikube), adicione a flag `--network=host` para o container conseguir acessar o cluster.

### Rodar o scan em um cluster Kubernetes

Monte seu kubeconfig local como volume e passe o caminho para o comando:

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig
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
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --all-namespaces
```

## Exemplos com Novas Funcionalidades

### Diferentes formatos de saída

```sh
# Saída em JSON
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --output json

# Saída em YAML
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --output yaml
```

### Verificar tipos específicos de probe

```sh
# Verificar apenas liveness probes
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --probe-type liveness

# Verificar readiness probes com recomendações
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --probe-type readiness --recommendation
```

### Opções avançadas

```sh
# Falhar em avisos (exit code 1 para warnings)
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan --kubeconfig /kubeconfig --fail-on-warn

# Exemplo completo com todas as opções
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.2.0 scan \
    --kubeconfig /kubeconfig \
    --all-namespaces \
    --probe-type liveness \
    --recommendation \
    --output json \
    --fail-on-warn
```

## Exemplos de saída

### Saída texto (padrão)
```
[WARNING] Pod default/vote-595d458c7c-fmj9g (container: vote) is missing liveness probe
  Recommendation: Add a liveness probe to ensure the container is running correctly.
[WARNING] Pod default/vote-595d458c7c-fmj9g (container: vote) is missing readiness probe
  Recommendation: Add a readiness probe to ensure the container is ready to accept traffic.
```

### Saída JSON
```json
{
  "issues": [
    {
      "namespace": "default",
      "podName": "vote-595d458c7c-fmj9g",
      "containerName": "vote",
      "probeType": "liveness",
      "message": "missing liveness probe",
      "recommendation": "Add a liveness probe to ensure the container is running correctly."
    }
  ],
  "summary": "Found 1 probe issues in default",
  "namespace": "default",
  "exitCode": 0
}
```

Se não houver pods no namespace:
```
No pods found in namespace NOME_DO_NAMESPACE
```

Consulte as opções disponíveis com:

```sh
docker run --rm kubeprobes:1.2.0 scan --help
```