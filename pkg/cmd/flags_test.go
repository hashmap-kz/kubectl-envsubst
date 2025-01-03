package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestResolveFilenames2(t *testing.T) {
	// Create a temporary test directory structure
	tempDir, err := os.MkdirTemp("", "test_resolveFilenames2")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir) // Clean up after test

	testFile1 := filepath.Join(tempDir, "file1.txt")
	testFile2 := filepath.Join(tempDir, "file2.yaml")
	testSubDir := filepath.Join(tempDir, "subdir")
	testFile3 := filepath.Join(testSubDir, "file3.json")
	testFile4 := filepath.Join(testSubDir, "file4.xml")

	// Create test files and directories
	if err := os.WriteFile(testFile1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(testSubDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile3, []byte("content3"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile4, []byte("content4"), 0644); err != nil {
		t.Fatal(err)
	}

	// Define test cases
	tests := []struct {
		name        string
		path        string
		recursive   bool
		expected    []string
		expectError bool
	}{
		{
			name:        "Single file valid extension",
			path:        testFile2,
			recursive:   false,
			expected:    []string{testFile2},
			expectError: false,
		},
		{
			name:        "Single file invalid extension",
			path:        testFile1,
			recursive:   false,
			expected:    nil,
			expectError: false,
		},
		{
			name:        "Glob pattern valid extensions",
			path:        filepath.Join(tempDir, "*.yaml"),
			recursive:   false,
			expected:    []string{testFile2},
			expectError: false,
		},
		{
			name:        "Directory non-recursive",
			path:        tempDir,
			recursive:   false,
			expected:    []string{testFile2},
			expectError: false,
		},
		{
			name:        "Directory recursive",
			path:        tempDir,
			recursive:   true,
			expected:    []string{testFile2, testFile3},
			expectError: false,
		},
		{
			name:        "URL handling",
			path:        "http://example.com/file.yaml",
			recursive:   false,
			expected:    []string{"http://example.com/file.yaml"},
			expectError: false,
		},
		{
			name:        "Nonexistent file",
			path:        filepath.Join(tempDir, "nonexistent.txt"),
			recursive:   false,
			expected:    nil,
			expectError: true,
		},
	}

	// Execute test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := resolveFilenames(test.path, test.recursive)
			if test.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}
				if !reflect.DeepEqual(result, test.expected) {
					t.Errorf("expected %v, got %v", test.expected, result)
				}
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid HTTP URL", "http://example.com", true},
		{"Valid HTTPS URL", "https://example.com", true},
		{"Invalid URL", "example.com", false},
		{"Empty String", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsURL(tt.input)
			if result != tt.expected {
				t.Errorf("IsURL(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIgnoreFile(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		extensions []string
		expected   bool
	}{
		{"Allowed Extension", "file.yaml", []string{".json", ".yaml"}, false},
		{"Disallowed Extension", "file.txt", []string{".json", ".yaml"}, true},
		{"Empty Extensions", "file.yaml", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ignoreFile(tt.path, tt.extensions)
			if result != tt.expected {
				t.Errorf("ignoreFile(%q, %v) = %v, expected %v", tt.path, tt.extensions, result, tt.expected)
			}
		})
	}
}

func TestResolveSingle(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "file.yaml")
	os.WriteFile(tmpFile, []byte{}, 0644)

	proxy := &CmdFlagsProxy{Filenames: []string{}}
	err := resolveSingle(tmpFile, proxy)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(proxy.Filenames) != 1 || proxy.Filenames[0] != tmpFile {
		t.Errorf("resolveSingle(%q) did not add the file to CmdFlagsProxy.Filenames", tmpFile)
	}
}
