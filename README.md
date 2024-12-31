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

# example:
export APP_NAME=nginx
export APP_IMAGE_NAME=nginx
export APP_IMAGE_TAG=latest
kubectl envsubst apply -f testdata/subst/01.yaml --dry-run=client -oyaml --envsubst-allowed-prefixes=APP_
```

### Options:

- --envsubst-allowed-vars     : cmd flag, that consumes a list of names that allowed for expansion
- --envsubst-allowed-prefixes : cmd flag, that consumes a list of prefixes (APP_), and variables that not match will be ignored 

### Implementation details

Substitution of environmental variables without checking for inclusion in one of the filter list is not used intentionally, because this behavior can lead to subtle mistakes.

If a variable is not found in the filter list and strict mode is not set, an error will not be returned, and this variable will not be replaced in the source text.

If the variable is not found in the filter list and strict mode is set, an error will be returned.
