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

	mykube "github.com/jedipunkz/kubectl-login-pod/internal/kubernetes"
)

var namespace string

var RootCmd = &cobra.Command{
	Use:   "kubectl login pod",
	Short: "kubectl plugin to login to a kubernetes pod",
	Run: func(cmd *cobra.Command, args []string) {

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

		podIndex, err := fuzzyfinder.Find(
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

		selectedPod := podList[podIndex]

		containerGetter := &mykube.ContainerGetterImpl{}
		containerList, err := containerGetter.GetContainers(clientset, selectedPod.Name, selectedPod.Namespace)

		if err != nil {
			fmt.Printf("Error getting containers: %v\n", err)
			os.Exit(1)
		}

		containerIndex, err := fuzzyfinder.Find(
			containerList,
			func(i int) string {
				return containerList[i].Name
			},
			fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
				if i == -1 {
					return ""
				}
				return fmt.Sprintf(
					"%s: %s\n%s: %s\n",
					color.CyanString("Name"), containerList[i].Name,
					color.MagentaString("Image"), containerList[i].Image,
				)
			}))

		if err != nil {
			fmt.Printf("Error finding container: %v\n", err)
			os.Exit(1)
		}

		selectedContainer := containerList[containerIndex]

		fmt.Printf("Logging into container %s in pod %s in namespace %s...\n", selectedContainer.Name, selectedPod.Name, selectedPod.Namespace)

		podExecutor := &mykube.PodExecutorImpl{}
		if err := podExecutor.ExecInPod(clientset, config, selectedPod.Name, selectedPod.Namespace, selectedContainer.Name); err != nil {
			fmt.Printf("Error executing command in container: %v\n", err)
		}
	},
}

func init() {
	RootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace")
}
