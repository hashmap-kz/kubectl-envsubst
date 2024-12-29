package app

import (
	"os"
	"regexp"
	"strings"
)

// SubstituteEnvs replaces only the specified environment variables in the given text.
func SubstituteEnvs(text string, allowedEnvs []string) string {
	envMap := make(map[string]string)
	for _, env := range allowedEnvs {
		if value, exists := os.LookupEnv(env); exists {
			envMap[env] = value
		}
	}

	// Match placeholders like ${VAR} or $VAR
	re := regexp.MustCompile(`\$\{?([a-zA-Z_][a-zA-Z0-9_]*)\}?`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the variable name
		varName := strings.Trim(match, "${}")
		if value, ok := envMap[varName]; ok {
			return value
		}
		return match // Leave the original text if not in allowedEnvs
	})
}
