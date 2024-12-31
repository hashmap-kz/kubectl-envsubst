package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

func IsRunningAsPlugin() bool {
	return strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-")
}

func UsageMessage() string {
	var usageRaw = `
Expands environment variables in manifests, before applying them

Usage:
  # substitute variables whose names start with one of the prefixes
  kubectl envsubst apply -f manifests/ --envsubst-allowed-prefixes=CI_,APP_
  
  # substitute well-defined variables
  kubectl envsubst apply -f manifests/ --envsubst-allowed-vars=CI_PROJECT_NAME,CI_COMMIT_REF_NAME,APP_IMAGE
  
  # mixed mode, check both full match and prefix match 
  kubectl envsubst apply -f manifests/ --envsubst-allowed-prefixes=CI_,APP_ --envsubst-allowed-vars=HOME,USER

Examples:
  # example with other flags
  kubectl envsubst apply -f testdata/subst/01.yaml --dry-run=client -oyaml --envsubst-allowed-prefixes=APP_

Flags:
  --envsubst-allowed-vars     flag, that consumes a list of comma-separated names that allowed for expansion
  --envsubst-allowed-prefixes flag, that consumes a list of comma-separated prefixes that allowed for expansion
`

	return strings.TrimSpace(usageRaw)
}
