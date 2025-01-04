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

// basic

func TestEnvsubstIntegration_SubstApplyFromFile(t *testing.T) {

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
		t.Errorf("Expected resource %s to exist, got %s", resourceName, string(validateOutput))
	}

}

func TestEnvsubstIntegration_SubstApplyFromStdin(t *testing.T) {

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
		t.Errorf("Expected resource %s to exist, got %s", resourceName, string(validateOutput))
	}

}

func TestEnvsubstIntegration_SubstApplyFromUrl(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("Integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	const url = "https://raw.githubusercontent.com/hashmap-kz/kubectl-envsubst/refs/heads/master/integration/immutable_data/01_deployment.yaml"
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
		t.Errorf("Expected resource %s to exist, got %s", resourceName, string(validateOutput))
	}

}

// mixed substitution

func TestEnvsubstIntegration_SubstApplyConfigmapMixed(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
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

func TestEnvsubstIntegration_SubstMixedManifestsCombined(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("Integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	namespaceName := randomIdent(32)
	createNs(t, namespaceName)
	defer cleanupResource(t, "ns", namespaceName)

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
	defer setContextNs(t, "default")

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/resolve/subst-combined")
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

func TestEnvsubstIntegration_SubstMixedManifestsSeparated(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("Integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	namespaceName := randomIdent(32)
	createNs(t, namespaceName)
	defer cleanupResource(t, "ns", namespaceName)

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
	defer setContextNs(t, "default")

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/resolve/subst-separated")
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

func TestEnvsubstIntegration_Subst_MixedManifests_MixedExtensions_MixedFileFormats(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("Integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	namespaceName := randomIdent(32)
	createNs(t, namespaceName)
	defer cleanupResource(t, "ns", namespaceName)

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
	defer setContextNs(t, "default")

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/resolve/subst-yaml-json")
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

func TestEnvsubstIntegration_NoSubst_MixedManifests(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("Integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	namespaceName := "kubectl-envsubst-integration-tests-ns-1"
	defer cleanupResource(t, "ns", namespaceName)

	// Setup context
	setContextNs(t, namespaceName)
	defer setContextNs(t, "default")

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

func TestEnvsubstIntegration_NoSubst_MixedManifests_MixedFileFormats(t *testing.T) {

	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("Integration test was skipped due to configuration")
		return
	}

	t.Log("running integration test: ", t.Name())
	printEnvsubstVersionInfo(t)

	namespaceName := "kubectl-envsubst-integration-tests-ns-2"
	defer cleanupResource(t, "ns", namespaceName)

	// Setup context
	setContextNs(t, namespaceName)
	defer setContextNs(t, "default")

	// Run kubectl-envsubst

	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "immutable_data/resolve/plain-yaml-json")
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

func createNs(t *testing.T, ns string) {
	cmd := exec.Command("kubectl", "create", "ns", ns)
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
	cmd := exec.Command("kubectl", "krew", "info", "envsubst")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "VERSION") {
			t.Logf("*** kubectl-envsubst %s ***", line)
			t.Logf("*** kubectl-envsubst path: %s ***", getEnvsubstPath(t))
		}
	}
}
