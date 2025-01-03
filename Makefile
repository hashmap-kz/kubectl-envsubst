# Variables
SOURCES := $(shell find . -name '*.go')
BINARY := kubectl-envsubst
COV_REPORT := coverage.txt
TEST_FLAGS := -v -race -timeout 30s

# Default target
.PHONY: all
all: build

# Build the binary (GOARCH=amd64 GOOS=linux; -o $(BINARY))
.PHONY: build
build: $(SOURCES)
	CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/$(BINARY).go

# Run unit tests
.PHONY: test
test:
	go test $(TEST_FLAGS) ./...

# Run tests with coverage
.PHONY: test-cov
test-cov:
	go test -coverprofile=$(COV_REPORT) ./...
	go tool cover -html=$(COV_REPORT)

# Run integration tests
.PHONY: test-integration
test-integration:
	KUBECTL_ENVSUBST_INTEGRATION_TESTS_AVAILABLE=0xcafebabe go test -v integration/*.go

# Lint the code
.PHONY: lint
lint:
	golangci-lint run ./...

# Format the code
.PHONY: format
format:
	go fmt ./...

# Clean build artifacts
.PHONY: clean
clean:
	@rm -rf $(BINARY) $(COV_REPORT)
