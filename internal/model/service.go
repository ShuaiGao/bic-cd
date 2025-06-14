package model

import (
	"fmt"
	"gorm.io/gorm"
	"path"
	"strings"
)

type Service struct {
	gorm.Model
	Name        string `gorm:"unique;index;size:255"`
	Description string `gorm:"comment:服务描述"`
	WorkingDir  string `gorm:"comment:工作目录"`
	User        string `gorm:"comment:运行用户"`
	PortMin     uint16 `gorm:"comment:最小端口号"`
	PortMax     uint16 `gorm:"comment:最大端口号"`
	Config      string `gorm:"comment:配置文件"`
}

func (s *Service) GetService(version string) string {
	return fmt.Sprintf("%s-%s.service", s.Name, version)
}

type ServiceInstance struct {
	gorm.Model
	ServiceID uint
	Service   Service
	ExecStart string `gorm:"comment:启动命令"`
	Port      uint16 `gorm:"comment:端口"`
	Version   string `gorm:"comment:版本;size:32"`
}

func (i *ServiceInstance) GetService() string {
	v := strings.Replace(i.Version, ".", "-", -1)
	return fmt.Sprintf("%s-%s.service", i.Service.Name, v)
}

func (i *ServiceInstance) SetExecStart(port uint16) {
	i.Port = port
	execStart := path.Join(i.Service.WorkingDir, i.Service.Name)
	i.ExecStart = fmt.Sprintf("%s bic --port %d", execStart, port)
}
