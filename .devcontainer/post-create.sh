#!/bin/bash

set -e
set -x

# Função para verificar se o Docker está funcionando
check_docker() {
  echo "Verificando conectividade com o Docker..."
  if ! sudo docker info > /dev/null 2>&1; then
    echo "ERRO: Não foi possível conectar ao Docker. Verifique se o socket está mapeado corretamente."
    return 1
  fi
  echo "Conectividade com Docker OK!"
  return 0
}

# Função para criar um cluster KinD
create_kind_cluster() {
  echo "Criando cluster KinD..."

  # Verificar se já existe um cluster com este nome
  if sudo kind get clusters 2>/dev/null | grep -q "devcontainer-cluster"; then
    echo "Cluster 'devcontainer-cluster' já existe."
    return 0
  fi

  # Configuração do cluster KinD com mapeamento de portas
  cat > /tmp/kind-config.yaml << EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 8081
    hostPort: 8081
    protocol: TCP
  - containerPort: 443
    hostPort: 8443
    protocol: TCP
  - containerPort: 6443
    hostPort: 16443    # novo hostPort para evitar conflito
    protocol: TCP
EOF

  sudo kind create cluster --name devcontainer-cluster --config /tmp/kind-config.yaml

  echo "Configurando o contexto do kubectl..."
  mkdir -p $HOME/.kube
  sudo kind get kubeconfig --name devcontainer-cluster > $HOME/.kube/config

  # Corrige o endpoint do servidor para localhost:16443 e nomes de contexto/cluster/user
  sed -i 's|server: https://127.0.0.1:[0-9]*|server: https://127.0.0.1:16443|' $HOME/.kube/config
  sed -i 's|name: kind-.*|name: kind-devcontainer-cluster|g' $HOME/.kube/config
  sed -i 's|cluster: kind-.*|cluster: kind-devcontainer-cluster|g' $HOME/.kube/config
  sed -i 's|user: kind-.*|user: kind-devcontainer-cluster|g' $HOME/.kube/config
  sed -i 's|current-context: kind-.*|current-context: kind-devcontainer-cluster|g' $HOME/.kube/config

  sudo chown vscode:vscode $HOME/.kube/config

  return 0
}

# Instalar dependências do projeto
install_project_deps() {
  echo "Instalando dependências do projeto Go..."
  
  # Verificar se Go está instalado
  if ! command -v go &> /dev/null; then
    echo "Go não encontrado. Instalando..."
    return 1
  fi
  
  # Navegar para o diretório do projeto
  cd /workspaces/kubeprobes
  
  # Instalar dependências do Go
  if [ -f "go.mod" ]; then
    echo "Instalando dependências do Go..."
    go mod download
    go mod tidy
    
    echo "Executando go vet..."
    go vet ./...
    
    echo "Compilando o projeto..."
    make build || {
      echo "Falha no build com Makefile, tentando build direto..."
      go build -o bin/kubeprobes ./cmd/kubeprobes
    }
    
    echo "Executando testes..."
    go test -v ./... || { echo "Tests failed"; exit 1; }
    
    echo "Projeto Go configurado com sucesso!"
  else
    echo "go.mod não encontrado no diretório do projeto"
  fi
}

# Execução principal
main() {
  echo "Iniciando configuração do ambiente de desenvolvimento..."

  if check_docker; then
    create_kind_cluster
    install_project_deps
    echo "Ambiente de desenvolvimento configurado com sucesso!"
  else
    echo "Falha ao configurar o ambiente de desenvolvimento."
    exit 1
  fi
}

# Executar o script principal
main