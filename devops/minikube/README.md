# Minikube

For Kubernetes development, we're using [Minikube](https://minikube.sigs.k8s.io/docs/start/).

Feel free to use your preference, but documentation will be based on Minikube.

## Installing
Follow instructions on [the web site](https://minikube.sigs.k8s.io/docs/start/).

## Environment Variables

Set

`MINIKUBE_ROOTLESS=false`

in your shell.

At this time, we're using traditional (root) docker and k8s.  In the future, a move to Rootless would be wonderful

TODO: Move to Rootless