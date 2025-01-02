# **kubectl-envsubst**

_A `kubectl` plugin for substituting environment variables in Kubernetes manifests before applying them._

[![License](https://img.shields.io/github/license/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hashmap-kz/kubectl-envsubst)](https://goreportcard.com/report/github.com/hashmap-kz/kubectl-envsubst)
[![GitHub Issues](https://img.shields.io/github/issues/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/issues)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/blob/master/go.mod#L3)
[![Latest Release](https://img.shields.io/github/v/release/hashmap-kz/kubectl-envsubst)](https://github.com/hashmap-kz/kubectl-envsubst/releases/latest)

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

## **Usage**

### Basic Variable Substitution

```bash
export IMAGE_TAG=v1.2.3

kubectl envsubst apply -f manifests/ \
    --envsubst-allowed-vars=IMAGE_TAG
```

### Using Prefixes

```bash
export APP_PROJECT_NAME=auth-svc
export APP_PROJECT_NAMESPACE=trade-system-prod
export INFRA_DOMAIN_NAME=company.org

kubectl envsubst apply -f manifests/ \
    --envsubst-allowed-prefixes=APP_,INFRA_
```

### Mixed mode

```bash
export APP_PROJECT_NAME=auth-svc
export APP_PROJECT_NAMESPACE=trade-system-prod
export INFRA_DOMAIN_NAME=company.org
export HOST=svc.company.org
export PORT=1024

kubectl envsubst apply -f manifests/ \
    --envsubst-allowed-prefixes=APP_,INFRA_ \
    --envsubst-allowed-vars=HOST,PORT
```

---

## **Flags**:

```
--envsubst-allowed-vars=HOME,USER,PKEY_PATH,DB_PASS,IMAGE_NAME,IMAGE_TAG
    Accepts a comma-separated list of variable names allowed for substitution. 
    Variables not included in this list will not be substituted.

--envsubst-allowed-prefixes=APP_,CI_ 
    Accepts a comma-separated list of prefixes. 
    Only variables with names starting with one of these prefixes will be substituted; others will be ignored.
    
```

---

## **Variable Expansion and Filtering Behavior**

### **Description**

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

### **Behavior**:

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

## **Usage scenario in CI/CD**

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

We welcome contributions! To contribute:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request describing your changes.

---

## **License**

This project is licensed under the [MIT License](LICENSE).

---

## **Additional Resources**

- [Kubernetes Documentation: Managing Resources](https://kubernetes.io/docs/concepts/cluster-administration/manage-deployment/)
- [Using Environment Variables in CI/CD](https://12factor.net/config)

For more information, visit the [project repository](https://github.com/hashmap-kz/kubectl-envsubst).

---
