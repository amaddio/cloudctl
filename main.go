package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeConfigPath string
var config *rest.Config
var clientset *kubernetes.Clientset

func doesKubeConfigExist(path string) bool {
	file, err := os.Open(path)
	defer file.Close()
	if err == nil {
		return true
	} else {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("configuration does not exist at:", path)
		} else {
			fmt.Println("file exists but can not be opened:", err)
		}
		return false
	}
}

func init() {
	if home := homedir.HomeDir(); home != "" {
		kubeConfigPath = filepath.Join(home, ".kube", "config")
	} else {
		kubeConfigPath = ""
	}
	configNotFound := fmt.Errorf("🔥 no kubeconfig found at %s. ➡️ Configure your cluster and restart program", kubeConfigPath)

	if !doesKubeConfigExist(kubeConfigPath) {
		fmt.Println(configNotFound)
		os.Exit(1)
	}

	// Build the kubeConfigPath
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "cloudctl",
		Short: "cloudctl is a CLI for interacting with Kubernetes clusters",
	}

	// Add subcommands
	rootCmd.AddCommand(getCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var getCmd = &cobra.Command{
	Use:   "get [resource]",
	Short: "Get resources from the Kubernetes cluster",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		switch resource {
		case "pods":
			getPods()
		default:
			fmt.Printf("Unknown resource: %s\n", resource)
		}
	},
}

func getPods() {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error getting pods: %v\n", err)
		return
	}

	for _, pod := range pods.Items {
		fmt.Printf("Namespace: %s, Name: %s\n", pod.Namespace, pod.Name)
	}
}
