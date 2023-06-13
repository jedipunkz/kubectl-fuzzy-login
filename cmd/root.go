package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	mykube "github.com/jedipunkz/kubectl-fuzzy-login/internal/kubernetes"
)

var namespace string
var shell string

var RootCmd = &cobra.Command{
	Use:   "kubectl login pod",
	Short: "kubectl plugin to login to a kubernetes pod",
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return fmt.Errorf("error building kubeconfig: %v", err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("error creating clientset: %v", err)
		}

		selectedPod, err := selectPod(clientset)
		if err != nil {
			return err
		}

		selectedContainer, err := selectContainer(clientset, selectedPod)
		if err != nil {
			return err
		}

		fmt.Printf("Logging into container %s in pod %s in namespace %s...\n", selectedContainer.Name, selectedPod.Name, selectedPod.Namespace)

		podExecutor := &mykube.PodExecutorImpl{}
		if err := podExecutor.ExecInPod(clientset, config, selectedPod.Name, selectedPod.Namespace, selectedContainer.Name, shell); err != nil {
			return fmt.Errorf("error executing command in container: %v", err)
		}

		return nil
	},
}

func selectPod(clientset *kubernetes.Clientset) (*v1.Pod, error) {
	podGetter := &mykube.PodGetterImpl{}

	podList, err := podGetter.GetPods(clientset, namespace)
	if err != nil {
		return nil, fmt.Errorf("error getting pods: %v", err)
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
		return nil, fmt.Errorf("error finding pod: %v", err)
	}

	return &podList[podIndex], nil
}

func selectContainer(clientset *kubernetes.Clientset, selectedPod *v1.Pod) (*v1.Container, error) {
	containerGetter := &mykube.ContainerGetterImpl{}
	containerList, err := containerGetter.GetContainers(clientset, selectedPod.Name, selectedPod.Namespace)

	if err != nil {
		return nil, fmt.Errorf("error getting containers: %v", err)
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
		return nil, fmt.Errorf("error finding container: %v", err)
	}

	return &containerList[containerIndex], nil
}

func init() {
	RootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace")
	RootCmd.Flags().StringVarP(&shell, "shell", "s", "/bin/sh", "Shell to use")
}
