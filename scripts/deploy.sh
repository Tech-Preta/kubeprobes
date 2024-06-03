# Build the Docker image
docker build -t nataliagranato/kubeprobes:0.1.0 .

# Push the Docker image to a registry
docker push nataliagranato/kubeprobes:0.1.0

# Deploy the application to Kubernetes
kubectl apply -f deployment.yaml