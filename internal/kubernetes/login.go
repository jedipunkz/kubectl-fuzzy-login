package kubecli

import (
	"context"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	defaultShell = "/bin/sh"
)

type PodGetter interface {
	GetPods(clientset kubernetes.Interface) ([]corev1.Pod, error)
}

type PodExecutor interface {
	ExecInPod(clientset *kubernetes.Clientset, config *rest.Config, podName string, namespace string) error
}

type PodGetterImpl struct{}

func (p *PodGetterImpl) GetPods(clientset kubernetes.Interface) ([]corev1.Pod, error) {
	podList, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

type PodExecutorImpl struct{}

func (p *PodExecutorImpl) ExecInPod(clientset kubernetes.Interface, config *rest.Config, podName string, namespace string) error {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true").
		Param("command", defaultShell)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return err
	}

	return nil
}
