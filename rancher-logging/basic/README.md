This directory contains what you need to get a basic logging setup
up and running. It's for debugging and sanity checks, not for production.

1. Start a fresh Rancher instance on k3s.
2. Install Rancher Logging through UI with defaults
3. kubectl create namespace quickstart
4. kubectl apply -f http-echo.yaml
5. helm upgrade --install --wait --namespace quickstart log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
6. kubectl apply -f manifests.yaml

You can then follow the logs on the http-echo pod. Note that it will
take a bit of time to come up.
