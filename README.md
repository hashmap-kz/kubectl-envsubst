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

```bash
# substitute variables whose names start with one of the prefixes
kubectl envsubst apply -f manifests/ --envsubst-allowed-prefixes=CI_,APP_

# substitute well-defined variables
kubectl envsubst apply -f manifests/ \
    --envsubst-allowed-vars=CI_PROJECT_NAME,CI_COMMIT_REF_NAME,APP_IMAGE

# mixed mode, check both full match and prefix match
kubectl envsubst apply -f manifests/ \
    --envsubst-allowed-prefixes=CI_,APP_ \
    --envsubst-allowed-vars=HOME,USER

# example:
export APP_NAME=nginx
export APP_IMAGE_NAME=nginx
export APP_IMAGE_TAG=latest
kubectl envsubst apply -f testdata/subst/01.yaml \
    --dry-run=client -oyaml \
    --envsubst-allowed-prefixes=APP_
```

### Flags:

```
--envsubst-allowed-vars 
    Description: 
        Consumes comma-separated list of names that allowed for expansion.
        If a variable not in allowed list, it won't be expanded.
    Example: --envsubst-allowed-prefixes=APP_,CI_

--envsubst-allowed-prefixes 
    Description: 
        Consumes comma-separated list of prefixes.
        Variables whose names do not start with any prefix will be ignored.
    Example: --envsubst-allowed-vars=HOME,USER,PKEY_PATH,DB_PASS,IMAGE_NAME,IMAGE_TAG

--envsubst-no-strict
    Description: 
        Strict mode is ON by default. 
        In 99% of cases, this is exactly what is required.
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

### Usage scenario in CI/CD

A typical setup (with dev/stage/prod environments) for a typical microservice may look like this:

```
.
├── cmd
│   └── auth-svc.go
├── go.mod
├── go.sum
├── .gitlab-ci.yml
├── k8s-manifests
│   ├── dev
│   │   └── manifests.yaml
│   ├── prod
│   │   ├── hpa.yaml
│   │   ├── manifests.yaml
│   │   └── secrets.yaml
│   └── stage
│       └── manifests.yaml
├── README.md
```

Where for each environment (dev, stage, prod) there are a bunch of k8s-manifests.

They may vary between environments just in naming (image name in registry), or may contain specific resources (hpa and secrets for prod).

But the main part is duplicated (service, deployment, ingress, secrets, labels, naming conventions, etc...).

And all these duplicated parts may be substituted from env-vars in a pipeline you use.

The tool that you use for CI/CD does not matter, each tool may provide different var-names and patterns.

In this example, a project name, its path, labels are the same for all environments.

Only the image name will be different, and perhaps you may also include specific resources, depends on your needs.

```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: &app ${CI_PROJECT_NAME}
  labels:
    app: *app
spec:
  ports:
    - port: 8080
      targetPort: 8080
      name: http
  selector:
    app: *app

---
apiVersion: v1
kind: Secret
metadata:
  name: &app ${CI_PROJECT_NAME}
  labels:
    app: *app
type: Opaque
stringData:
  vault_path: "secret/${CI_PROJECT_PATH}/${CI_COMMIT_REF_NAME}"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: &app ${CI_PROJECT_NAME}
  labels:
    app: *app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: *app
  template:
    metadata:
      labels:
        app: *app
    spec:
      containers:
        - name: *app
          # image name for each environment is different (dev, stage, prod)
          image: ${APP_IMAGE}
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
              name: http
          envFrom:
            - secretRef:
                name: *app
```

Let's assume that we deploy our application by using gitlab-ci (the tool does not matter, just an example).

This bash script is used as an example that illustrates if we have env-vars in the context we deploy our application - then we may use these vars in our manifests.

```bash
# gitlab-specific variables, that defined in each pipeline 
# we simulate that pipeline by setting them by hand
export CI_REGISTRY=mirror-0.company.org:5000
export CI_PROJECT_ROOT_NAMESPACE=banking-system-envsubst-test
export CI_PROJECT_NAMESPACE=backend/auth-svc
export CI_PROJECT_PATH=banking-system-envsubst-test/backend/auth-svc
export CI_PROJECT_NAME=auth-svc
export CI_COMMIT_REF_NAME=dev

# project-specific variables (may be passed to pipeline from secrets, project/group variables, combined from other vars, etc...)
# nginx image is used for simplicity, of course there will be your images, that was pushed to your repo in build-step
export APP_IMAGE=nginx
export APP_NAMESPACE="${CI_PROJECT_ROOT_NAMESPACE}-${CI_COMMIT_REF_NAME}"
export INFRA_DOMAIN_NAME=company.org

# prepare namespace and context
kubectl create ns "${APP_NAMESPACE}" --dry-run=client -oyaml | kubectl apply -f -
kubectl config set-context --current --namespace="${APP_NAMESPACE}"

# expand and apply manifests
kubectl envsubst apply -f "k8s-manifests/${CI_COMMIT_REF_NAME}" \
  --envsubst-allowed-prefixes=CI_,APP_,INFRA_
```

The CI/CD stage may look like this:

```yaml
deploy:
  stage: deploy
  before_script:
    - apk update && apk add --no-cache bash curl jq
    # setup kubectl
    - curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
    - chmod +x ./kubectl
    - cp ./kubectl /usr/local/bin
    # setup kubectl-envsubst plugin
    - wget -O- "https://github.com/hashmap-kz/kubectl-envsubst/releases/download/v1.0.6/kubectl-envsubst_v1.0.6_linux_amd64.tar.gz" | \
      tar -xzf - -C /usr/local/bin && chmod +x /usr/local/bin/kubectl-envsubst
  tags:
    - dind
  environment:
    name: $CI_COMMIT_REF_NAME
  script:
    - kubectl create ns "${APP_NAMESPACE}" --dry-run=client -oyaml | kubectl apply -f -
    - kubectl config set-context --current --namespace="${APP_NAMESPACE}"
    - kubectl envsubst apply -f "k8s-manifests/${CI_COMMIT_REF_NAME}" --envsubst-allowed-prefixes=CI_,APP_,INFRA_
    - kubectl rollout restart deploy "${CI_PROJECT_NAME}"
```

> As a result - we use plain kubernetes manifests, without any 'magic-preprocessing-tricks', and environment variables for
> handling application deploy in different environments (dev, stage, prod, etc...).
> 
> The main point - that we're totally sure that the variable substitution is managed and predictable.
> 
> If our application contains a configmap for nginx (that uses $ sign a LOT), we're absolutely sure that there will be
> no substitutions, if we're not allow it explicitly by adding collide variables in allow-list of prefix-list.

