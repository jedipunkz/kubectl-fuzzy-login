package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	mykube "github.com/jedipunkz/kubecli/internal/kubernetes"
)

// Namespace variable for command flags
var namespace string

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

	podList, err := podGetter.GetPods(clientset, namespace)
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
			return fmt.Sprintf(
				"%s: %s\n%s: %s\n%s: %s\n%s: %s\n",
				color.CyanString("Name"), podList[i].Name,
				color.MagentaString("Namespace"), podList[i].Namespace,
				color.YellowString("Creating Timestamp"), podList[i].CreationTimestamp,
				color.RedString("UUID"), podList[i].UID,
			)
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
	loginCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace")
	rootCmd.AddCommand(loginCmd)
}
