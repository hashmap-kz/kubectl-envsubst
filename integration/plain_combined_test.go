package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestEnvsubstIntegration_NoSubst_MixedManifests(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	printEnvsubstVersionInfo(t)

	namespaceName := "kubectl-envsubst-plain-combined"
	createNs(t, namespaceName)

	// Setup context
	setContextNs(t, namespaceName)

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/resolve/plain-combined")
	output, err := cmdEnvsubstApply.CombinedOutput()
	stringOutput := string(output)
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, stringOutput)
	}
	t.Log(stringOutput)

	expectResources := []string{
		"serviceaccount/my-app",
		"role.rbac.authorization.k8s.io/my-app",
		"rolebinding.rbac.authorization.k8s.io/my-app",
		"configmap/my-app",
		"secret/my-app",
		"deployment.apps/my-app",
		"service/my-app",
	}

	for _, er := range expectResources {
		// Check result (it should be created/updated/unchanged, etc...)
		expectedOutput := strings.Contains(stringOutput, er)
		if !expectedOutput {
			t.Errorf("Expected substituted output to contain '%s', got %s", er, stringOutput)
		}
	}

}
