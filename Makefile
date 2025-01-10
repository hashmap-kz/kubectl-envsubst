# Variables
SOURCES := $(shell find . -name '*.go')
BINARY := kubectl-envsubst
COV_REPORT := coverage.txt
TEST_FLAGS := -v -race -timeout 30s
KIND_CLUSTER_NAME := kubectl-envsubst
INSTALL_DIR := /usr/local/bin

# Default target
.PHONY: all
all: build

# Build the binary (GOARCH=amd64 GOOS=linux; -o $(BINARY))
.PHONY: build
build: $(SOURCES)
	CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/kubectl-envsubst

# Install the binary to /usr/local/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY) to $(INSTALL_DIR)..."
	@install -m 0755 $(BINARY) $(INSTALL_DIR)

# Run unit tests
.PHONY: test
test:
	go test $(TEST_FLAGS) ./...

# Run tests with coverage
.PHONY: test-cov
test-cov:
	go test -coverprofile=$(COV_REPORT) -coverpkg=./pkg/...,!./integration/... ./...
	go tool cover -html=$(COV_REPORT)

# Setup kind-cluster (for running integration tests in a sandbox)
.PHONY: kind-setup
kind-setup:
	kind create cluster --name $(KIND_CLUSTER_NAME)
	kubectl config set-context kind-$(KIND_CLUSTER_NAME)

# Cleanup kind-cluster
.PHONY: kind-teardown
kind-teardown:
	kind delete clusters $(KIND_CLUSTER_NAME)

# Run integration tests
.PHONY: test-integration
test-integration: install kind-setup
	KUBECTL_ENVSUBST_INTEGRATION_TESTS_AVAILABLE=0xcafebabe go test -v integration/*.go
	KUBECTL_ENVSUBST_INTEGRATION_TESTS_AVAILABLE=0xcafebabe go test -v cmd/kubectl-envsubst/*.go
	$(MAKE) kind-teardown

# Lint the code
.PHONY: lint
lint:
	golangci-lint run ./...

# Check goreleaser
.PHONY: snapshot
snapshot:
	goreleaser release --skip sign --skip publish --snapshot --clean

# Format the code
.PHONY: format
format:
	go fmt ./...

# Clean build artifacts
.PHONY: clean
clean:
	@rm -rf $(BINARY) $(COV_REPORT)
