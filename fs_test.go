package gutils

import (
	"testing"
)

func TestMkdirAll(t *testing.T) {
	type args struct {
		path    string
		options []MkDirOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exits",
			args: args{
				path:    "tmp/demo/demo/demo.png",
				options: []MkDirOptions{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MkdirAll(tt.args.path, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("MkdirAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
