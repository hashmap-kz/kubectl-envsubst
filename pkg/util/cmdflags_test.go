package util

import (
	"reflect"
	"testing"
)

func TestParseCmdFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *CmdFlagsProxy
		wantErr bool
	}{
		{
			name: "No arguments, default namespace",
			args: []string{},
			want: &CmdFlagsProxy{
				Filenames: []string{},
				Namespace: "default",
				Others:    []string{},
			},
			wantErr: false,
		},
		{
			name: "Single filename and namespace",
			args: []string{"--filename", "file1.yaml", "--namespace", "test-namespace"},
			want: &CmdFlagsProxy{
				Filenames: []string{"file1.yaml"},
				Namespace: "test-namespace",
				Others:    []string{},
			},
			wantErr: false,
		},
		{
			name: "Multiple filenames",
			args: []string{"-f", "file1.yaml", "-f", "file2.yaml"},
			want: &CmdFlagsProxy{
				Filenames: []string{"file1.yaml", "file2.yaml"},
				Namespace: "default",
				Others:    []string{},
			},
			wantErr: false,
		},
		{
			name: "Unrecognized arguments",
			args: []string{"--filename", "file1.yaml", "extra1", "extra2"},
			want: &CmdFlagsProxy{
				Filenames: []string{"file1.yaml"},
				Namespace: "default",
				Others:    []string{"extra1", "extra2"},
			},
			wantErr: false,
		},
		{
			name:    "Missing value for filename",
			args:    []string{"--filename"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Missing value for namespace",
			args:    []string{"--namespace"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Mixed flags and extra arguments-1",
			args: []string{"--filename", "file1.yaml", "-n", "test-namespace", "extra1"},
			want: &CmdFlagsProxy{
				Filenames: []string{"file1.yaml"},
				Namespace: "test-namespace",
				Others:    []string{"extra1"},
			},
			wantErr: false,
		},
		{
			name: "Mixed flags and extra arguments-2",
			args: []string{"--filename", "file1.yaml", "-n", "test-namespace", "extra1", "--dry-run=client", "-oyaml"},
			want: &CmdFlagsProxy{
				Filenames: []string{"file1.yaml"},
				Namespace: "test-namespace",
				Others:    []string{"extra1", "--dry-run=client", "-oyaml"},
			},
			wantErr: false,
		},
		{
			name: "Namespace override",
			args: []string{"--filename", "file1.yaml", "-n", "test-namespace", "-n", "tmp", "extra1", "--dry-run=client", "-oyaml"},
			want: &CmdFlagsProxy{
				Filenames: []string{"file1.yaml"},
				Namespace: "tmp",
				Others:    []string{"extra1", "--dry-run=client", "-oyaml"},
			},
			wantErr: false,
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
