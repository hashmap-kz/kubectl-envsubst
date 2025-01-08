package cmd

import (
	"fmt"
	"os"
	"strings"
)

const (
	envsubstAllowedVarsEnv     = "ENVSUBST_ALLOWED_VARS"
	envsubstAllowedPrefixesEnv = "ENVSUBST_ALLOWED_PREFIXES"
)

type ArgsRawRecognized struct {
	Filenames             []string
	EnvsubstAllowedVars   []string
	EnvsubstAllowedPrefix []string
	Recursive             bool
	Help                  bool
	Others                []string
	HasStdin              bool
	Version               bool
}

func allEmpty(values []string) bool {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}

func ParseArgs() (ArgsRawRecognized, error) {
	args := os.Args[1:] // Skip the program name
	var result ArgsRawRecognized

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		// Handle --filename= or -f=
		case strings.HasPrefix(arg, "--filename="), strings.HasPrefix(arg, "-f="):
			if err := handleFilename(strings.SplitN(arg, "=", 2)[1], &result); err != nil {
				return result, err
			}

		// Handle --filename or -f with a separate value
		case arg == "--filename" || arg == "-f":
			if i+1 >= len(args) || args[i+1] == "" {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}
			if err := handleFilename(args[i+1], &result); err != nil {
				return result, err
			}
			i++ // Skip the next argument

		// Handle --envsubst-allowed-vars=
		case strings.HasPrefix(arg, "--envsubst-allowed-vars="):
			list, err := appendList(strings.TrimPrefix(arg, "--envsubst-allowed-vars="))
			if err != nil {
				return result, err
			}
			result.EnvsubstAllowedVars = append(result.EnvsubstAllowedVars, list...)

		// Handle --envsubst-allowed-vars with a separate value
		case arg == "--envsubst-allowed-vars":
			if i+1 >= len(args) || args[i+1] == "" {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}
			list, err := appendList(args[i+1])
			if err != nil {
				return result, err
			}
			result.EnvsubstAllowedVars = append(result.EnvsubstAllowedVars, list...)
			i++ // Skip the next argument

		// Handle --envsubst-allowed-prefixes=
		case strings.HasPrefix(arg, "--envsubst-allowed-prefixes="):
			list, err := appendList(strings.TrimPrefix(arg, "--envsubst-allowed-prefixes="))
			if err != nil {
				return result, err
			}
			result.EnvsubstAllowedPrefix = append(result.EnvsubstAllowedPrefix, list...)

		// Handle --envsubst-allowed-prefixes with a separate value
		case arg == "--envsubst-allowed-prefixes":
			if i+1 >= len(args) || args[i+1] == "" {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}
			list, err := appendList(args[i+1])
			if err != nil {
				return result, err
			}
			result.EnvsubstAllowedPrefix = append(result.EnvsubstAllowedPrefix, list...)
			i++ // Skip the next argument

		// Handle boolean flags

		case arg == "--recursive" || arg == "-R":
			result.Recursive = true

		case arg == "--help" || arg == "-h":
			result.Help = true

		case arg == "--version":
			result.Version = true

		// Handle unrecognized arguments
		default:
			result.Others = append(result.Others, arg)
		}
	}

	// Load allowed vars and prefixes from environment variables
	if len(result.EnvsubstAllowedVars) == 0 {
		if err := loadEnvVars(envsubstAllowedVarsEnv, &result.EnvsubstAllowedVars); err != nil {
			return result, err
		}
	}
	if len(result.EnvsubstAllowedPrefix) == 0 {
		if err := loadEnvVars(envsubstAllowedPrefixesEnv, &result.EnvsubstAllowedPrefix); err != nil {
			return result, err
		}
	}

	return result, nil
}

func handleFilename(filename string, result *ArgsRawRecognized) error {
	if filename == "" {
		return fmt.Errorf("missing filename value")
	}
	if filename == "-" {
		if result.HasStdin {
			return fmt.Errorf("multiple redirection to stdin detected")
		}
		result.HasStdin = true
	} else {
		result.Filenames = append(result.Filenames, filename)
	}
	return nil
}

func appendList(value string) ([]string, error) {
	split := strings.Split(value, ",")
	if value == "" || allEmpty(split) {
		return nil, fmt.Errorf("empty list value")
	}
	return split, nil
}

func loadEnvVars(envKey string, target *[]string) error {
	value, exists := os.LookupEnv(envKey)
	if !exists {
		return nil
	}

	split := strings.Split(value, ",")
	if allEmpty(split) {
		return fmt.Errorf("missing value for env: %s", envKey)
	}

	*target = split
	return nil
}
