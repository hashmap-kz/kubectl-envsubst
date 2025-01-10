#!/bin/bash
set -euo pipefail

# prepare config for the 'kind' cluster
cat <<EOF >kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: "kubectl-envsubst-kustomize"
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 31355
        hostPort: 31355
        protocol: TCP
EOF

# setup cluster with kind, to safely test in a sandbox
if kind get clusters | grep "kubectl-envsubst-kustomize"; then
  kind delete clusters "kubectl-envsubst-kustomize"
fi
kind create cluster --config=kind-config.yaml
kubectl config set-context "kind-kubectl-envsubst-kustomize"
rm -f kind-config.yaml
