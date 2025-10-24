This directory contains a basic demonstration of collecting logs
via opentelemetry-operator.

1. helm install cert-manager jetstack/cert-manager --namespace cert-manager --version v1.19.1 --set crds.enabled=true --wait --create-namespace
2. kubectl apply -f otel-operator.yaml
3. kubectl apply -f manifest.yaml

Then:

1. kubectl get otelcol,daemonset,pod
2. kubectl logs <pod_name> -f
