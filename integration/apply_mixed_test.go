package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestEnvsubstIntegration_SubstApplyConfigmapMixed(t *testing.T) {
	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	printEnvsubstVersionInfo(t)

	resourceName := randomIdent(32)
	defer cleanupResource(t, "configmap", resourceName)

	// Setup environment variables that was used in substitution
	os.Setenv("CI_PROJECT_ROOT_NAMESPACE", "trade-system-application")
	os.Setenv("CI_PROJECT_NAME", resourceName)
	os.Setenv("CI_COMMIT_REF_NAME", "dev")
	os.Setenv("INFRA_DOMAIN_NAME", "company.org")
	defer os.Unsetenv("CI_PROJECT_ROOT_NAMESPACE")
	defer os.Unsetenv("CI_PROJECT_NAME")
	defer os.Unsetenv("CI_COMMIT_REF_NAME")
	defer os.Unsetenv("INFRA_DOMAIN_NAME")

	// configure CLI
	os.Setenv("ENVSUBST_ALLOWED_PREFIXES", "CI_,IMAGE_,INFRA_")
	defer os.Unsetenv("ENVSUBST_ALLOWED_PREFIXES")

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/02_configmap.yaml")
	output, err := cmdEnvsubstApply.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, string(output))
	}
	t.Log(string(output))

	// Check result (it should be created/updated/unchanged, etc...)
	expectedOutput := strings.Contains(string(output), fmt.Sprintf("configmap/%s", resourceName))
	if !expectedOutput {
		t.Errorf("Expected substituted output to contain 'configmap/%s', got %s", resourceName, string(output))
	}

	// Validate applied resource
	validateCmd := exec.Command("kubectl", "get", "cm", resourceName)
	validateOutput, err := validateCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to validate applied resource: %v, output: %s", err, string(validateOutput))
	}
	if !strings.Contains(string(validateOutput), resourceName) {
		t.Errorf("Expected resource %s to exist, got %s", resourceName, string(validateOutput))
	}

	// check data
	getCmData := exec.Command("kubectl", "get", "cm", resourceName, "-ojsonpath={.data.trade-system-application-dev}")
	cmDataOut, err := getCmData.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to validate applied resource: %v, output: %s", err, string(validateOutput))
	}

	// expecting that config was not substituted
	expectCmOut := strings.TrimSpace(`
server {
    listen 80;
    server_name localhost;
    return 301 https://$server_name$request_uri;
    server_tokens off;
    access_log off;
    error_log off;
}
server {
    listen 443 ssl;
    server_name localhost;
    access_log /var/log/nginx/access.log json_combined;
    error_log /var/log/nginx/error.log warn;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-Host $host;
    proxy_set_header X-Forwarded-Server $host;
    location /api/ {
        proxy_pass http://gateway-service-http:8080/api/;
    }
}
`)
	if strings.TrimSpace(string(cmDataOut)) != strings.TrimSpace(expectCmOut) {
		t.Fatalf("configmaps data are different")
	}
}
