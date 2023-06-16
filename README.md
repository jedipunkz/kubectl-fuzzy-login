# kubectl-fuzzy-login
The kubectl-fuzzy-login is a kubectl plugin that allows you to fuzzy find pods/containers and login to container across namespaces in your kubernetes cluster.

## Screenshot

<img src="https://raw.githubusercontent.com/jedipunkz/kubecli/main/static/kubectl-fuzzy-login.gif">

## Installation
To install kubectl-fuzzy-login, follow these steps:

### Install with Krew

- require: [Install krew](https://krew.sigs.k8s.io/docs/user-guide/setup/install/)

```shell
git clone https://github.com/jedipunkz/kubectl-fuzzy-login.git
kubectl krew install --manifest=./kubectl-fuzzy-login/krew/fuzzy-login.yaml
```

### Manual Install

#### 1. Build the kubectl-fuzzy-login binary using the Go compiler:
```bash
go build
```

#### 2. Copy the generated binary to a directory in your PATH:

```bash
cp kubectl-fuzzy-login /your/bin/path
```

Replace /your/bin/path with the directory in your PATH where you want to copy the binary.

## Usage

There are two main ways to use kubectl-fuzzy-login:

### Login to a pod across all namespaces:

```bash
kubectl fuzzy login
```

This command provides access to all pods within your Kubernetes cluster, irrespective of the namespace they are in. This is useful when you need a broad overview of your cluster's pods.

### Login to a pod within a specific namespace:

```bash
kubectl fuzzy login -n <namespace>
```

This command limits access to the pods within the specified namespace. This is helpful when you want to focus on a specific subset of your cluster's pods.

Replace <namespace> with the name of the namespace you want to target.

  ## LICENSE

Apache LICENSE

## Author

- jedipunkz
