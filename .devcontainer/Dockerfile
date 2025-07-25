FROM mcr.microsoft.com/vscode/devcontainers/base:1.0-ubuntu-22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    software-properties-common \
    git \
    unzip \
    jq \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# Docker CLI
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list \
    && apt-get update && apt-get install -y docker-ce-cli docker-buildx-plugin docker-compose-plugin \
    && rm -rf /var/lib/apt/lists/*

# kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/local/bin/

# KinD
RUN curl -Lo /usr/local/bin/kind https://kind.sigs.k8s.io/dl/v0.21.0/kind-linux-amd64 \
    && chmod +x /usr/local/bin/kind

# Helm
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Go programming language
RUN GO_VERSION="1.24.5" \
    && curl -Lo go.tar.gz "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
    && sudo rm -rf /usr/local/go \
    && sudo tar -C /usr/local -xzf go.tar.gz \
    && rm go.tar.gz

# Set Go environment variables
ENV GOROOT=/usr/local/go
ENV GOPATH=/home/vscode/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

# Create Go workspace directories
RUN mkdir -p $GOPATH/{bin,src,pkg} \
    && chown -R vscode:vscode $GOPATH

# Install common Go tools
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# K9s
RUN curl -Lo k9s.tar.gz https://github.com/derailed/k9s/releases/latest/download/k9s_Linux_amd64.tar.gz \
    && tar -xf k9s.tar.gz \
    && mv k9s /usr/local/bin/ \
    && rm k9s.tar.gz LICENSE README.md

# Permissões e sudo sem senha para vscode
RUN groupadd -f docker \
    && usermod -aG docker vscode \
    && mkdir -p /home/vscode/.kube /home/vscode/.kind \
    && chown -R vscode:vscode /home/vscode/.kube /home/vscode/.kind \
    && echo "vscode ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/99-vscode

USER vscode