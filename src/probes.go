package main

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		log.Fatalf("Error: KUBECONFIG environment variable not set")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %s", err.Error())
	}

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if container.LivenessProbe == nil {
				log.Printf("Pod %s in namespace %s does not have a liveness probe\n", pod.Name, pod.Namespace)
				log.Println("Recommendation: Add a liveness probe to ensure the container is running correctly.")
			}
			if container.ReadinessProbe == nil {
				log.Printf("Pod %s in namespace %s does not have a readiness probe\n", pod.Name, pod.Namespace)
				log.Println("Recommendation: Add a readiness probe to ensure the container is ready to accept traffic.")
			}
			if container.StartupProbe == nil {
				log.Printf("Pod %s in namespace %s does not have a startup probe\n", pod.Name, pod.Namespace)
				log.Println("Recommendation: Add a startup probe to ensure the container has started successfully.")
			}
		}
	}
}
