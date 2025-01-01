package cmd

import (
	"os"
	"regexp"
	"strings"
)

import (
	"fmt"
)

var (
	// Match placeholders like ${VAR} or $VAR
	envVarRegex = regexp.MustCompile(`\$\{?([a-zA-Z_][a-zA-Z0-9_]*)\}?`)
)

type Envsubst struct {
	allowedVars     []string
	allowedPrefixes []string
	strict          bool
}

func NewEnvsubst(allowedVars []string, allowedPrefixes []string, strict bool) *Envsubst {
	return &Envsubst{
		allowedVars:     allowedVars,
		allowedPrefixes: allowedPrefixes,
		strict:          strict,
	}
}

func (p *Envsubst) SubstituteEnvs(text string) (string, error) {

	// Collect allowed environment variables and prefixes
	envMap := make(map[string]string)
	for _, env := range p.allowedVars {
		if value, exists := os.LookupEnv(env); exists {
			envMap[env] = value
		}
	}

	// Add variables with allowed prefixes
	for _, prefix := range p.allowedPrefixes {
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 && strings.HasPrefix(parts[0], prefix) {
				envMap[parts[0]] = parts[1]
			}
		}
	}

	// Perform substitution using regex
	substituted := envVarRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the variable name
		// alternate: varName := envVarRegex.FindStringSubmatch(match)[1]
		varName := strings.Trim(match, "${}")

		// get value, according to filters
		if value, ok := envMap[varName]; ok {
			return value
		}

		return match
	})

	// Handle strict mode by detecting unresolved variables
	if p.strict {
		unresolved := envVarRegex.FindAllString(substituted, -1)
		if len(unresolved) > 0 {
			return "", fmt.Errorf("undefined variables: %v", unresolved)
		}
	}

	return substituted, nil
}
