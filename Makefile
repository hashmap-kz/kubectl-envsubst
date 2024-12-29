.PHONY: all kubectl-envsubst test

all: kubectl-envsubst

kubectl-envsubst: *.go
	go build -ldflags="-s -w" .

test:
	go test ./... -coverprofile=cover.out -v

run:
	go run ./main.go

run-linter:
	echo "Starting linters"
	golangci-lint run ./...
