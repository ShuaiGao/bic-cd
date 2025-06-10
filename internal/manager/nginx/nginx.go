package nginx

import (
	"bic-cd/internal/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const (
	nginxConfigDir = "/etc/nginx/conf.d/"
)

type Service struct {
	Name string
	IP   string
	Port int
}

type NginxConfig struct {
	Domain          string    `json:"domain"`
	Services        []Service `json:"services"`
	EnableHTTPS     bool      `json:"enable_https"`      // 是否启用HTTPS
	SSLCertPath     string    `json:"ssl_cert_path"`     // SSL证书路径
	SSLKeyPath      string    `json:"ssl_key_path"`      // SSL私钥路径
	HSTS            bool      `json:"hsts"`              // 是否启用HSTS
	SSLProtocols    string    `json:"ssl_protocols"`     // SSL协议版本
	SSLCiphers      string    `json:"ssl_ciphers"`       // SSL密码套件
	SSLPreferServer bool      `json:"ssl_prefer_server"` // SSL优先使用服务器密码套件
}

var nginxConfigTemplate = template.Must(template.New("nginx").Parse(`upstream {{.Domain}}_backend {
{{- range .Services}}
    server {{.IP}}:{{.Port}};
{{- end}}
}

{{- if .EnableHTTPS}}
server {
    listen 80;
    server_name {{.Domain}};
    
    location / {
        proxy_pass http://{{.Domain}}_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 443 ssl http2;
    server_name {{.Domain}};
    
    ssl_certificate {{.SSLCertPath}};
    ssl_certificate_key {{.SSLKeyPath}};
    
    {{- if .HSTS}}
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    {{- end}}
    
    {{- if .SSLProtocols}}
    ssl_protocols {{.SSLProtocols}};
    {{- end}}
    
    {{- if .SSLCiphers}}
    ssl_ciphers {{.SSLCiphers}};
    {{- end}}
    
    {{- if .SSLPreferServer}}
    ssl_prefer_server_ciphers on;
    {{- end}}
    
    location / {
        proxy_pass http://{{.Domain}}_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
        proxy_set_header X-Forwarded-Port 443;
    }
}
{{- else}}
server {
    listen 80;
    server_name {{.Domain}};

    location / {
        proxy_pass http://{{.Domain}}_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
{{- end}}
`))

func configureNginx(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var config NginxConfig
	if err := json.Unmarshal(body, &config); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validateNginxConfig(config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	configFile := filepath.Join(nginxConfigDir, config.Domain+".conf")
	if err := createNginxConfig(configFile, config); err != nil {
		http.Error(w, "Failed to create nginx configuration", http.StatusInternalServerError)
		return
	}

	if err := executeNginx("reload"); err != nil {
		http.Error(w, "Failed to reload nginx", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Nginx configured for %s", config.Domain)
}

func validateNginxConfig(config NginxConfig) error {
	if config.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if len(config.Services) == 0 {
		return fmt.Errorf("at least one service is required")
	}
	for _, s := range config.Services {
		if !util.IsValidName(s.Name) {
			return fmt.Errorf("invalid service name: %s", s.Name)
		}
	}
	return nil
}

// 使用模板生成Nginx配置文件
func createNginxConfig(path string, config NginxConfig) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return nginxConfigTemplate.Execute(file, config)
}

func executeNginx(args ...string) error {
	cmd := exec.Command("nginx", append([]string{"-s"}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Nginx error: %s, Output: %s", err, string(output))
		return err
	}
	return nil
}
