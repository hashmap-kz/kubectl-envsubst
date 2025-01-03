SOURCES := $(shell find . -name '*.go')
BINARY := kubectl-envsubst
COV_REPORT := "coverage.txt"

build: kubectl-envsubst

test: $(SOURCES)
	go test -v -short -race -timeout 30s ./...

test-cov:
	go test ./... -coverprofile=$(COV_REPORT)
	go tool cover -html=$(COV_REPORT)

test-integration:
	KUBECTL_ENVSUBST_INTEGRATION_TESTS_AVAILABLE=0xcafebabe go test -v test/integration/*.go

clean:
	@rm -rf $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o $(BINARY) -ldflags="-s -w" ./cmd/$(BINARY).go
