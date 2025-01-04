package integration

import (
	"os/exec"
	"strings"
	"testing"
)

func TestEnvsubstIntegration_NoSubst_MixedManifests(t *testing.T) {

	// if os.Getenv(integrationTestEnv) != integrationTestFlag {
	// 	t.Log("Integration test was skipped due to configuration")
	// 	return
	// }

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	namespaceName := "kubectl-envsubst-integration-tests-ns-1"
	cleanupResource(t, "ns", namespaceName)
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
		"serviceaccount/my-app created",
		"role.rbac.authorization.k8s.io/my-app created",
		"rolebinding.rbac.authorization.k8s.io/my-app created",
		"configmap/my-app created",
		"secret/my-app created",
		"deployment.apps/my-app created",
		"service/my-app created",
	}

	for _, er := range expectResources {
		// Check result (it should be created/updated/unchanged, etc...)
		expectedOutput := strings.Contains(stringOutput, er)
		if !expectedOutput {
			t.Errorf("Expected substituted output to contain '%s', got %s", er, stringOutput)
		}
	}

}
