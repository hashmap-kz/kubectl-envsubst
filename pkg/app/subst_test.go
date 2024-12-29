package app

import (
	"os"
	"testing"
)

func TestSubstituteEnvs(t *testing.T) {
	// Set environment variables for testing
	_ = os.Setenv("FOO", "foo_value")
	_ = os.Setenv("BAR", "bar_value")

	tests := []struct {
		name        string
		text        string
		allowedEnvs []string
		expected    string
	}{
		{
			name:        "Basic substitution with single variable",
			text:        "Hello $FOO!",
			allowedEnvs: []string{"FOO"},
			expected:    "Hello foo_value!",
		},
		{
			name:        "Basic substitution with multiple variables",
			text:        "Hello $FOO and ${BAR}!",
			allowedEnvs: []string{"FOO", "BAR"},
			expected:    "Hello foo_value and bar_value!",
		},
		{
			name:        "Variable not in allowed list",
			text:        "Hello $FOO and $BAR!",
			allowedEnvs: []string{"FOO"},
			expected:    "Hello foo_value and $BAR!",
		},
		{
			name:        "No variables allowed",
			text:        "Hello $FOO!",
			allowedEnvs: []string{},
			expected:    "Hello $FOO!",
		},
		{
			name:        "Variable not set in environment",
			text:        "Hello $UNSET_VAR!",
			allowedEnvs: []string{"UNSET_VAR"},
			expected:    "Hello $UNSET_VAR!",
		},
		{
			name:        "No substitution needed",
			text:        "Hello world!",
			allowedEnvs: []string{"FOO", "BAR"},
			expected:    "Hello world!",
		},
		{
			name:        "Empty text",
			text:        "",
			allowedEnvs: []string{"FOO", "BAR"},
			expected:    "",
		},
		{
			name:        "Partially valid and invalid variables",
			text:        "Hello ${FOO} and $BAZ!",
			allowedEnvs: []string{"FOO"},
			expected:    "Hello foo_value and $BAZ!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SubstituteEnvs(test.text, test.allowedEnvs)
			if result != test.expected {
				t.Errorf("Test %q failed: expected %q, got %q", test.name, test.expected, result)
			}
		})
	}
}
