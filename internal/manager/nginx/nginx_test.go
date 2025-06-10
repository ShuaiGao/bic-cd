package nginx

import "testing"

func Test_createNginxConfig(t *testing.T) {
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
			path: "/etc/nginx/conf.d/hello_v1.conf",
			config: NginxConfig{
				Domain:      "hello.farmergao.cn",
				EnableHTTPS: true,
				SSLCertPath: "./ssl/hello.ssl.cert",
				SSLKeyPath:  "./ssl/hello.ssl.key",
				Services:    []Service{{Name: "hello.service", IP: "127.0.0.1", Port: 7788}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createNginxConfig(tt.args.path, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("createNginxConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
