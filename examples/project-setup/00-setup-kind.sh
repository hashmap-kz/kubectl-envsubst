#!/bin/bash
set -euo pipefail

# prepare config for the 'kind' cluster
cat <<EOF >kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: "kubectl-envsubst"
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 32501
        hostPort: 32501
        protocol: TCP
      - containerPort: 32502
        hostPort: 32502
        protocol: TCP
      - containerPort: 32503
        hostPort: 32503
        protocol: TCP
EOF

# setup cluster with kind, to safely test in a sandbox
if kind get clusters | grep "kubectl-envsubst"; then
  kind delete clusters "kubectl-envsubst"
fi
kind create cluster --config=kind-config.yaml
kubectl config set-context "kind-kubectl-envsubst"
rm -f kind-config.yaml
