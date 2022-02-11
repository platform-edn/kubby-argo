package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/platform-edn/kubby"
)

func main() {
	baseDir := getDir()
	argoNamespace := "argo"
	clusterName := "argo-cluster"
	kubeconfigFilePath := filepath.Join(baseDir, "local", "kubeconfig.yaml")

	cluster, err := kubby.NewKubeCluster(
		kubby.WithName(clusterName),
		kubby.WithControlNodes(1),
		kubby.WithWorkerNodes(2),
		kubby.WithKubeConfigPath(kubeconfigFilePath),
		kubby.ShouldStartOnCreation(true),
		kubby.WithMaxAttempts(10),
		kubby.WithNodePorts(&kubby.NodePort{
			Container: "32746",
			Host:      "32746",
		}),
		kubby.WithNamespaces(argoNamespace),
		kubby.WithHelmCharts(
			&kubby.HelmChart{
				Name:      "argo-events",
				Namespace: argoNamespace,
				Path:      filepath.Join(baseDir, "charts", "argo-events"),
			},
			&kubby.HelmChart{
				Name:      "argo-workflows",
				Namespace: argoNamespace,
				Path:      filepath.Join(baseDir, "charts", "argo-workflows"),
			},
			&kubby.HelmChart{
				Name:      "argo-event-bus",
				Namespace: argoNamespace,
				Path:      filepath.Join(baseDir, "charts", "event-bus"),
			},
		),
	)
	if err != nil {
		log.Fatalf("error creating cluster: %s", err)
	}

	fmt.Println("cluster created! Images can now be pushed to localhost:5000 ...")
	fmt.Println("Argo Workflows can be accessed at http://localhost:32746/workflows ...")
	fmt.Printf("set KUBECONFIG to %s to connect with kubectl ...\n", kubeconfigFilePath)

	config := os.Getenv("KUBECONFIG")
	os.Setenv("KUBECONFIG", kubeconfigFilePath)

	fmt.Println("press any key to destroy the cluster...")
	fmt.Scanln()

	fmt.Println("Destroying cluster ...")
	err = cluster.Delete()
	if err != nil {
		log.Fatalf("error deleting cluster: %s", err)
	}

}

func getDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return dir
}
