# kubectl-login-pod
The kubectl-login-pod is a kubectl plugin that allows you to access pods across amespaces in your Kubernetes cluster. This makes it easy to debug and manage your applications.

## Installation
To install kubectl-login-pod, follow these steps:

1. Build the kubectl-login-pod binary using the Go compiler:
```bash
go build
```

2.Copy the generated binary to a directory in your PATH:

```bash
cp kubectl-login-pod /your/bin/path
```

Replace /your/bin/path with the directory in your PATH where you want to copy the binary.

## Usage

There are two main ways to use kubectl-login-pod:

### Access pods across all namespaces:

```bash
kubectl login pod
```

This command provides access to all pods within your Kubernetes cluster, irrespective of the namespace they are in. This is useful when you need a broad overview of your cluster's pods.

### Access pods within a specific namespace:

```bash
kubectl login pod -n <namespace>
```

This command limits access to the pods within the specified namespace. This is helpful when you want to focus on a specific subset of your cluster's pods.

Replace <namespace> with the name of the namespace you want to target.

## Screenshot

<img src="https://raw.githubusercontent.com/jedipunkz/kubecli/main/static/kubectl-login-pod.gif">

## LICENSE

Apache LICENSE

## Author

- jedipunkz
