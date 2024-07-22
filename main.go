package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"time"
)

var KubeConfigPath string
var config *rest.Config
var clientset *kubernetes.Clientset

var configNotFound = fmt.Errorf("🔥 no kubeconfig found at %s. ➡️ Configure your cluster and restart program", KubeConfigPath)

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
		KubeConfigPath = filepath.Join(home, ".kube", "config")
	} else {
		KubeConfigPath = ""
	}
}

func main() {
	if !doesKubeConfigExist(KubeConfigPath) {
		fmt.Println(configNotFound)
		os.Exit(1)
	}

	// Build the KubeConfigPath
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", KubeConfigPath)
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

	// add header
	printFormat := "%-14v %-41v %-8v %-10v %-11v %-3vh\n"
	fmt.Printf(printFormat, "NAMESPACE", "NAME", "READY", "STATUS", "RESTARTS", "AGE")
	// add row per pod
	for _, pod := range pods.Items {
		//prettyPrint("debugfile", pod)
		containerRestarts := int32(0)
		lengthOfContainerStatuses := len(pod.Status.ContainerStatuses)
		ageOfContainer := "unknown"
		readyStatus := "0/1"
		if lengthOfContainerStatuses > 0 {
			containerRestarts = pod.Status.ContainerStatuses[lengthOfContainerStatuses-1].RestartCount
			lastContainerStartTime := pod.Status.ContainerStatuses[lengthOfContainerStatuses-1].State.Running.StartedAt.Time
			ageOfContainer = fmt.Sprintf("%v", int(time.Now().Sub(lastContainerStartTime).Hours()))
			if pod.Status.ContainerStatuses[lengthOfContainerStatuses-1].Ready {
				readyStatus = "1/1"
			}
		}
		fmt.Printf(printFormat, pod.Namespace, pod.Name, readyStatus, pod.Status.Phase, containerRestarts, ageOfContainer)
	}
}
