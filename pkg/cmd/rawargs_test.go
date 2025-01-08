package cmd

import (
	"os"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedResult ArgsRawRecognized
		expectedError  bool
	}{
		{
			name:           "No arguments",
			args:           []string{},
			expectedResult: ArgsRawRecognized{},
			expectedError:  false,
		},
		{
			name:           "Single filename with =",
			args:           []string{"--filename=file1.yaml"},
			expectedResult: ArgsRawRecognized{Filenames: []string{"file1.yaml"}},
			expectedError:  false,
		},
		{
			name:           "Recursive flag with long option",
			args:           []string{"--recursive"},
			expectedResult: ArgsRawRecognized{Recursive: true},
			expectedError:  false,
		},
		{
			name:           "Recursive flag with short option",
			args:           []string{"-R"},
			expectedResult: ArgsRawRecognized{Recursive: true},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars",
			args:           []string{"--envsubst-allowed-vars=HOME,USER"},
			expectedResult: ArgsRawRecognized{EnvsubstAllowedVars: []string{"HOME", "USER"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars (no =)",
			args:           []string{"--envsubst-allowed-vars", "HOME,USER"},
			expectedResult: ArgsRawRecognized{EnvsubstAllowedVars: []string{"HOME", "USER"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars, with append",
			args:           []string{"--envsubst-allowed-vars=HOME,USER", "--envsubst-allowed-vars=PWD"},
			expectedResult: ArgsRawRecognized{EnvsubstAllowedVars: []string{"HOME", "USER", "PWD"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars, with append (no =)",
			args:           []string{"--envsubst-allowed-vars", "HOME,USER", "--envsubst-allowed-vars", "PWD"},
			expectedResult: ArgsRawRecognized{EnvsubstAllowedVars: []string{"HOME", "USER", "PWD"}},
			expectedError:  false,
		},
		{
			name:           "Missing value for --filename",
			args:           []string{"--filename"},
			expectedResult: ArgsRawRecognized{},
			expectedError:  true,
		},
		{
			name:           "Unknown flag",
			args:           []string{"--unknown-flag"},
			expectedResult: ArgsRawRecognized{Others: []string{"--unknown-flag"}},
			expectedError:  false,
		},
		{
			name: "Mix of valid and invalid args",
			args: []string{"--filename=file.yaml", "--unknown"},
			expectedResult: ArgsRawRecognized{
				Filenames: []string{"file.yaml"},

				Others: []string{"--unknown"},
			},
			expectedError: false,
		},
		{
			name:           "Multiple filenames",
			args:           []string{"--filename=file1.yaml", "--filename=file2.yaml"},
			expectedResult: ArgsRawRecognized{Filenames: []string{"file1.yaml", "file2.yaml"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes",
			args:           []string{"--envsubst-allowed-prefixes=CI_,APP"},
			expectedResult: ArgsRawRecognized{EnvsubstAllowedPrefix: []string{"CI_", "APP"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes (no =)",
			args:           []string{"--envsubst-allowed-prefixes", "CI_,APP"},
			expectedResult: ArgsRawRecognized{EnvsubstAllowedPrefix: []string{"CI_", "APP"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes, with append",
			args:           []string{"--envsubst-allowed-prefixes=CI_,APP", "--envsubst-allowed-prefixes=TF_VAR_"},
			expectedResult: ArgsRawRecognized{EnvsubstAllowedPrefix: []string{"CI_", "APP", "TF_VAR_"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes with append (no =)",
			args:           []string{"--envsubst-allowed-prefixes", "CI_,APP", "--envsubst-allowed-prefixes", "TF_VAR_"},
			expectedResult: ArgsRawRecognized{EnvsubstAllowedPrefix: []string{"CI_", "APP", "TF_VAR_"}},
			expectedError:  false,
		},
		{
			name:           "Empty value for --filename",
			args:           []string{"--filename="},
			expectedResult: ArgsRawRecognized{},
			expectedError:  true,
		},
		{
			name:           "Empty value for --envsubst-allowed-vars",
			args:           []string{"--envsubst-allowed-vars="},
			expectedResult: ArgsRawRecognized{},
			expectedError:  true,
		},
		{
			name:           "Single filename with short flag",
			args:           []string{"-f=file3.yaml"},
			expectedResult: ArgsRawRecognized{Filenames: []string{"file3.yaml"}},
			expectedError:  false,
		},
		{
			name:           "Unrecognized argument without prefix",
			args:           []string{"random-arg"},
			expectedResult: ArgsRawRecognized{Others: []string{"random-arg"}},
			expectedError:  false,
		},
		{
			name:           "Multiple unrecognized arguments",
			args:           []string{"random-arg1", "random-arg2"},
			expectedResult: ArgsRawRecognized{Others: []string{"random-arg1", "random-arg2"}},
			expectedError:  false,
		},
		{
			name: "Mixed valid and unrecognized arguments",
			args: []string{"random-arg", "--filename=file.yaml"},
			expectedResult: ArgsRawRecognized{
				Filenames: []string{"file.yaml"},
				Others:    []string{"random-arg"},
			},
			expectedError: false,
		},
		{
			name:           "Unrecognized argument resembling a flag",
			args:           []string{"-notarealflag"},
			expectedResult: ArgsRawRecognized{Others: []string{"-notarealflag"}},
			expectedError:  false,
		},
		{
			name: "Unrecognized argument with spaces",
			args: []string{"random-arg", "--filename=file.yaml", "another-random-arg"},
			expectedResult: ArgsRawRecognized{
				Filenames: []string{"file.yaml"},
				Others:    []string{"random-arg", "another-random-arg"},
			},
			expectedError: false,
		},
		{
			name: "Valid args with unrecognized argument that looks like a short flag (1)",
			args: []string{"-R", "-xyz"},
			expectedResult: ArgsRawRecognized{
				Recursive: true,
				Others:    []string{"-xyz"},
			},
			expectedError: false,
		},
		{
			name: "Valid args with unrecognized argument that looks like a short flag (2)",
			args: []string{"-h", "-xyz"},
			expectedResult: ArgsRawRecognized{
				Help:   true,
				Others: []string{"-xyz"},
			},
			expectedError: false,
		},
		{
			name: "Valid args with unrecognized argument that looks like a short flag (3)",
			args: []string{"-h", "-xyz", "--version"},
			expectedResult: ArgsRawRecognized{
				Help:    true,
				Version: true,
				Others:  []string{"-xyz"},
			},
			expectedError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate os.CmdArgsRawRecognized
			osArgs := append([]string{"program"}, tc.args...)
			os.Args = osArgs

			result, err := ParseArgs()

			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf("expected result: %+v, got: %+v", tc.expectedResult, result)
			}
		})
	}
}

func TestParseArgs_SingleStdin(t *testing.T) {
	os.Args = []string{"cmd", "--filename", "-"}
	result, err := ParseArgs()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.HasStdin {
		t.Errorf("Expected HasStdin to be true")
	}

	if len(result.Filenames) != 0 {
		t.Errorf("Expected Filenames to be empty, got: %v", result.Filenames)
	}
}

func TestParseArgs_MultipleStdin(t *testing.T) {
	os.Args = []string{"cmd", "--filename", "-", "-f", "-"}
	_, err := ParseArgs()

	if err == nil {
		t.Fatal("Expected an error for multiple stdin redirection, but got none")
	}

	expectedError := "multiple redirection to stdin detected"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadEnvVars(t *testing.T) {
	tests := []struct {
		name      string
		envKey    string
		envValue  string
		initial   []string
		expected  []string
		expectErr bool
	}{
		{
			name:      "Valid Environment Variable",
			envKey:    "TEST_ENV",
			envValue:  "VAR1,VAR2,VAR3",
			initial:   nil,
			expected:  []string{"VAR1", "VAR2", "VAR3"},
			expectErr: false,
		},
		{
			name:      "Missing Environment Variable",
			envKey:    "MISSING_ENV",
			envValue:  "",
			initial:   nil,
			expected:  nil,
			expectErr: false,
		},
		// not an error, it means that configuration ENVS were not set
		// an error will be handled in subst stage
		{
			name:      "Empty Environment Variable Value",
			envKey:    "EMPTY_ENV",
			envValue:  "",
			initial:   nil,
			expected:  nil,
			expectErr: false,
		},
		{
			name:      "Whitespace Only Environment Variable",
			envKey:    "WHITESPACE_ENV",
			envValue:  "   ,   ,   ",
			initial:   nil,
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Pre-existing Values in Target",
			envKey:    "TEST_ENV",
			envValue:  "NEW1,NEW2",
			initial:   []string{"OLD1", "OLD2"},
			expected:  []string{"NEW1", "NEW2"},
			expectErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup environment variable
			if test.envValue != "" {
				os.Setenv(test.envKey, test.envValue)
				defer os.Unsetenv(test.envKey)
			}

			// Initialize target
			target := test.initial

			// Call function
			err := loadEnvVars(test.envKey, &target)

			// Check for errors
			if test.expectErr {
				if err == nil {
					t.Errorf("Test '%s': expected error but got none", test.name)
				}
			} else {
				if err != nil {
					t.Errorf("Test '%s': unexpected error: %v", test.name, err)
				}
			}

			// Check result
			if !reflect.DeepEqual(target, test.expected) {
				t.Errorf("Test '%s': result = %v; want %v", test.name, target, test.expected)
			}
		})
	}
}

func TestParseArgs_ErrorsAndReturns(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		envVars   map[string]string
		expectErr string
		validate  func(t *testing.T, result ArgsRawRecognized)
	}{
		{
			name:      "Missing value for --filename=",
			args:      []string{"app", "--filename="},
			expectErr: "missing filename value",
		},
		{
			name:      "Missing value for --filename",
			args:      []string{"app", "--filename"},
			expectErr: "missing value for flag --filename",
		},
		{
			name:      "Multiple redirection to stdin",
			args:      []string{"app", "--filename=-", "--filename=-"},
			expectErr: "multiple redirection to stdin detected",
		},
		{
			name:      "Empty list value for --envsubst-allowed-vars=",
			args:      []string{"app", "--envsubst-allowed-vars="},
			expectErr: "empty list value",
		},
		{
			name:      "Empty list value for --envsubst-allowed-vars",
			args:      []string{"app", "--envsubst-allowed-vars"},
			expectErr: "missing value for flag --envsubst-allowed-vars",
		},
		{
			name:      "Empty list value for --envsubst-allowed-prefixes=",
			args:      []string{"app", "--envsubst-allowed-prefixes="},
			expectErr: "empty list value",
		},
		{
			name:      "Empty list value for --envsubst-allowed-prefixes",
			args:      []string{"app", "--envsubst-allowed-prefixes"},
			expectErr: "missing value for flag --envsubst-allowed-prefixes",
		},
		{
			name: "Empty environment variable for ENVSUBST_ALLOWED_VARS",
			args: []string{"app"},
			envVars: map[string]string{
				"ENVSUBST_ALLOWED_VARS": "",
			},
			expectErr: "missing value for env: ENVSUBST_ALLOWED_VARS",
		},
		{
			name: "Empty environment variable for ENVSUBST_ALLOWED_PREFIXES",
			args: []string{"app"},
			envVars: map[string]string{
				"ENVSUBST_ALLOWED_PREFIXES": "",
			},
			expectErr: "missing value for env: ENVSUBST_ALLOWED_PREFIXES",
		},
		{
			name: "Successful parsing with all flags",
			args: []string{"app", "--filename=test.yaml", "--envsubst-allowed-vars=VAR1,VAR2", "--envsubst-allowed-prefixes=PREFIX1,PREFIX2", "--recursive", "--help"},
			validate: func(t *testing.T, result ArgsRawRecognized) {
				if len(result.Filenames) != 1 || result.Filenames[0] != "test.yaml" {
					t.Errorf("Expected filename 'test.yaml', got %v", result.Filenames)
				}
				if len(result.EnvsubstAllowedVars) != 2 || result.EnvsubstAllowedVars[0] != "VAR1" || result.EnvsubstAllowedVars[1] != "VAR2" {
					t.Errorf("Expected EnvsubstAllowedVars to contain [VAR1, VAR2], got %v", result.EnvsubstAllowedVars)
				}
				if len(result.EnvsubstAllowedPrefix) != 2 || result.EnvsubstAllowedPrefix[0] != "PREFIX1" || result.EnvsubstAllowedPrefix[1] != "PREFIX2" {
					t.Errorf("Expected EnvsubstAllowedPrefix to contain [PREFIX1, PREFIX2], got %v", result.EnvsubstAllowedPrefix)
				}
				if !result.Recursive {
					t.Errorf("Expected Recursive to be true")
				}
				if !result.Help {
					t.Errorf("Expected Help to be true")
				}
			},
		},
	}

	originalEnv := make(map[string]string)
	for _, k := range os.Environ() {
		originalEnv[k] = os.Getenv(k)
	}
	defer func() {
		for k, v := range originalEnv {
			os.Setenv(k, v)
		}
	}()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set environment variables for the test
			for k, v := range test.envVars {
				os.Setenv(k, v)
			}
			defer func(envVars map[string]string) {
				for k := range envVars {
					os.Unsetenv(k)
				}
			}(test.envVars)

			// Simulate command-line arguments
			os.Args = test.args

			// Run ParseArgs
			result, err := ParseArgs()
			if test.expectErr != "" {
				if err == nil || err.Error() != test.expectErr {
					t.Fatalf("Expected error '%s', got '%v'", test.expectErr, err)
				}
			} else if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Validate results
			if test.validate != nil {
				test.validate(t, result)
			}
		})
	}
}
