# kubectl-envsubst

[![License](https://img.shields.io/github/license/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hashmap-kz/kubectl-envsubst)](https://goreportcard.com/report/github.com/hashmap-kz/kubectl-envsubst)
[![GitHub Issues](https://img.shields.io/github/issues/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/issues)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/blob/master/go.mod#L3)
[![Latest Release](https://img.shields.io/github/v/release/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/releases/latest)

A kubectl plugin that substitutes (manageable and predictable) env-vars in k8s manifests before applying them.

### Features:

- Expand environment-variables in manifests passed to kubectl, before applying them
- Uses prefixes, and allowed variables (you're able to be totally sure what it's going to be substituted)
- Has a strict mode (if some variables are not expanded, it fails)
- Uses all known kubectl args (by just proxying them as is, without any additional actions)
- ZERO dependencies, and I mean literally zero

## Installation

Use [krew](https://krew.sigs.k8s.io/) plugin manager to install:

    kubectl krew install envsubst
    kubectl envsubst --help

Download the binary from [GitHub Releases](https://github.com/hashmap-kz/kubectl-envsubst/releases).

```bash
# Other available architectures are linux_arm64, darwin_amd64, darwin_arm64, windows_amd64.
export ARCH=linux_amd64
# Check the latest version, https://github.com/hashmap-kz/kubectl-envsubst/releases/latest
export VERSION=1.0.2
wget -O- "https://github.com/hashmap-kz/kubectl-envsubst/releases/download/v${VERSION}/kubectl-envsubst_${VERSION}_${ARCH}.tar.gz" | \
  sudo tar -xzf - -C /usr/local/bin && \
  sudo chmod +x /usr/local/bin/kubectl-envsubst
```

From source.

```bash
go install github.com/hashmap-kz/kubectl-envsubst@latest
sudo mv $GOPATH/bin/kubectl-envsubst /usr/local/bin
```

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

### Flags:

```
--envsubst-allowed-vars: consumes comma-separated list of names that allowed for expansion
    Example: --envsubst-allowed-prefixes=APP_,CI_

--envsubst-allowed-prefixes: consumes comma-separated list of prefixes, variables that not match will be ignored
    Example: --envsubst-allowed-vars=HOME,USER,PKEY_PATH,DB_PASS,IMAGE_NAME,IMAGE_TAG

--envsubst-no-strict: strict mode is ON by default, in 99% of cases, this is exactly what is required.
```

### Implementation details

Substitution of environmental variables without checking for inclusion in one of the filter list is not used
intentionally, because this behavior can lead to subtle mistakes.

If a variable is not found in the filter list and strict mode is not set, an error will not be returned, and this
variable will not be replaced in the source text.

If the variable is not found in the filter list and strict mode is set, an error will be returned.

It's totally fine - expanding manifests with all available env-vars, if your manifest contains a service and deployment
with a few vars to substitute.

But it will be a hard to debug, it you need to apply a few dozens manifests with config-maps, secrets, CRD's, etc...

It that case you won't be absolutely sure that everything will be placed as expected.

That's why it's better to control the amount of variables that may be substituted.
