# **kubectl-envsubst**

_A `kubectl` plugin for substituting environment variables in Kubernetes manifests before applying them._

[![License](https://img.shields.io/github/license/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hashmap-kz/kubectl-envsubst)](https://goreportcard.com/report/github.com/hashmap-kz/kubectl-envsubst)
[![Workflow Status](https://img.shields.io/github/actions/workflow/status/hashmap-kz/kubectl-envsubst/ci.yml?branch=master)](https://github.com/hashmap-kz/kubectl-envsubst/actions/workflows/ci.yml?query=branch:master)
[![codecov](https://codecov.io/gh/hashmap-kz/kubectl-envsubst/branch/master/graph/badge.svg)](https://codecov.io/gh/hashmap-kz/kubectl-envsubst)
[![GitHub Issues](https://img.shields.io/github/issues/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/issues)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/blob/master/go.mod#L3)
[![Latest Release](https://img.shields.io/github/v/release/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/releases/latest)

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
    - [Using Krew](#using-krew)
    - [Manual Installation](#manual-installation)
    - [Package-Based Installation](#package-based-installation-for-cicd-pipelines-example-for-alpine-linux)
- [Flags](#flags)
- [Usage](#usage-examples)
    - [Basic Usage](#basic-substitution-example)
    - [Substitution Along with Other `kubectl apply` Options](#substitution-along-with-other-kubectl-apply-options)
    - [Advanced Usage](#advanced-usage-typical-scenario-in-cicd)
- [Implementation details](#implementation-details)
    - [Variable expansion behaviour](#variable-expansion-and-filtering-behavior)
- [Brief conclusion](#brief-conclusion)
- [Contributing](#contributing)
- [License](#license)
- [Additional resources](#additional-resources)

---

## **Features**

- **Environment Variable Substitution**: Replaces placeholders in Kubernetes manifests with environment variable values.
- **Allowed Variable Filtering**: Allows you to control substitutions by specifying prefixes or permitted
  variables.
- **Strict Mode**: Fails if any placeholders remain unexpanded, ensuring deployment predictability.
- **Seamless Integration**: Works with all `kubectl` arguments, that may be used in `apply`.
- **Zero Dependencies**: No external tools or libraries are required, simplifying installation and operation. The plugin
  is designed primarily for CI/CD, with a focus on minimizing binary size.

---

## **Installation**

### Using `krew`

1. Install the [Krew](https://krew.sigs.k8s.io/docs/user-guide/setup/) plugin manager if you haven’t already.
2. Run the following command:
   ```bash
   kubectl krew install envsubst
   ```
3. Verify installation:
   ```bash
   kubectl envsubst --version
   ```

### Manual Installation

1. Download the latest binary for your platform from
   the [Releases page](https://github.com/hashmap-kz/kubectl-envsubst/releases).
2. Place the binary in your system's `PATH` (e.g., `/usr/local/bin`).
3. Example installation script for Unix-Based OS:
   ```bash
   (
     set -euo pipefail

     OS="$(uname | tr '[:upper:]' '[:lower:]')"
     ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')"
     TAG="$(curl -s https://api.github.com/repos/hashmap-kz/kubectl-envsubst/releases/latest | jq -r .tag_name)"

     curl -L "https://github.com/hashmap-kz/kubectl-envsubst/releases/download/${TAG}/kubectl-envsubst_${TAG}_${OS}_${ARCH}.tar.gz" |
       tar -xzf - -C /usr/local/bin && chmod +x /usr/local/bin/kubectl-envsubst
   )
   ```
4. Verify installation:
   ```bash
   kubectl envsubst --version
   ```

### Package-Based Installation (for CI/CD pipelines, example for Alpine-Linux)

```bash
apk update && apk add --no-cache bash curl
curl -LO https://github.com/hashmap-kz/kubectl-envsubst/releases/latest/download/kubectl-envsubst_linux_amd64.apk
apk add kubectl-envsubst_linux_amd64.apk --allow-untrusted
```

---

## **Flags**:

### **`--envsubst-allowed-vars`**

- **Description**: Specifies a comma-separated list of variable names that are explicitly allowed for substitution.
- **Corresponding environment variable**: **`ENVSUBST_ALLOWED_VARS`**
- **Usage**:
  ```bash
  # Using CLI options
  kubectl envsubst apply -f deployment.yaml \
    --envsubst-allowed-vars=IMAGE_NAME,IMAGE_TAG,APP_NAME,PKEY_PATH

  # Using environment variables
  export ENVSUBST_ALLOWED_VARS='IMAGE_NAME,IMAGE_TAG,APP_NAME,PKEY_PATH'
  kubectl envsubst apply -f deployment.yaml
  ```
- **Behavior**:
    - Variables not included in this list will not be substituted.
    - Useful for ensuring only specific variables are processed, preventing accidental substitutions.

---

### **`--envsubst-allowed-prefixes`**

- **Description**: Specifies a comma-separated list of prefixes to filter variables by name.
- **Corresponding environment variable**: **`ENVSUBST_ALLOWED_PREFIXES`**
- **Usage**:
  ```bash
  # Using CLI options
  kubectl envsubst apply -f deployment.yaml \
    --envsubst-allowed-prefixes=CI_,APP_,IMAGE_

  # Using environment variables
  export ENVSUBST_ALLOWED_PREFIXES='CI_,APP_,IMAGE_'
  kubectl envsubst apply -f deployment.yaml
  ```
- **Behavior**:
    - Only variables with names starting with one of the specified prefixes will be substituted.
    - Variables without matching prefixes will be ignored.

---

### **Note: CLI Takes Precedence Over Environment Variables**

- **Priority**: If both CLI flags and environment variables are set:
    - **CLI flags** (`--envsubst-allowed-vars`, `--envsubst-allowed-prefixes`) will **override** their respective
      environment variable values.
    - This ensures explicit command-line options have the highest priority.

---

## Usage Examples

### **Basic Substitution Example**

Given a manifest file `deployment.yaml`:

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${APP_NAME}
  labels:
    app: ${APP_NAME}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ${APP_NAME}
  template:
    metadata:
      labels:
        app: ${APP_NAME}
    spec:
      containers:
        - name: ${APP_NAME}
          image: ${IMAGE_NAME}:${IMAGE_TAG}
```

#### Export required environment variables:

```bash
export APP_NAME=my-app
export IMAGE_NAME=nginx
export IMAGE_TAG='1.27.3-bookworm'
```

#### Substitution with Allowed Variables

```bash
kubectl envsubst apply --filename=deployment.yaml \
  --envsubst-allowed-vars=APP_NAME,IMAGE_NAME,IMAGE_TAG
```

#### Substitution with Allowed Prefixes

```bash
kubectl envsubst apply --filename=deployment.yaml \
  --envsubst-allowed-prefixes=APP_,IMAGE_
```

---

### **Substitution Along with Other `kubectl apply` Options**

```bash
# Configure CLI for Prefix Substitution:
export ENVSUBST_ALLOWED_PREFIXES='APP_,IMAGE_'
```

```bash
# Apply resources in dry-run mode to see the expected output before applying:
kubectl envsubst apply -f deployment.yaml \
  --dry-run=client -o yaml
```

```bash
# Recursively process all files in a directory while using dry-run mode:
kubectl envsubst apply -f manifests/ --recursive \
  --dry-run=client -o yaml
```

```bash
# Use redirection from stdin to apply the manifest:
cat deployment.yaml | kubectl envsubst apply -f -
```

```bash
# Process and apply a manifest located on a remote server:
kubectl envsubst apply \
  -f https://raw.githubusercontent.com/user/repo/refs/heads/master/manifests/deployment.yaml
```

```bash
# Using with 'kubectl kustomize'
kubectl kustomize manifests/ | kubectl envsubst apply -f -
```

---

### **Advanced usage (typical scenario in CI/CD)**

A typical setup for a microservice with dev, stage, and prod environments may look like this:

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

Each environment (dev, stage, prod) typically includes a set of Kubernetes manifests.
These manifests may differ between environments in minor ways, such as image names in the registry,
or more significantly, with specific resources like HPAs or secrets for production.

However, most parts of the manifests, such as services, deployments, ingresses, secrets,
labels, and naming conventions, are often duplicated across environments.
These duplicated elements can easily be replaced using environment variables in your CI/CD pipeline.

The specific CI/CD tool you use doesn’t matter, as each tool may provide different
variable names and patterns. In this example, project details like name, path, and
labels remain consistent across environments, while variables like image names or
environment-specific resources can be adjusted based on your needs.

A set of application deployment manifests:

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

Let’s assume we’re deploying our application using GitLab CI (the specific tool doesn’t matter, it’s just an example).

The CI/CD stage may look like this:

```yaml
deploy:
  stage: deploy
  before_script:
    - apk update && apk add --no-cache bash curl
    # setup kubectl
    - curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
    - chmod +x ./kubectl && cp ./kubectl /usr/local/bin
    # setup kubectl-envsubst plugin (using latest release and *.apk package)
    - curl -LO https://github.com/hashmap-kz/kubectl-envsubst/releases/latest/download/kubectl-envsubst_linux_amd64.apk
    - apk add kubectl-envsubst_linux_amd64.apk --allow-untrusted
  tags:
    - dind
  environment:
    name: $CI_COMMIT_REF_NAME
  script:
    - export APP_NAMESPACE="${CI_PROJECT_ROOT_NAMESPACE}-${CI_COMMIT_REF_NAME}"
    - export ENVSUBST_ALLOWED_PREFIXES='CI_,APP_,INFRA_,IMAGE_'
    # create namespace, setup context
    - kubectl create ns "${APP_NAMESPACE}" --dry-run=client -oyaml | kubectl apply -f -
    - kubectl config set-context --current --namespace="${APP_NAMESPACE}"
    # substitute and apply manifests
    - kubectl envsubst apply -f "k8s-manifests/${CI_COMMIT_REF_NAME}"
    - kubectl rollout restart deploy "${CI_PROJECT_NAME}"
```

---

## **Implementation details**

### **Variable Expansion and Filtering Behavior**

#### **Description**

The behavior of variable substitution is determined by the inclusion of variables in the `--envsubst-allowed-vars` or
`--envsubst-allowed-prefixes` lists and whether the variables remain unresolved.

Substitution of environment variables without verifying their inclusion in a filter list is intentionally
avoided, as this behavior can lead to subtle errors.

If a variable is not found in the filter list, an error will be returned.

Expanding manifests with all available environment variables can work fine
for simple cases, such as when your manifest contains only a service and a deployment
with a few variables to substitute.
However, this approach becomes challenging to debug when dealing with more complex
scenarios, such as applying dozens of manifests involving ConfigMaps, Secrets, CRDs, and other resources.

In such cases, you may not have complete confidence that all substitutions were
performed as expected.
For this reason, it’s better to explicitly control which variables are substituted by using a filter list.

#### **Behavior**

1. **Variables Included in Filters (Allowed for Substitution):**
    - Variables listed in `--envsubst-allowed-vars` or matching a prefix in `--envsubst-allowed-prefixes`:
        - **If Unexpanded**: This will result in an **error** during the substitution process.
        - **Reason**: The error ensures all explicitly allowed variables are resolved to avoid deployment issues.

2. **Variables Excluded from Filters (Not Allowed for Substitution):**
    - Variables not listed in `--envsubst-allowed-vars` and not matching any prefix in `--envsubst-allowed-prefixes`:
        - **Behavior**:
            - The variable remains unexpanded.
            - This does **not trigger an error** during substitution.
            - Kubernetes deployment may fail if unresolved placeholders are incompatible with the manifest structure.

3. **Expected Behavior for Specific Use Cases:**
    - Certain placeholders, such as those in annotations, are intentionally **not expanded** unless explicitly allowed.
    - Example:
      ```yaml
      annotations:
        some.controller.annotation/snippet: |
          set $agentflag 0;
          if ($http_user_agent ~* "(Android|iPhone|Windows Phone|UC|Kindle)" ) {
            set $agentflag 1;
          }
          if ( $agentflag = 1 ) {
            return 301 http://m.company.org;
          }
      ```
    - In the above case, placeholders remain unchanged unless explicitly allowed
      via `--envsubst-allowed-vars` or `--envsubst-allowed-prefixes`.
      This ensures manifest consistency and aligns with expected behavior.

---

## **Brief conclusion**

We use plain Kubernetes manifests without any complex preprocessing, relying on environment variables to
manage deployments across different environments (dev, stage, prod, etc.).

This approach ensures variable substitution is controlled and predictable.
For example, if an application includes a ConfigMap for Nginx (which uses $ frequently),
substitutions won’t occur unless explicitly allowed by adding the relevant variables to
an allow-list or prefix-list.

---

## **Contributing**

We welcome contributions! To contribute: see the [Contribution](CONTRIBUTING.md) guidelines.

---

## **License**

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

## **Additional Resources**

- [Kubernetes Documentation: Managing Resources](https://kubernetes.io/docs/concepts/cluster-administration/manage-deployment/)
- [Using Environment Variables in CI/CD](https://12factor.net/config)

For more information, visit the [project repository](https://github.com/hashmap-kz/kubectl-envsubst).

---
