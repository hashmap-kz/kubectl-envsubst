package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Match placeholders like ${VAR} or $VAR
var envVarRegex = regexp.MustCompile(`\$\{?([a-zA-Z_][a-zA-Z0-9_]*)\}?`)

type Envsubst struct {
	allowedVars     []string
	allowedPrefixes []string
	strict          bool
	verbose         bool
}

func NewEnvsubst(allowedVars, allowedPrefixes []string, strict bool) *Envsubst {
	return &Envsubst{
		allowedVars:     allowedVars,
		allowedPrefixes: allowedPrefixes,
		strict:          strict,
	}
}

func (p *Envsubst) SubstituteEnvs(text string) (string, error) {
	// Collect allowed environment variables
	envMap := p.collectAllowedEnvVars()

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

	// Handle unresolved variables in strict mode
	// Returns error, if and only if an unresolved variable is from one of the filter-list.
	// Ignoring other unexpanded variables, that may be a parts of config-maps, etc...
	//
	if err := p.checkUnresolvedStrictMode(substituted); err != nil {
		return "", err
	}

	// Log unresolved variables in verbose mode
	// if there are unexpanded placeholders, it's not an error, just debug-info
	// it's not an error, because these placeholders are not in filter lists, so they remain unchanged
	unresolved := envVarRegex.FindAllString(substituted, -1)
	p.logUnresolvedVariables(unresolved)

	return substituted, nil
}

func (p *Envsubst) SetVerbose(value bool) {
	p.verbose = value
}

// Helper Functions

// collectAllowedEnvVars collects variables and prefixes allowed for substitution
func (p *Envsubst) collectAllowedEnvVars() map[string]string {
	envMap := make(map[string]string)

	// Collect variables in the allowedVars list
	for _, env := range p.allowedVars {
		if value, exists := os.LookupEnv(env); exists {
			envMap[env] = value
		}
	}

	// Collect variables matching allowed prefixes
	globalEnv := preprocessEnv()
	for _, prefix := range p.allowedPrefixes {
		for key, value := range globalEnv {
			if strings.HasPrefix(key, prefix) {
				envMap[key] = value
			}
		}
	}

	return envMap
}

// preprocessEnv preprocesses environment variables into a map
func preprocessEnv() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap
}

// checkUnresolvedStrictMode checks unresolved variables in strict mode
func (p *Envsubst) checkUnresolvedStrictMode(substituted string) error {
	unresolved := envVarRegex.FindAllString(substituted, -1)
	if p.strict {
		filtered := p.filterUnresolvedByAllowedLists(unresolved)
		if len(filtered) > 0 {
			return fmt.Errorf("undefined variables: [%s]", strings.Join(filtered, ", "))
		}
	}
	return nil
}

// logUnresolvedVariables logs unresolved variables in verbose mode
func (p *Envsubst) logUnresolvedVariables(unresolved []string) {
	if p.verbose {
		for _, variable := range p.sortUnresolved(unresolved) {
			log.Printf("DEBUG: an unresolved variable that is not in the filter list remains unchanged: %s", variable)
		}
	}
}

// filterUnresolvedByAllowedLists filters unresolved variables based on allowed lists
func (p *Envsubst) filterUnresolvedByAllowedLists(input []string) []string {
	result := []string{}
	for _, v := range input {
		v := strings.Trim(v, "${}")
		if p.isInFilter(v) && !varInSlice(v, result) {
			result = append(result, v)
		}
	}
	sort.Strings(result)
	return result
}

// isInFilter checks if a variable is in the allowed lists
func (p *Envsubst) isInFilter(e string) bool {
	for _, allowed := range p.allowedVars {
		if e == allowed {
			return true
		}
	}
	for _, prefix := range p.allowedPrefixes {
		if strings.HasPrefix(e, prefix) {
			return true
		}
	}
	return false
}

// sortUnresolved removes duplicates and sorts unresolved variables
func (p *Envsubst) sortUnresolved(input []string) []string {
	result := []string{}
	for _, v := range input {
		v := strings.Trim(v, "${}")
		if !varInSlice(v, result) {
			result = append(result, v)
		}
	}
	sort.Strings(result)
	return result
}

// varInSlice checks if a string is in a slice
func varInSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
