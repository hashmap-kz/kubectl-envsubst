package integration

import (
	"bufio"
	"math/rand"
	"os/exec"
	"strings"
	"testing"
	"time"
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
