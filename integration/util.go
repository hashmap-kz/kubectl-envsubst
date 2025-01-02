package integration

import (
	"bufio"
	"math/rand"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const (
	integrationTestEnv  = "KUBECTL_ENVSUBST_INTEGRATION_TESTS_AVAILABLE"
	integrationTestFlag = "0xcafebabe"
)

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

func printInfo(t *testing.T) {
	cmd := exec.Command("kubectl", "krew", "info", "envsubst")
	output, _ := cmd.CombinedOutput()

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "VERSION") {
			t.Logf("*** kubectl-envsubst %s ***", line)
		}
	}
}
