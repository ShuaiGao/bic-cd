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
Description={{.Service.Description}}
After=network.target

[Service]
Type=simple
ExecStart={{.ExecStart}}
{{- if .Service.WorkingDir}}
WorkingDirectory={{.Service.WorkingDir}}
{{- end}}
{{- if .Service.User}}
User={{.Service.User}}
{{- end}}

[Install]
WantedBy=multi-user.target
`))

type Config struct {
	Instance model.ServiceInstance
}

// CreateServiceFile 使用模板生成systemd服务文件
func CreateServiceFile(path string, config Config) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return systemdServiceTemplate.Execute(file, config.Instance)
}

func CreateService(config Config) (string, error) {
	var buffer bytes.Buffer
	err := systemdServiceTemplate.Execute(&buffer, config.Instance)
	return buffer.String(), err
}

// GetAvailablePort 获取一个可用的port
func GetAvailablePort(minPort, maxPort uint16) (uint16, error) {
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

func BuildService(config Config) error {
	// 生成服务文件
	servicePath := filepath.Join(systemdDir, config.Instance.GetService())
	fmt.Println("Building service", servicePath)
	if err := CreateServiceFile(servicePath, config); err != nil {
		fmt.Println("Building service err", err)
		return err
	}
	// reload service config
	fmt.Println("reload service", config.Instance.GetService())
	if err := reloadService(); err != nil {
		fmt.Println("reload service err", err)
		return err
	}
	fmt.Println("enable service", config.Instance.GetService())
	// enable service
	if err := EnableService(config.Instance); err != nil {
		fmt.Println("enable service err", err)
		return err
	}
	fmt.Println("start service", config.Instance.GetService())
	// start service
	if err := StartService(config.Instance); err != nil {
		fmt.Println("start service err", err)
		return err
	}
	return nil
}

func reloadService() error {
	cmd := exec.Command("systemctl", "daemon-reload")
	return cmd.Run()
}

func EnableService(service model.ServiceInstance) error {
	cmd := exec.Command("systemctl", "enable", service.GetService())
	return cmd.Run()
}

func StartService(service model.ServiceInstance) error {
	cmd := exec.Command("systemctl", "start", service.GetService())
	return cmd.Run()
}

func StopService(service model.ServiceInstance) error {
	cmd := exec.Command("systemctl", "stop", service.GetService())
	return cmd.Run()
}

func StatusService(service model.ServiceInstance) (string, error) {
	fmt.Println("service status 111 ", service.GetService())
	cmd := exec.Command("systemctl", "status", service.GetService())
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func RemoveService(instance model.ServiceInstance) error {
	cmd := exec.Command("systemctl", "stop", instance.GetService())
	if err := cmd.Run(); err != nil {
		return err
	}
	servicePath := filepath.Join(systemdDir, instance.GetService())
	if err := os.RemoveAll(servicePath); err != nil {
		return err
	}
	return reloadService()
}
