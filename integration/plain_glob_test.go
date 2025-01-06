package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestEnvsubstIntegration_NoSubst_GlobPatterns_Yaml(t *testing.T) {
	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	printEnvsubstVersionInfo(t)

	namespaceName := "kubectl-envsubst-plain-glob"
	createNs(t, namespaceName)

	// Setup context
	setContextNs(t, namespaceName)

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/resolve/plain-glob/*.yaml")
	output, err := cmdEnvsubstApply.CombinedOutput()
	stringOutput := string(output)
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, stringOutput)
	}
	t.Logf("\n%s\n", strings.TrimSpace(stringOutput))

	expectResources := []string{
		"cm-1-level-0", // *.yaml
		"cm-2-level-0", // *.yaml
	}

	unExpectResources := []string{
		"cm-3-level-0", // *.json, should be ignored
	}

	for _, er := range expectResources {
		expectedOutput := strings.Contains(stringOutput, er)
		if !expectedOutput {
			t.Errorf("Expected substituted output to contain '%s', got %s", er, stringOutput)
		}
	}

	for _, er := range unExpectResources {
		unExpectedOutput := strings.Contains(stringOutput, er)
		if unExpectedOutput {
			t.Errorf("Expected substituted output does not contain '%s'", er)
		}
	}
}

func TestEnvsubstIntegration_NoSubst_GlobPatterns_Json(t *testing.T) {
	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	printEnvsubstVersionInfo(t)

	namespaceName := "kubectl-envsubst-plain-glob"
	createNs(t, namespaceName)

	// Setup context
	setContextNs(t, namespaceName)

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/resolve/plain-glob/*.json")
	output, err := cmdEnvsubstApply.CombinedOutput()
	stringOutput := string(output)
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, stringOutput)
	}
	t.Logf("\n%s\n", strings.TrimSpace(stringOutput))

	expectResources := []string{
		"cm-3-level-0", // *.json
	}

	unExpectResources := []string{
		"cm-1-level-0", // *.yaml
		"cm-2-level-0", // *.yaml
	}

	for _, er := range expectResources {
		expectedOutput := strings.Contains(stringOutput, er)
		if !expectedOutput {
			t.Errorf("Expected substituted output to contain '%s', got %s", er, stringOutput)
		}
	}

	for _, er := range unExpectResources {
		unExpectedOutput := strings.Contains(stringOutput, er)
		if unExpectedOutput {
			t.Errorf("Expected substituted output does not contain '%s'", er)
		}
	}
}
