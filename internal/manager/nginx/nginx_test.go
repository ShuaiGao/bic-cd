package nginx

import (
	"testing"
)

func Test_CreateNginxConfig(t *testing.T) {
	type args struct {
		path   string
		config NginxConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "hello", args: args{
			//path: "/etc/nginx/conf.d/hello_v1.conf",
			path: "./hello_v1.conf",
			config: NginxConfig{
				Domain:      "hello.farmergao.cn",
				EnableHTTPS: true,
				SSLCertPath: "./ssl/hello.ssl.cert",
				SSLKeyPath:  "./ssl/hello.ssl.key",
				Port:        8080,
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateNginxConfig(tt.args.path, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("createNginxConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
