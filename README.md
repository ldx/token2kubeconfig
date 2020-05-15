# token2kubeconfig

Create a kubeconfig file from a Kubernetes service account token.

## Build

    $ go build

## Usage

This is meant to be used in a pod:

    $ ./token2kubeconfig > kubeconfig

It will generate a kubeconfig file (usable by e.g. kubectl) from the pod service account token.
