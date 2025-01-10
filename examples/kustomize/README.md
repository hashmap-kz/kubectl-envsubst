## A brief example of usage along with 'kubectl kustomize'

### Prerequisites:

You can use a Kind cluster for testing. If you haven’t installed it yet, check out the [Kind installation guide](https://kind.sigs.k8s.io/).

### Deploy:

Execute scripts one by one:

```
# prepare kind cluster
bash 00-setup-kind.sh

# kustomize, envsubst, apply manifests
bash 01-kustomize-envsubst-apply.sh
```

### Verification:

Check the result: http://localhost:31355

_If your Kind cluster is running on a remote machine, replace ‘localhost’ with the machine’s IP address._
