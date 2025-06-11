package model

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Name        string `gorm:"unique;index;size:255"`
	Description string `gorm:"comment:服务描述"`
	ExecStart   string `gorm:"comment:启动命令"`
	WorkingDir  string `gorm:"comment:工作目录"`
	User        string `gorm:"comment:运行用户"`
	PortMin     uint16 `gorm:"comment:最小端口号"`
	PortMax     uint16 `gorm:"comment:最大端口号"`
	Config      string `gorm:"comment:配置文件"`
}

type ServiceInstance struct {
	gorm.Model
	ServiceID uint
	Service   Service
	Port      uint16 `gorm:"comment:端口"`
	Version   string `gorm:"comment:版本;size:32"`
}
