package util

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseCmdFlags(t *testing.T) {
	// Test data setup
	testDir := "testdata"
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
				Namespace: "default",
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
				Namespace: "default",
				Others:    []string{},
				Recursive: true,
			},
			wantErr: false,
		},
		//{
		//	name: "Glob pattern",
		//	args: []string{"--filename", filepath.Join(testDir, "*.yaml")},
		//	want: &CmdFlagsProxy{
		//		Filenames: []string{file1},
		//		Namespace: "default",
		//		Others:    []string{},
		//		Recursive: false,
		//	},
		//	wantErr: false,
		//},
		{
			name: "Namespace and extra args",
			args: []string{"--filename", file1, "--namespace", "dev", "extra-arg"},
			want: &CmdFlagsProxy{
				Filenames: []string{file1},
				Namespace: "dev",
				Others:    []string{"extra-arg"},
				Recursive: false,
			},
			wantErr: false,
		},
		{
			name: "URL input",
			args: []string{"--filename", "http://example.com/file.yaml"},
			want: &CmdFlagsProxy{
				Filenames: []string{"http://example.com/file.yaml"},
				Namespace: "default",
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
