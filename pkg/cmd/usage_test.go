package cmd

import (
	"strings"
	"testing"
)

// To prevent accidentally make a typo in a usage, just duplicate it. 

func TestUsageMessage(t *testing.T) {

	expected := strings.TrimSpace(`
Expands environment variables in manifests, before applying them

Usage:
  # substitute variables whose names start with one of the prefixes
  kubectl envsubst apply -f manifests/ --envsubst-allowed-prefixes=CI_,APP_
  
  # substitute well-defined variables
  kubectl envsubst apply -f manifests/ --envsubst-allowed-vars=CI_PROJECT_NAME,CI_COMMIT_REF_NAME,APP_IMAGE
  
  # mixed mode, check both full match and prefix match 
  kubectl envsubst apply -f manifests/ --envsubst-allowed-prefixes=CI_,APP_ --envsubst-allowed-vars=HOME,USER

Examples:
  # example usage with other kubectl flags
  kubectl envsubst apply -f manifests/ --dry-run=client -oyaml --envsubst-allowed-prefixes=APP_

Flags:
  --envsubst-allowed-vars
      Accepts a comma-separated list of variable names allowed for substitution. 
      Variables not included in this list will not be substituted.

  --envsubst-allowed-prefixes
      Accepts a comma-separated list of prefixes. 
      Only variables with names starting with one of these prefixes will be substituted; others will be ignored.
`)

	actual := UsageMessage()

	if actual != expected {
		t.Errorf("UsageMessage() returned unexpected result.\nExpected:\n%s\n\nGot:\n%s", expected, actual)
	}
}
