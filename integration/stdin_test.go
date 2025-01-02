package integration

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestEnvsubstIntegrationFromStdin(t *testing.T) {

	if os.Getenv("KUBECTL_ENVSUBST_INTEGRATION_TESTS_AVAILABLE") != "KUBECTL_ENVSUBST_INTEGRATION_TESTS_AVAILABLE" {
		log.Printf("Integration test was skipped due to configuration")
		return
	}

	// Setup environment variables that was used in substitution
	os.Setenv("IMAGE_NAME", "nginx")
	os.Setenv("IMAGE_TAG", "latest")
	os.Setenv("CI_PROJECT_NAME", "kubectl-envsubst-integration-test")
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
	cmd := exec.Command("kubectl", "envsubst", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(manifest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, string(output))
	}

	// Check result (it should be created/updated/unchanged, etc...)
	expectedOutput := strings.Contains(string(output), "deployment.apps/kubectl-envsubst-integration-test")
	if !expectedOutput {
		t.Errorf("Expected substituted output to contain 'deployment.apps/kubectl-envsubst-integration-test', got %s", string(output))
	}

	// Validate applied resource
	validateCmd := exec.Command("kubectl", "get", "deployment", "kubectl-envsubst-integration-test")
	validateOutput, err := validateCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to validate applied resource: %v, output: %s", err, string(validateOutput))
	}
	if !strings.Contains(string(validateOutput), "kubectl-envsubst-integration-test") {
		t.Errorf("Expected deployment 'kubectl-envsubst-integration-test' to exist, got %s", string(validateOutput))
	}
}
