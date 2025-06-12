package service

import (
	"bic-cd/internal/model"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const systemdDir = "/etc/systemd/system/"

// 添加systemd服务模板
var systemdServiceTemplate = template.Must(template.New("systemd").Parse(`[Unit]
Description={{.Description}}
After=network.target

[Service]
Type=simple
ExecStart={{.ExecStart}}
{{- if .WorkingDir}}
WorkingDirectory={{.WorkingDir}}
{{- end}}
{{- if .User}}
User={{.User}}
{{- end}}

[Install]
WantedBy=multi-user.target
`))

type Config struct {
	Service  model.Service
	Instance model.ServiceInstance
}

// 使用模板生成systemd服务文件
func createServiceFile(path string, config Config) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return systemdServiceTemplate.Execute(file, config.Service)
}

func CreateService(config Config) (string, error) {
	var buffer bytes.Buffer
	err := systemdServiceTemplate.Execute(&buffer, config.Service)
	return buffer.String(), err
}

// getAvaliablePort 获取一个可用的port
func getAvaliablePort(minPort, maxPort int16) (int16, error) {
	port := minPort
	for port <= maxPort {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			port++
			continue
		}
		listener.Close()
		return port, nil
	}
	return 0, fmt.Errorf("no avaliable port")
}

func buildService(config Config, minPort, maxPort int16) (string, error) {
	port, err := getAvaliablePort(minPort, maxPort)
	if err != nil {
		return "", err
	}
	fmt.Println(port)
	// 生成服务文件
	serviceFile := fmt.Sprintf("%s.service", config.Service.Name)
	servicePath := filepath.Join(systemdDir, serviceFile)
	err = createServiceFile(servicePath, config)
	if err != nil {
		return "", err
	}
	// enable service

	// 启动服务
	err = startService(servicePath)
	if err != nil {
		return "", err
	}
	return "", nil
}
func enableService(servicePath string) error {
	cmd := exec.Command("systemctl", "enable", servicePath)
	return cmd.Run()
}

func startService(servicePath string) error {
	cmd := exec.Command("systemctl", "start", servicePath)
	return cmd.Run()
}
