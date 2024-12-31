package util

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseCmdFlags(t *testing.T) {
	// Test data setup
	testDir := "testdata_tmpdir"
	subDir := filepath.Join(testDir, "subdir")
	file1 := filepath.Join(testDir, "file1.yaml")
	file2 := filepath.Join(subDir, "file2.yaml")

	os.MkdirAll(subDir, 0755)
	defer os.RemoveAll(testDir)

	os.WriteFile(file1, []byte("content1"), 0644)
	os.WriteFile(file2, []byte("content2"), 0644)

	tests := []struct {
		name    string
		args    []string
		want    *CmdFlagsProxy
		wantErr bool
	}{
		{
			name: "Single file",
			args: []string{"--filename", file1},
			want: &CmdFlagsProxy{
				Filenames: []string{file1},
				Others:    []string{},
				Recursive: false,
			},
			wantErr: false,
		},
		{
			name: "Single file with equal (long form)",
			args: []string{"--filename=" + file1},
			want: &CmdFlagsProxy{
				Filenames: []string{file1},
				Others:    []string{},
				Recursive: false,
			},
			wantErr: false,
		},
		{
			name: "Single file with equal (short form)",
			args: []string{"-f=" + file1},
			want: &CmdFlagsProxy{
				Filenames: []string{file1},
				Others:    []string{},
				Recursive: false,
			},
			wantErr: false,
		},
		{
			name: "Recursive directory",
			args: []string{"--filename", testDir, "--recursive"},
			want: &CmdFlagsProxy{
				Filenames: []string{file1, file2},
				Others:    []string{},
				Recursive: true,
			},
			wantErr: false,
		},
		{
			name: "Glob pattern",
			args: []string{"--filename", filepath.Join(testDir, "*.yaml")},
			want: &CmdFlagsProxy{
				Filenames: []string{file1},
				Others:    []string{},
				Recursive: false,
			},
			wantErr: false,
		},
		{
			name: "Namespace and extra args",
			args: []string{"--filename", file1, "--namespace", "dev", "extra-arg"},
			want: &CmdFlagsProxy{
				Filenames: []string{file1},
				Others:    []string{"--namespace", "dev", "extra-arg"},
				Recursive: false,
			},
			wantErr: false,
		},
		{
			name: "URL input",
			args: []string{"--filename", "http://example.com/file.yaml"},
			want: &CmdFlagsProxy{
				Filenames: []string{"http://example.com/file.yaml"},
				Others:    []string{},
				Recursive: false,
			},
			wantErr: false,
		},
		{
			name:    "Missing value for filename",
			args:    []string{"--filename"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCmdFlags(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCmdFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCmdFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
