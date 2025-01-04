package cmd

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestJoinFiles_StdinUrlFilename(t *testing.T) {
	// Create a temporary file to mimic stdin
	tmpFile, err := os.CreateTemp("", "mock_stdin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the mock input data to the temp file
	mockInput := `
---
apiVersion: v1
kind: Secret
metadata:
  name: test-secret
type: Opaque
stringData:
  pass: "admin"
`
	if _, err := tmpFile.WriteString(mockInput); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Reset the read pointer of the temp file
	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("Failed to reset temp file pointer: %v", err)
	}

	// Replace os.Stdin with the temporary file
	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()
	os.Stdin = tmpFile

	files, err := JoinFiles([]string{
		"../../testdata/immutable_data/pod.yaml",
		"https://raw.githubusercontent.com/hashmap-kz/kubectl-envsubst/refs/heads/master/testdata/immutable_data/configmap.yaml",
	}, true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(files) == 0 {
		t.Error("Unexpected empty buffer")
	}

	content := string(files)
	if !strings.Contains(content, "name: nginx-container") {
		t.Error("Expecting pod manifests is read")
	}
	if !strings.Contains(content, "name: test-secret") {
		t.Error("Expecting secret manifests is read")
	}
	if !strings.Contains(content, "executor = LocalExecutor") {
		t.Error("Expecting configmap manifests is read")
	}
}

func TestReadRemote(t *testing.T) {

	t.Run("Successful Request", func(t *testing.T) {
		mockHTTPResponse := "Remote file content"
		http.DefaultClient = &http.Client{
			Transport: roundTripper(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(mockHTTPResponse)),
				}
			}),
		}
		result, err := readRemote("http://example.com/data")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if string(result) != mockHTTPResponse {
			t.Errorf("Expected: %s, Got: %s", mockHTTPResponse, string(result))
		}
	})

	t.Run("Failed Request", func(t *testing.T) {
		http.DefaultClient = &http.Client{
			Transport: roundTripper(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 404,
					Body:       io.NopCloser(strings.NewReader("Not Found")),
				}
			}),
		}
		_, err := readRemote("http://example.com/not-found")
		if err == nil {
			t.Error("Expected error but got none")
		}
	})
}

// Helper for mocking http.Client
type roundTripper func(req *http.Request) *http.Response

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
