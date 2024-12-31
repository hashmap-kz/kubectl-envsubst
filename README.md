# kubectl-envsubst - a plugin for kubectl, used for expand env-vars in manifests

### Features:

- Expand environment-variables in manifests, passed to kubectl before applying them
- Uses prefixes, and allowed variables (you're able to be totally sure what it's going to be substituted)
- Has --strict mode (if some variables are not expanded, it fails)
- Uses all known kubectl args (by just proxying them as is, without any additional actions)
- ZERO dependencies, and I mean literally zero

### Usage:

```
# substitute variables whose names start with one of the prefixes
kubectl envsubst apply -f manifests/ --envsubst-allowed-prefixes=CI_,APP_

# substitute well-defined variables
kubectl envsubst apply -f manifests/ --envsubst-allowed-vars=CI_PROJECT_NAME,CI_COMMIT_REF_NAME,APP_IMAGE

# mixed mode, check both full match and prefix match 
kubectl envsubst apply -f manifests/ --envsubst-allowed-prefixes=CI_,APP_ --envsubst-allowed-vars=HOME,USER
```

### Options:

- --envsubst-allowed-vars     : cmd flag, that consumes a list of names that allowed for expansion
- --envsubst-allowed-prefixes : cmd flag, that consumes a list of prefixes (APP_), and variables that not match will be ignored 


