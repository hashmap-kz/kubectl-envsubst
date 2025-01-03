# **kubectl-envsubst**

_A `kubectl` plugin for substituting environment variables in Kubernetes manifests before applying them._

[![License](https://img.shields.io/github/license/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hashmap-kz/kubectl-envsubst)](https://goreportcard.com/report/github.com/hashmap-kz/kubectl-envsubst)
[![GitHub Issues](https://img.shields.io/github/issues/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/issues)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/blob/master/go.mod#L3)
[![Latest Release](https://img.shields.io/github/v/release/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/releases/latest)

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
    - [Using Krew](#using-krew)
    - [Manual Installation](#manual-installation)
- [Flags](#flags)
    - [Configure Using Command-Line Flags](#1-configure-using-command-line-flags)
    - [Configure Using Environment Variables](#2-configure-using-environment-variables)
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
- **Zero Dependencies**: No external tools or libraries are required, simplifying installation and operation.

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
3. Verify installation:
   ```bash
   kubectl envsubst --version
   ```

---

## **Flags**:

### **1. Configure Using Command-Line Flags**

#### `--envsubst-allowed-vars`

- **Description**: Specifies a comma-separated list of variable names that are explicitly allowed for substitution.
- **Usage**:
  ```bash
  kubectl envsubst apply -f deployment.yaml \
    --envsubst-allowed-vars=HOME,USER,PKEY_PATH,DB_PASS,IMAGE_NAME,IMAGE_TAG
  ```
- **Behavior**:
    - Variables not included in this list will not be substituted.
    - Useful for ensuring only specific variables are processed, preventing accidental substitutions.

#### `--envsubst-allowed-prefixes`

- **Description**: Specifies a comma-separated list of prefixes to filter variables by name.
- **Usage**:
  ```bash
  kubectl envsubst apply -f deployment.yaml \
    --envsubst-allowed-prefixes=APP_,CI_
  ```
- **Behavior**:
    - Only variables with names starting with one of the specified prefixes will be substituted.
    - Variables without matching prefixes will be ignored.

---

### **2. Configure Using Environment Variables**

#### `ENVSUBST_ALLOWED_VARS`

- **Description**: Alternative way to define a list of allowed variables for substitution.
- **Usage**:
  ```bash
  export ENVSUBST_ALLOWED_VARS='HOST,USER,PKEY_PATH'
  kubectl envsubst apply -f deployment.yaml
  ```

#### `ENVSUBST_ALLOWED_PREFIXES`

- **Description**: Alternative way to define a list of allowed prefixes for variable substitution.
- **Usage**:
  ```bash
  export ENVSUBST_ALLOWED_PREFIXES='CI_,APP_'
  kubectl envsubst apply -f deployment.yaml
  ```

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

Configure CLI for Prefix Substitution:

```bash
export ENVSUBST_ALLOWED_PREFIXES='APP_,IMAGE_'
```

Apply resources in dry-run mode to see the expected output before applying:

```bash
kubectl envsubst apply -f deployment.yaml --dry-run=client -o yaml
```

Recursively process all files in a directory while using dry-run mode:

```bash
kubectl envsubst apply -f manifests/ --recursive --dry-run=client -o yaml
```

Use redirection from stdin to apply the manifest:

```bash
cat deployment.yaml | kubectl envsubst apply -f -
```

Process and apply a manifest located on a remote server:

```bash
kubectl envsubst apply \
  -f https://raw.githubusercontent.com/user/repo/refs/heads/master/manifests/deployment.yaml
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
    - apk update && apk add --no-cache bash curl jq
    # setup kubectl
    - curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
    - chmod +x ./kubectl
    - cp ./kubectl /usr/local/bin
    # setup kubectl-envsubst plugin
    - wget -O- "https://github.com/hashmap-kz/kubectl-envsubst/releases/download/v1.0.14/kubectl-envsubst_v1.0.14_linux_amd64.tar.gz" | \
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

#### **Behavior**:

1. **Variables in Filters (Allowed for Substitution):**
    - If a variable is included in either the `--envsubst-allowed-vars` list or matches a prefix in the
      `--envsubst-allowed-prefixes` list but remains unexpanded:
        - This will result in an **error** during the substitution process.
        - The error occurs to ensure that all explicitly allowed variables are resolved.

2. **Variables Not in Filters (Not Allowed for Substitution):**
    - If a variable is not included in the `--envsubst-allowed-vars` list and does not match any prefixes in the
      `--envsubst-allowed-prefixes` list:
        - The variable will remain unexpanded.
        - This will not trigger an error during the substitution process but **may cause an error during the deploying
          of the manifest** if the unresolved placeholder is incompatible with Kubernetes.
        - In this example, placeholders within annotations are
          not expanded unless explicitly allowed via `--envsubst-allowed-vars` or `--envsubst-allowed-prefixes`.
          And this behavior is exactly what is expected.
          ```yaml
          some.controller.annotation/snippet: |
            set $agentflag 0;
            if ($http_user_agent ~* "(Android|iPhone|Windows Phone|UC|Kindle)" ) {
              set $agentflag 1;
            }
            if ( $agentflag = 1 ) {
              return 301 http://m.company.org;
            }
          ```

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
