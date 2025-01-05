package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestEnvsubstIntegration_NoSubst_Recursive(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	printEnvsubstVersionInfo(t)

	namespaceName := "kubectl-envsubst-plain-recursive"
	createNs(t, namespaceName)

	// Setup context
	setContextNs(t, namespaceName)

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl",
		"envsubst",
		"apply",
		"-f",
		"immutable_data/resolve/plain-recursive",
		"--recursive")

	output, err := cmdEnvsubstApply.CombinedOutput()
	stringOutput := string(output)
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, stringOutput)
	}
	t.Logf("\n%s\n", strings.TrimSpace(stringOutput))

	expectResources := []string{
		"cm-1-level-0",
		"cm-2-level-0",
		"cm-3-level-0",
		"cm-1-level-1",
		"cm-2-level-1",
		"cm-1-level-2",
		"cm-2-level-2",
	}

	for _, er := range expectResources {
		expectedOutput := strings.Contains(stringOutput, er)
		if !expectedOutput {
			t.Errorf("Expected substituted output to contain '%s', got %s", er, stringOutput)
		}
	}

}
