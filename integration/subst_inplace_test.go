package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestSubstFromFile(t *testing.T) {
	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	printEnvsubstVersionInfo(t)

	testCases := []struct {
		name                 string
		inputContent         string
		envVars              []string
		expectedOutput       string
		expectedImage        string
		expectedResourceName string
	}{
		{
			name: "Service with variable substitution",
			inputContent: `
apiVersion: v1
kind: Service
metadata:
  name: ${SERVICE_NAME}
spec:
  selector:
    app: ${APP_NAME}
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
`,
			envVars: []string{
				"SERVICE_NAME=my-service",
				"APP_NAME=my-app",
				"ENVSUBST_ALLOWED_PREFIXES=SERVICE_,APP_",
			},

			expectedOutput: "service/my-service",
		},

		{
			name: "Deployment with variable substitution",
			inputContent: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${DEPLOYMENT_NAME}
spec:
  replicas: ${REPLICAS}
  selector:
    matchLabels:
      app: ${APP_NAME}
  template:
    metadata:
      labels:
        app: ${APP_NAME}
    spec:
      containers:
        - name: ${CONTAINER_NAME}
          image: ${IMAGE_NAME}:${IMAGE_TAG}
`,
			envVars: []string{
				"DEPLOYMENT_NAME=my-deployment",
				"REPLICAS=3",
				"APP_NAME=my-app",
				"CONTAINER_NAME=my-container",
				"IMAGE_NAME=nginx",
				"IMAGE_TAG=1.27.3-bookworm",
				"ENVSUBST_ALLOWED_PREFIXES=DEP,REP,APP,CON,IM,EN,",
			},
			expectedOutput:       "deployment.apps/my-deployment",
			expectedImage:        "nginx:1.27.3-bookworm",
			expectedResourceName: "my-deployment",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// prepare tmp file with content
			tmpFile, err := createTempFile(t, tc.inputContent, "yaml")
			if err != nil {
				t.Fatal("failed to create temp-file")
			}
			defer os.Remove(tmpFile)

			// run plugin
			cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "--filename", tmpFile)
			cmdEnvsubstApply.Env = append(os.Environ(), tc.envVars...)
			output, err := cmdEnvsubstApply.CombinedOutput()
			stringOutput := string(output)
			if err != nil {
				t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, stringOutput)
			}
			t.Logf("\n%s\n", strings.TrimSpace(stringOutput))

			// check that resources were applied
			expectedOutput := strings.Contains(stringOutput, tc.expectedOutput)
			if !expectedOutput {
				t.Errorf("Expected substituted output to contain '%s', got %s", tc.expectedOutput, stringOutput)
			}

			if tc.expectedImage != "" {
				imageName := getDeploymentImageName(t, tc.expectedResourceName)
				if imageName != tc.expectedImage {
					t.Errorf("Expected image: %s, got %s", tc.expectedImage, imageName)
				}
			}
		})
	}
}
