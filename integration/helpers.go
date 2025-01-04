package integration

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const (
	integrationTestEnv  = "KUBECTL_ENVSUBST_INTEGRATION_TESTS_AVAILABLE"
	integrationTestFlag = "0xcafebabe"
)

func randomIdent(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return strings.ToLower("kubectl-envsubst-" + string(b))
}

func cleanupResource(t *testing.T, kind, name string) {
	t.Logf("running: kubectl delete %s %s --ignore-not-found", kind, name)
	cmd := exec.Command("kubectl", "delete", kind, name, "--ignore-not-found")
	output, _ := cmd.CombinedOutput()
	t.Log(string(output))
}

func createNs(t *testing.T, ns string) {
	template := `
---
apiVersion: v1
kind: Namespace
metadata:
  name: "%s"
`
	nsp := fmt.Sprintf(template, ns)
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader([]byte(nsp))
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))
}

func setContextNs(t *testing.T, ns string) {
	cmd := exec.Command("kubectl", "config", "set-context", "--current", "--namespace", ns)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))
}

func getEnvsubstPath(t *testing.T) string {
	path, err := exec.LookPath("kubectl-envsubst")
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func printEnvsubstVersionInfo(t *testing.T) {
	t.Logf("*** kubectl-envsubst path: %s ***", getEnvsubstPath(t))

	cmd := exec.Command("kubectl", "krew", "info", "envsubst")
	output, err := cmd.CombinedOutput()
	// NOTE: '== nil'; in ci-jobs we're testing with 'make install'
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "VERSION") {
				t.Logf("*** kubectl-envsubst %s ***", line)
			}
		}
	}

}

func createTempFile(content string, extension string) (string, error) {

	tempFile, err := os.CreateTemp("", "kubectl-envsubst-tmp-*."+extension)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := tempFile.WriteString(content); err != nil {
		tempFile.Close()
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	return tempFile.Name(), nil
}
