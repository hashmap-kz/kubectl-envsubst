package integration

import (
	"bufio"
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

func TestEnvsubstIntegrationFromFile(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	resourceName := randomIdent(32)
	defer cleanupResource(t, "deployment", resourceName)

	// Setup environment variables that was used in substitution
	os.Setenv("IMAGE_NAME", "nginx")
	os.Setenv("IMAGE_TAG", "latest")
	os.Setenv("CI_PROJECT_NAME", resourceName)
	defer os.Unsetenv("IMAGE_NAME")
	defer os.Unsetenv("IMAGE_TAG")
	defer os.Unsetenv("CI_PROJECT_NAME")

	// configure CLI
	os.Setenv("ENVSUBST_ALLOWED_PREFIXES", "CI_,IMAGE_")
	defer os.Unsetenv("ENVSUBST_ALLOWED_PREFIXES")

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/01_deployment.yaml")
	output, err := cmdEnvsubstApply.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, string(output))
	}
	t.Log(string(output))

	// Check result (it should be created/updated/unchanged, etc...)
	expectedOutput := strings.Contains(string(output), fmt.Sprintf("deployment.apps/%s", resourceName))
	if !expectedOutput {
		t.Errorf("Expected substituted output to contain 'deployment.apps/%s', got %s", resourceName, string(output))
	}

	// Validate applied resource
	validateCmd := exec.Command("kubectl", "get", "deployment", resourceName)
	validateOutput, err := validateCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to validate applied resource: %v, output: %s", err, string(validateOutput))
	}
	if !strings.Contains(string(validateOutput), resourceName) {
		t.Errorf("Expected deployment 'kubectl-envsubst-integration-test' to exist, got %s", string(validateOutput))
	}

}

func TestEnvsubstIntegrationFromStdin(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	resourceName := randomIdent(32)
	defer cleanupResource(t, "deployment", resourceName)

	// Setup environment variables that was used in substitution
	os.Setenv("IMAGE_NAME", "nginx")
	os.Setenv("IMAGE_TAG", "latest")
	os.Setenv("CI_PROJECT_NAME", resourceName)
	defer os.Unsetenv("IMAGE_NAME")
	defer os.Unsetenv("IMAGE_TAG")
	defer os.Unsetenv("CI_PROJECT_NAME")

	// configure CLI
	os.Setenv("ENVSUBST_ALLOWED_PREFIXES", "CI_,IMAGE_")
	defer os.Unsetenv("ENVSUBST_ALLOWED_PREFIXES")

	// Prepare input manifest
	manifest := `
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
          image: $IMAGE_NAME:$IMAGE_TAG
          imagePullPolicy: Always
`

	// Run kubectl-envsubst
	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "-")
	cmdEnvsubstApply.Stdin = strings.NewReader(manifest)
	output, err := cmdEnvsubstApply.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, string(output))
	}
	t.Log(string(output))

	// Check result (it should be created/updated/unchanged, etc...)
	expectedOutput := strings.Contains(string(output), fmt.Sprintf("deployment.apps/%s", resourceName))
	if !expectedOutput {
		t.Errorf("Expected substituted output to contain 'deployment.apps/%s', got %s", resourceName, string(output))
	}

	// Validate applied resource
	validateCmd := exec.Command("kubectl", "get", "deployment", resourceName)
	validateOutput, err := validateCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to validate applied resource: %v, output: %s", err, string(validateOutput))
	}
	if !strings.Contains(string(validateOutput), resourceName) {
		t.Errorf("Expected deployment 'kubectl-envsubst-integration-test' to exist, got %s", string(validateOutput))
	}

}

func TestEnvsubstIntegrationFromUrl(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("Integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	const url = "https://raw.githubusercontent.com/hashmap-kz/kubectl-envsubst/refs/heads/integ/integration/immutable_data/01_deployment.yaml"
	resourceName := randomIdent(32)
	defer cleanupResource(t, "deployment", resourceName)

	// Setup environment variables that was used in substitution
	os.Setenv("IMAGE_NAME", "nginx")
	os.Setenv("IMAGE_TAG", "latest")
	os.Setenv("CI_PROJECT_NAME", resourceName)
	defer os.Unsetenv("IMAGE_NAME")
	defer os.Unsetenv("IMAGE_TAG")
	defer os.Unsetenv("CI_PROJECT_NAME")

	// configure CLI
	os.Setenv("ENVSUBST_ALLOWED_PREFIXES", "CI_,IMAGE_")
	defer os.Unsetenv("ENVSUBST_ALLOWED_PREFIXES")

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", url)
	output, err := cmdEnvsubstApply.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, string(output))
	}
	t.Log(string(output))

	// Check result (it should be created/updated/unchanged, etc...)
	expectedOutput := strings.Contains(string(output), fmt.Sprintf("deployment.apps/%s", resourceName))
	if !expectedOutput {
		t.Errorf("Expected substituted output to contain 'deployment.apps/%s', got %s", resourceName, string(output))
	}

	// Validate applied resource
	validateCmd := exec.Command("kubectl", "get", "deployment", resourceName)
	validateOutput, err := validateCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to validate applied resource: %v, output: %s", err, string(validateOutput))
	}
	if !strings.Contains(string(validateOutput), resourceName) {
		t.Errorf("Expected deployment 'kubectl-envsubst-integration-test' to exist, got %s", string(validateOutput))
	}

}

// helpers

func randomIdent(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return strings.ToLower("I" + string(b))
}

func cleanupResource(t *testing.T, kind, name string) {
	t.Logf("running: kubectl delete %s %s --ignore-not-found", kind, name)
	cmd := exec.Command("kubectl", "delete", kind, name, "--ignore-not-found")
	output, _ := cmd.CombinedOutput()
	t.Log(string(output))
}

func printEnvsubstVersionInfo(t *testing.T) {
	cmd := exec.Command("kubectl", "krew", "info", "envsubst")
	output, _ := cmd.CombinedOutput()

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "VERSION") {
			t.Logf("*** kubectl-envsubst %s ***", line)
		}
	}
}
