# kubectl-login-pod
The kubectl-login-pod is a kubectl plugin that allows you to access pods across amespaces in your Kubernetes cluster. This makes it easy to debug and manage your applications.

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
