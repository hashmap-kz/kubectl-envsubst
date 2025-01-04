package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

// NOTE: This package contains integration tests.
//
// Key Characteristics:
// 1. **No Mocks**: These tests interact with actual external dependencies rather than mocked versions.
//    The goal is to simulate real-world conditions as closely as possible.
// 2. **Sandboxed Environment**: To ensure safety and consistency, these tests are executed in a controlled and isolated environment (a kind-cluster sandbox).
//    This prevents unintended side effects on production systems, local configurations, or external resources.
//
// Safeguards:
// - Before each test run, appropriate safeguards are implemented to validate the environment and prevent harmful actions.
//
// Expectations:
// - These tests may take longer to execute compared to unit tests because they involve real components.
//

const (
	integrationTestEnv  = "KUBECTL_ENVSUBST_INTEGRATION_TESTS_AVAILABLE"
	integrationTestFlag = "0xcafebabe"
)

func TestApp_ApplyFromStdin(t *testing.T) {
	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	os.Args = []string{
		"kubectl-envsubst", "apply", "-f", "-",
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "mock-stdin-*")
	if err != nil {
		t.Fatal("Failed to create temporary file:", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the file afterwards

	// Write mock input to the temporary file
	mockInput := `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-c4a3d54857ef43398b9a557050a7c83c
data:
  key1: value1
`
	if _, err := tempFile.Write([]byte(mockInput)); err != nil {
		t.Fatal("Failed to write to temporary file:", err)
	}

	// Reset the file offset to the beginning for reading
	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		t.Fatal("Failed to seek temporary file:", err)
	}

	// Replace os.Stdin with the temporary file
	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()
	os.Stdin = tempFile

	// Create a temporary file for capturing stdout
	stdoutFile, err := os.CreateTemp("", "mock-stdout-*")
	if err != nil {
		t.Fatalf("Failed to create temporary stdout file: %v", err)
	}
	defer os.Remove(stdoutFile.Name()) // Clean up the temp file after the test

	// Replace os.Stdout with the temp file
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()
	os.Stdout = stdoutFile

	// Run application
	err = runApp()
	if err != nil {
		t.Fatal(err)
	}

	// Flush stdout and reset pointer for reading
	os.Stdout.Sync()
	stdoutFile.Seek(0, io.SeekStart)

	// Read and validate stdout content
	output, err := io.ReadAll(stdoutFile)
	if err != nil {
		t.Fatalf("Failed to read from temporary stdout file: %v", err)
	}

	// capture expected output
	strOut := string(output)
	if !strings.Contains(strOut, "configmap/cm-c4a3d54857ef43398b9a557050a7c83c") {
		t.Errorf("expected 'configmap/cm-c4a3d54857ef43398b9a557050a7c83c', got: %s", strOut)
	}

	t.Log(strOut)
}
