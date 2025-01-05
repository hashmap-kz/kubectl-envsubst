package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestEnvsubstIntegration_SubstMixedManifestsCombined(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	printEnvsubstVersionInfo(t)

	namespaceName := randomIdent(32)
	createNs(t, namespaceName)

	// Setup environment variables that was used in substitution
	os.Setenv("IMAGE_NAME", "nginx")
	os.Setenv("IMAGE_TAG", "latest")
	os.Setenv("APP_NAMESPACE", namespaceName)
	os.Setenv("APP_NAME", "my-app")
	defer os.Unsetenv("IMAGE_NAME")
	defer os.Unsetenv("IMAGE_TAG")
	defer os.Unsetenv("APP_NAMESPACE")
	defer os.Unsetenv("APP_NAME")

	// configure CLI
	os.Setenv("ENVSUBST_ALLOWED_PREFIXES", "APP_,IMAGE_")
	defer os.Unsetenv("ENVSUBST_ALLOWED_PREFIXES")

	// Setup context
	setContextNs(t, namespaceName)

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/resolve/subst-combined")
	output, err := cmdEnvsubstApply.CombinedOutput()
	stringOutput := string(output)
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, stringOutput)
	}
	t.Logf("\n%s\n", strings.TrimSpace(stringOutput))

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
