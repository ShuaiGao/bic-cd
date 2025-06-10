package service

import (
	"os"
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

type ServiceConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ExecStart   string `json:"exec_start"`
	WorkingDir  string `json:"working_dir,omitempty"`
	User        string `json:"user,omitempty"`
}

// 使用模板生成systemd服务文件
func createServiceFile(path string, config ServiceConfig) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return systemdServiceTemplate.Execute(file, config)
}
