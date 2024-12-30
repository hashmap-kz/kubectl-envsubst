# kubectl-envsubst - a plugin for kubectl, used for expand env-vars in manifests

### Usage:

```
kubectl envsubst apply -f manifests/ --allowed-prefixes=CI_,APP_
```

### Features:

- --allowed-vars     : cmd flag, that consumes a list of names that allowed for expansion
- --allowed-prefixes : cmd flag, that consumes a list of prefixes (APP_), and variables that not match will be ignored 
