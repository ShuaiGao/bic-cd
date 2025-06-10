package service

import "testing"

func Test_createServiceFile(t *testing.T) {
	type args struct {
		path   string
		config ServiceConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "hello", args: args{
			path: "hello.service",
			config: ServiceConfig{
				Name:        "hello",
				Description: "hello world",
				ExecStart:   "echo \"hello world\"",
				WorkingDir:  "/home/hello",
				User:        "root",
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createServiceFile(tt.args.path, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("createServiceFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
