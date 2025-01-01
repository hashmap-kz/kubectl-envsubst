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
		expectedResult CmdArgsRawRecognized
		expectedError  bool
	}{
		{
			name:           "No arguments",
			args:           []string{},
			expectedResult: CmdArgsRawRecognized{Strict: true},
			expectedError:  false,
		},
		{
			name:           "Single filename with =",
			args:           []string{"--filename=file1.yaml"},
			expectedResult: CmdArgsRawRecognized{Strict: true, Filenames: []string{"file1.yaml"}},
			expectedError:  false,
		},
		{
			name:           "Strict flag",
			args:           []string{"--envsubst-no-strict"},
			expectedResult: CmdArgsRawRecognized{Strict: false},
			expectedError:  false,
		},
		{
			name:           "Recursive flag with long option",
			args:           []string{"--recursive"},
			expectedResult: CmdArgsRawRecognized{Strict: true, Recursive: true},
			expectedError:  false,
		},
		{
			name:           "Recursive flag with short option",
			args:           []string{"-R"},
			expectedResult: CmdArgsRawRecognized{Strict: true, Recursive: true},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars",
			args:           []string{"--envsubst-allowed-vars=HOME,USER"},
			expectedResult: CmdArgsRawRecognized{Strict: true, EnvsubstAllowedVars: []string{"HOME", "USER"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars (no =)",
			args:           []string{"--envsubst-allowed-vars", "HOME,USER"},
			expectedResult: CmdArgsRawRecognized{Strict: true, EnvsubstAllowedVars: []string{"HOME", "USER"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars, with append",
			args:           []string{"--envsubst-allowed-vars=HOME,USER", "--envsubst-allowed-vars=PWD"},
			expectedResult: CmdArgsRawRecognized{Strict: true, EnvsubstAllowedVars: []string{"HOME", "USER", "PWD"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed vars, with append (no =)",
			args:           []string{"--envsubst-allowed-vars", "HOME,USER", "--envsubst-allowed-vars", "PWD"},
			expectedResult: CmdArgsRawRecognized{Strict: true, EnvsubstAllowedVars: []string{"HOME", "USER", "PWD"}},
			expectedError:  false,
		},
		{
			name:           "Missing value for --filename",
			args:           []string{"--filename"},
			expectedResult: CmdArgsRawRecognized{Strict: true},
			expectedError:  true,
		},
		{
			name:           "Unknown flag",
			args:           []string{"--unknown-flag"},
			expectedResult: CmdArgsRawRecognized{Strict: true, Others: []string{"--unknown-flag"}},
			expectedError:  false,
		},
		{
			name: "Mix of valid and invalid args",
			args: []string{"--filename=file.yaml", "--envsubst-no-strict", "--unknown"},
			expectedResult: CmdArgsRawRecognized{
				Filenames: []string{"file.yaml"},
				Strict:    false,
				Others:    []string{"--unknown"},
			},
			expectedError: false,
		},
		{
			name:           "Multiple filenames",
			args:           []string{"--filename=file1.yaml", "--filename=file2.yaml"},
			expectedResult: CmdArgsRawRecognized{Strict: true, Filenames: []string{"file1.yaml", "file2.yaml"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes",
			args:           []string{"--envsubst-allowed-prefixes=CI_,APP"},
			expectedResult: CmdArgsRawRecognized{Strict: true, EnvsubstAllowedPrefix: []string{"CI_", "APP"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes (no =)",
			args:           []string{"--envsubst-allowed-prefixes", "CI_,APP"},
			expectedResult: CmdArgsRawRecognized{Strict: true, EnvsubstAllowedPrefix: []string{"CI_", "APP"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes, with append",
			args:           []string{"--envsubst-allowed-prefixes=CI_,APP", "--envsubst-allowed-prefixes=TF_VAR_"},
			expectedResult: CmdArgsRawRecognized{Strict: true, EnvsubstAllowedPrefix: []string{"CI_", "APP", "TF_VAR_"}},
			expectedError:  false,
		},
		{
			name:           "Envsubst allowed prefixes with append (no =)",
			args:           []string{"--envsubst-allowed-prefixes", "CI_,APP", "--envsubst-allowed-prefixes", "TF_VAR_"},
			expectedResult: CmdArgsRawRecognized{Strict: true, EnvsubstAllowedPrefix: []string{"CI_", "APP", "TF_VAR_"}},
			expectedError:  false,
		},
		{
			name:           "Empty value for --filename",
			args:           []string{"--filename="},
			expectedResult: CmdArgsRawRecognized{Strict: true},
			expectedError:  true,
		},
		{
			name:           "Empty value for --envsubst-allowed-vars",
			args:           []string{"--envsubst-allowed-vars="},
			expectedResult: CmdArgsRawRecognized{Strict: true},
			expectedError:  true,
		},
		{
			name:           "Single filename with short flag",
			args:           []string{"-f=file3.yaml"},
			expectedResult: CmdArgsRawRecognized{Strict: true, Filenames: []string{"file3.yaml"}},
			expectedError:  false,
		},
		{
			name:           "Recursive and strict flags",
			args:           []string{"-R", "--envsubst-no-strict"},
			expectedResult: CmdArgsRawRecognized{Recursive: true, Strict: false},
			expectedError:  false,
		},
		{
			name:           "Unrecognized argument without prefix",
			args:           []string{"random-arg"},
			expectedResult: CmdArgsRawRecognized{Strict: true, Others: []string{"random-arg"}},
			expectedError:  false,
		},
		{
			name:           "Multiple unrecognized arguments",
			args:           []string{"random-arg1", "random-arg2"},
			expectedResult: CmdArgsRawRecognized{Strict: true, Others: []string{"random-arg1", "random-arg2"}},
			expectedError:  false,
		},
		{
			name: "Mixed valid and unrecognized arguments",
			args: []string{"--envsubst-no-strict", "random-arg", "--filename=file.yaml"},
			expectedResult: CmdArgsRawRecognized{
				Filenames: []string{"file.yaml"},
				Strict:    false,
				Others:    []string{"random-arg"},
			},
			expectedError: false,
		},
		{
			name:           "Unrecognized argument resembling a flag",
			args:           []string{"-notarealflag"},
			expectedResult: CmdArgsRawRecognized{Strict: true, Others: []string{"-notarealflag"}},
			expectedError:  false,
		},
		{
			name: "Unrecognized argument with spaces",
			args: []string{"random-arg", "--filename=file.yaml", "another-random-arg"},
			expectedResult: CmdArgsRawRecognized{
				Strict:    true,
				Filenames: []string{"file.yaml"},
				Others:    []string{"random-arg", "another-random-arg"},
			},
			expectedError: false,
		},
		{
			name: "Valid args with unrecognized argument that looks like a short flag",
			args: []string{"-R", "-xyz"},
			expectedResult: CmdArgsRawRecognized{
				Strict:    true,
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
