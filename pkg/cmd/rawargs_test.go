package cmd

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedResult CmdArgsRawRecognized
		expectedError  bool
	}{
		{
			name:           "No arguments",
			args:           []string{},
			expectedResult: CmdArgsRawRecognized{},
			expectedError:  false,
		},
		{
			name:           "Single filename with =",
			args:           []string{"--filename=file1.yaml"},
			expectedResult: CmdArgsRawRecognized{Filenames: []string{"file1.yaml"}},
			expectedError:  false,
		},
		{
			name:           "Recursive flag with long option",
			args:           []string{"--recursive"},
			expectedResult: CmdArgsRawRecognized{Recursive: true},
			expectedError:  false,
		},
		{
			name:           "Recursive flag with short option",
			args:           []string{"-R"},
			expectedResult: CmdArgsRawRecognized{Recursive: true},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars",
			args:           []string{"--envsubst-allowed-vars=HOME,USER"},
			expectedResult: CmdArgsRawRecognized{EnvsubstAllowedVars: []string{"HOME", "USER"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars (no =)",
			args:           []string{"--envsubst-allowed-vars", "HOME,USER"},
			expectedResult: CmdArgsRawRecognized{EnvsubstAllowedVars: []string{"HOME", "USER"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars, with append",
			args:           []string{"--envsubst-allowed-vars=HOME,USER", "--envsubst-allowed-vars=PWD"},
			expectedResult: CmdArgsRawRecognized{EnvsubstAllowedVars: []string{"HOME", "USER", "PWD"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars, with append (no =)",
			args:           []string{"--envsubst-allowed-vars", "HOME,USER", "--envsubst-allowed-vars", "PWD"},
			expectedResult: CmdArgsRawRecognized{EnvsubstAllowedVars: []string{"HOME", "USER", "PWD"}},
			expectedError:  false,
		},
		{
			name:           "Missing value for --filename",
			args:           []string{"--filename"},
			expectedResult: CmdArgsRawRecognized{},
			expectedError:  true,
		},
		{
			name:           "Unknown flag",
			args:           []string{"--unknown-flag"},
			expectedResult: CmdArgsRawRecognized{Others: []string{"--unknown-flag"}},
			expectedError:  false,
		},
		{
			name: "Mix of valid and invalid args",
			args: []string{"--filename=file.yaml", "--unknown"},
			expectedResult: CmdArgsRawRecognized{
				Filenames: []string{"file.yaml"},

				Others: []string{"--unknown"},
			},
			expectedError: false,
		},
		{
			name:           "Multiple filenames",
			args:           []string{"--filename=file1.yaml", "--filename=file2.yaml"},
			expectedResult: CmdArgsRawRecognized{Filenames: []string{"file1.yaml", "file2.yaml"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes",
			args:           []string{"--envsubst-allowed-prefixes=CI_,APP"},
			expectedResult: CmdArgsRawRecognized{EnvsubstAllowedPrefix: []string{"CI_", "APP"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes (no =)",
			args:           []string{"--envsubst-allowed-prefixes", "CI_,APP"},
			expectedResult: CmdArgsRawRecognized{EnvsubstAllowedPrefix: []string{"CI_", "APP"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes, with append",
			args:           []string{"--envsubst-allowed-prefixes=CI_,APP", "--envsubst-allowed-prefixes=TF_VAR_"},
			expectedResult: CmdArgsRawRecognized{EnvsubstAllowedPrefix: []string{"CI_", "APP", "TF_VAR_"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes with append (no =)",
			args:           []string{"--envsubst-allowed-prefixes", "CI_,APP", "--envsubst-allowed-prefixes", "TF_VAR_"},
			expectedResult: CmdArgsRawRecognized{EnvsubstAllowedPrefix: []string{"CI_", "APP", "TF_VAR_"}},
			expectedError:  false,
		},
		{
			name:           "Empty value for --filename",
			args:           []string{"--filename="},
			expectedResult: CmdArgsRawRecognized{},
			expectedError:  true,
		},
		{
			name:           "Empty value for --envsubst-allowed-vars",
			args:           []string{"--envsubst-allowed-vars="},
			expectedResult: CmdArgsRawRecognized{},
			expectedError:  true,
		},
		{
			name:           "Single filename with short flag",
			args:           []string{"-f=file3.yaml"},
			expectedResult: CmdArgsRawRecognized{Filenames: []string{"file3.yaml"}},
			expectedError:  false,
		},
		{
			name:           "Unrecognized argument without prefix",
			args:           []string{"random-arg"},
			expectedResult: CmdArgsRawRecognized{Others: []string{"random-arg"}},
			expectedError:  false,
		},
		{
			name:           "Multiple unrecognized arguments",
			args:           []string{"random-arg1", "random-arg2"},
			expectedResult: CmdArgsRawRecognized{Others: []string{"random-arg1", "random-arg2"}},
			expectedError:  false,
		},
		{
			name: "Mixed valid and unrecognized arguments",
			args: []string{"random-arg", "--filename=file.yaml"},
			expectedResult: CmdArgsRawRecognized{
				Filenames: []string{"file.yaml"},
				Others:    []string{"random-arg"},
			},
			expectedError: false,
		},
		{
			name:           "Unrecognized argument resembling a flag",
			args:           []string{"-notarealflag"},
			expectedResult: CmdArgsRawRecognized{Others: []string{"-notarealflag"}},
			expectedError:  false,
		},
		{
			name: "Unrecognized argument with spaces",
			args: []string{"random-arg", "--filename=file.yaml", "another-random-arg"},
			expectedResult: CmdArgsRawRecognized{
				Filenames: []string{"file.yaml"},
				Others:    []string{"random-arg", "another-random-arg"},
			},
			expectedError: false,
		},
		{
			name: "Valid args with unrecognized argument that looks like a short flag",
			args: []string{"-R", "-xyz"},
			expectedResult: CmdArgsRawRecognized{
				Recursive: true,
				Others:    []string{"-xyz"},
			},
			expectedError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate os.CmdArgsRawRecognized
			osArgs := append([]string{"program"}, tc.args...)
			os.Args = osArgs

			result, err := parseArgs()

			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf("expected result: %+v, got: %+v", tc.expectedResult, result)
			}
		})
	}
}

func TestParseArgs_EnvFallback(t *testing.T) {
	os.Setenv("ENVSUBST_ALLOWED_VARS", "HOME,USER")
	os.Setenv("ENVSUBST_ALLOWED_PREFIXES", "APP_,CI_")
	defer os.Unsetenv("ENVSUBST_ALLOWED_VARS")
	defer os.Unsetenv("ENVSUBST_ALLOWED_PREFIXES")

	os.Args = []string{"cmd"}
	result, err := parseArgs()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedVars := []string{"HOME", "USER"}
	expectedPrefixes := []string{"APP_", "CI_"}
	if len(result.EnvsubstAllowedVars) != len(expectedVars) {
		t.Errorf("Expected allowed vars %v, got %v", expectedVars, result.EnvsubstAllowedVars)
	}
	for i, varName := range expectedVars {
		if result.EnvsubstAllowedVars[i] != varName {
			t.Errorf("Expected allowed var %s, got %s", varName, result.EnvsubstAllowedVars[i])
		}
	}
	if len(result.EnvsubstAllowedPrefix) != len(expectedPrefixes) {
		t.Errorf("Expected allowed prefixes %v, got %v", expectedPrefixes, result.EnvsubstAllowedPrefix)
	}
	for i, prefix := range expectedPrefixes {
		if result.EnvsubstAllowedPrefix[i] != prefix {
			t.Errorf("Expected allowed prefix %s, got %s", prefix, result.EnvsubstAllowedPrefix[i])
		}
	}
}

func TestParseArgs_CmdAndEnvFlags(t *testing.T) {
	// Set environment variables
	os.Setenv("ENVSUBST_ALLOWED_VARS", "HOME,USER")
	os.Setenv("ENVSUBST_ALLOWED_PREFIXES", "APP_,CI_")
	defer os.Unsetenv("ENVSUBST_ALLOWED_VARS")
	defer os.Unsetenv("ENVSUBST_ALLOWED_PREFIXES")

	// Set command-line arguments
	os.Args = []string{
		"cmd",
		"--envsubst-allowed-vars=CMD_VAR1,CMD_VAR2",
		"--envsubst-allowed-prefixes=CMD_",
	}

	result, err := parseArgs()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify that command-line flags take precedence over environment variables
	expectedVars := []string{"CMD_VAR1", "CMD_VAR2"}
	expectedPrefixes := []string{"CMD_"}

	if len(result.EnvsubstAllowedVars) != len(expectedVars) {
		t.Errorf("Expected allowed vars %v, got %v", expectedVars, result.EnvsubstAllowedVars)
	}
	for i, varName := range expectedVars {
		if result.EnvsubstAllowedVars[i] != varName {
			t.Errorf("Expected allowed var %s, got %s", varName, result.EnvsubstAllowedVars[i])
		}
	}

	if len(result.EnvsubstAllowedPrefix) != len(expectedPrefixes) {
		t.Errorf("Expected allowed prefixes %v, got %v", expectedPrefixes, result.EnvsubstAllowedPrefix)
	}
	for i, prefix := range expectedPrefixes {
		if result.EnvsubstAllowedPrefix[i] != prefix {
			t.Errorf("Expected allowed prefix %s, got %s", prefix, result.EnvsubstAllowedPrefix[i])
		}
	}
}

func TestParseArgs_EmptyEnvVars(t *testing.T) {
	// Set empty environment variables
	os.Setenv("ENVSUBST_ALLOWED_VARS", "")
	os.Setenv("ENVSUBST_ALLOWED_PREFIXES", "")
	defer os.Unsetenv("ENVSUBST_ALLOWED_VARS")
	defer os.Unsetenv("ENVSUBST_ALLOWED_PREFIXES")

	os.Args = []string{"cmd"}
	_, err := parseArgs()

	// Expect an error due to empty environment variables
	if err == nil {
		t.Fatal("Expected an error for empty environment variables, but got none")
	}

	expectedErrorVars := "missing value for env: ENVSUBST_ALLOWED_VARS"
	expectedErrorPrefixes := "missing value for env: ENVSUBST_ALLOWED_PREFIXES"

	if !strings.Contains(err.Error(), expectedErrorVars) && !strings.Contains(err.Error(), expectedErrorPrefixes) {
		t.Errorf("Expected error to mention missing values for ENVSUBST_ALLOWED_VARS or ENVSUBST_ALLOWED_PREFIXES, got '%s'", err.Error())
	}
}
