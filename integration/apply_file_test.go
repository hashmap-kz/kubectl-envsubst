package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestEnvsubstIntegration_SubstApplyFromFile(t *testing.T) {
	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

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
		t.Errorf("Expected resource %s to exist, got %s", resourceName, string(validateOutput))
	}
}
