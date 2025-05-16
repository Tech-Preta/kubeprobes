# Docker

Este documento descreve como usar o Docker para executar o kubeprobes.

## Construir a imagem Docker

Para construir a imagem Docker, execute o seguinte comando na raiz do projeto:

```sh
docker build -t kubeprobes:1.1.0 .
```

# Como executar o kubeprobes via Docker

## Exibir ajuda

```sh
docker run --rm kubeprobes:1.1.0 --help
```

## Rodar o scan em um cluster Kubernetes

> **Atenção:** Se estiver usando um cluster local (ex: kind, minikube), adicione a flag `--network=host` para o container conseguir acessar o cluster.

Monte seu kubeconfig local como volume e passe o caminho para o comando:

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.1.0 scan --kubeconfig /kubeconfig
```

### Analisar um namespace específico

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.1.0 scan --kubeconfig /kubeconfig --namespace NOME_DO_NAMESPACE
```

### Analisar todos os namespaces

```sh
docker run --rm --network=host \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.1.0 scan --kubeconfig /kubeconfig --namespace ""
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

Você pode adicionar outras opções, como namespace ou tipo de probe:

```sh
docker run --rm \
  -v $HOME/.kube/config:/kubeconfig:ro \
  kubeprobes:1.1.0 scan --kubeconfig /kubeconfig --namespace default --probe-type liveness
```

Consulte as opções disponíveis com:

```sh
docker run --rm kubeprobes:1.1.0 scan --help
```