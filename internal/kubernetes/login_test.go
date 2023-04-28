package kubecli

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

type MockPodGetter struct{}

func (mpg *MockPodGetter) GetPods(clientset kubernetes.Interface) ([]corev1.Pod, error) {
	podList := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-pod-1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-pod-2",
			},
		},
	}
	return podList, nil
}

type MockPodExecutor struct{}

func (mpe *MockPodExecutor) ExecInPod(clientset kubernetes.Interface, config *rest.Config, podName string, namespace string) error {
	return nil
}

func TestGetPods(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	podGetter := &MockPodGetter{}

	pods, err := podGetter.GetPods(clientset)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(pods) != 2 {
		t.Errorf("Expected 2 pods, got: %d", len(pods))
	}

	expectedNames := []string{"test-pod-1", "test-pod-2"}
	for i, pod := range pods {
		if pod.Name != expectedNames[i] {
			t.Errorf("Expected pod name %s, got: %s", expectedNames[i], pod.Name)
		}
	}
}

func TestExecInPod(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	config := &rest.Config{}
	podExecutor := &MockPodExecutor{}

	err := podExecutor.ExecInPod(clientset, config, "test-pod", "test-namespace")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
