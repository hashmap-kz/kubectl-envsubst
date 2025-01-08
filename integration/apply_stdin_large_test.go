package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var resourceTemplateForStdin = `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-%03d
data:
  app.conf: |
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
`

func TestEnvsubstIntegration_SubstApplyFromStdinWithLargeContent(t *testing.T) {
	if os.Getenv(integrationTestEnv) != integrationTestFlag {
		t.Log("integration test was skipped due to configuration")
		return
	}

	printEnvsubstVersionInfo(t)

	// Prepare input manifest
	sb := strings.Builder{}
	for i := 0; i < 100; i++ {
		sb.WriteString(fmt.Sprintf(resourceTemplateForStdin, i))
	}
	manifest := sb.String()

	// Run kubectl-envsubst
	cmdEnvsubstApply := exec.Command("kubectl", "envsubst", "apply", "-f", "-")
	cmdEnvsubstApply.Stdin = strings.NewReader(manifest)
	output, err := cmdEnvsubstApply.CombinedOutput()
	stringOutput := string(output)
	if err != nil {
		t.Fatalf("Failed to run kubectl envsubst: %v, output: %s", err, stringOutput)
	}

	expectResources := []string{
		"configmap/cm-000",
		"configmap/cm-050",
		"configmap/cm-089",
	}

	for _, er := range expectResources {
		// Check result (it should be created/updated/unchanged, etc...)
		expectedOutput := strings.Contains(stringOutput, er)
		if !expectedOutput {
			t.Errorf("Expected substituted output to contain '%s', got %s", er, stringOutput)
		}
	}
}
