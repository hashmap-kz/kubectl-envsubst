package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	// Match placeholders like ${VAR} or $VAR
	envVarRegex = regexp.MustCompile(`\$\{?([a-zA-Z_][a-zA-Z0-9_]*)\}?`)
)

type Envsubst struct {
	allowedVars     []string
	allowedPrefixes []string
	strict          bool
	verbose         bool
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
	// Returns error, if and only if an unresolved variable is from one of the filter-list.
	// Ignoring other unexpanded variables, that may be a parts of config-maps, etc...
	//
	if p.strict {
		unresolved := envVarRegex.FindAllString(substituted, -1)

		// if an unresolved variable is supposed to be substituted but is not, it is considered an error
		filterUnresolved := p.filterUnresolvedByAllowedLists(unresolved)
		if len(filterUnresolved) > 0 {
			sb := strings.Builder{}
			for _, k := range filterUnresolved {
				if !strings.Contains(sb.String(), k) {
					sb.WriteString(k + ", ")
				}
			}
			resultList := strings.TrimSpace(sb.String())
			return "", fmt.Errorf("undefined variables: [%s]", strings.TrimSuffix(resultList, ","))
		}

		// verbose mode here: if there are unexpanded placeholders, it's not an error, just debug-info
		// it's not an error, because these placeholders are not in filter lists, so they remain unchanged
		if p.verbose {
			sortUnresolved := p.sortUnresolved(unresolved)
			for _, u := range sortUnresolved {
				log.Printf("DEBUG: an unresolved variable that is not in the filter list remains unchanged: %s", u)
			}
		}
	}

	return substituted, nil
}

func (p *Envsubst) SetVerbose(value bool) {
	p.verbose = value
}

func (p *Envsubst) sortUnresolved(input []string) []string {
	result := []string{}
	for _, v := range input {
		v := strings.Trim(v, "${}")
		if varInSlice(v, result) {
			continue
		}
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}

func (p *Envsubst) filterUnresolvedByAllowedLists(input []string) []string {
	result := []string{}
	for _, v := range input {
		v := strings.Trim(v, "${}")
		if !p.isInFilter(v) {
			continue
		}
		if varInSlice(v, result) {
			continue
		}
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}

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

func varInSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
