package model

import (
	"fmt"
	"gorm.io/gorm"
	"path"
	"strings"
)

type Service struct {
	gorm.Model
	Domain      string `gorm:"comment:域名"`
	Name        string `gorm:"comment:服务名;size:128;uniqueIndex:idx_name,expression:CASE WHEN deleted_at IS NULL THEN name ELSE NULL END"`
	Description string `gorm:"comment:服务描述"`
	WorkingDir  string `gorm:"comment:工作目录"`
	User        string `gorm:"comment:运行用户"`
	PortMin     uint16 `gorm:"comment:最小端口号"`
	PortMax     uint16 `gorm:"comment:最大端口号"`
	Config      string `gorm:"comment:配置文件"`
	Version     string `gorm:"comment:当前服务版本号"`
	Instances   []*ServiceInstance
}

func (s *Service) GetInstanceName(inst *ServiceInstance) string {
	v := strings.Replace(inst.Version, ".", "-", -1)
	return fmt.Sprintf("%s-%s.service", s.Name, v)
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

func (i *ServiceInstance) GetNginx() string {
	return fmt.Sprintf("%s.conf", i.Service.Name)
}

func (i *ServiceInstance) GetName() string {
	v := strings.Replace(i.Version, ".", "-", -1)
	return fmt.Sprintf("%s-%s", i.Service.Name, v)
}

func (i *ServiceInstance) GetService() string {
	if i.Service.ID == 0 {
		return ""
	}
	return i.Service.GetInstanceName(i)
}

func (i *ServiceInstance) SetExecStart(version string, port uint16) {
	i.Port = port
	execStart := path.Join(i.Service.WorkingDir, version)
	i.ExecStart = fmt.Sprintf("%s bic --port %d", execStart, port)
}
