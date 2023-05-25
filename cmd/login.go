package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	mykube "github.com/jedipunkz/kubecli/internal/kubernetes"
)

// func main() {
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a Kubernetes Pod",
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
}

func login() {
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating clientset: %v\n", err)
		os.Exit(1)
	}

	podGetter := &mykube.PodGetterImpl{}

	podList, err := podGetter.GetPods(clientset)
	if err != nil {
		fmt.Printf("Error getting pods: %v\n", err)
		os.Exit(1)
	}

	index, err := fuzzyfinder.FindMulti(
		podList,
		func(i int) string {
			return podList[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Name: %s\nNamespace: %s\n", podList[i].Name, podList[i].Namespace)
		}))

	if err != nil {
		fmt.Printf("Error finding pod: %v\n", err)
		os.Exit(1)
	}

	selectedPod := podList[index[0]]
	fmt.Printf("Logging into pod %s in namespace %s...\n", selectedPod.Name, selectedPod.Namespace)
	podExecutor := &mykube.PodExecutorImpl{}
	if err := podExecutor.ExecInPod(clientset, config, selectedPod.Name, selectedPod.Namespace); err != nil {
		fmt.Printf("Error executing command in pod: %v\n", err)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
