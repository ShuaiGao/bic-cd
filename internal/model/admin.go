package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string   `gorm:"size:255;comment:姓名"`
	Username     string   `gorm:"size:56;unique;index;comment:用户名"`
	Phone        string   `gorm:"size:32;comment:电话号码"`
	Email        string   `gorm:"size:56;comment:邮箱"`
	Password     string   `gorm:"size:64;comment:密码"`
	Avatar       string   `gorm:"comment:用户头像url"`
	Groups       []*Group `gorm:"many2many:user_group;comment:用户关联权限组"`
	Forbidden    bool     `gorm:"comment:是否禁用"`
	UpdateUserID uint     `gorm:"comment:修改用户外键ID;default:null;"`
	UpdateUser   *User    `gorm:"default:galeone;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
}

type Group struct {
	gorm.Model
	Key          string       `gorm:"size:128;unique;index;comment:唯一值"`
	Alias        string       `gorm:"comment:别名"`
	Rank         int          `gorm:"comment:排序值"`
	Permissions  []Permission `gorm:"many2many:group_permission;comment:权限组对应权限列表;foreignKey:Key;References:Key"`
	FatherID     uint         `gorm:"comment:父角色;default:null;"`
	Father       *Group       `gorm:"default:galeone;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	CreateUserID uint         `gorm:"comment:创建用户ID外键"`
	CreateUser   User         ``
	UpdateUserID uint         `gorm:"comment:修改用户ID外键"`
	UpdateUser   User         ``
}

type PermissionType string

const (
	DATA PermissionType = "data" // 数据
	API  PermissionType = "api"  // api
	Menu PermissionType = "menu" // client-菜单
	Page PermissionType = "page" // client-page
)

type Permission struct {
	gorm.Model
	Key    string         `gorm:"size:128;unique;index;comment:唯一值"`
	Type   PermissionType `gorm:"comment:权限类型"`
	Obj    string         `gorm:"comment:父节点;uniqueIndex:idx_obj_action;size:56"`
	Action string         `gorm:"uniqueIndex:idx_obj_action;size:56"`
	Alias  string         `gorm:"comment:权限别名"`
	Rank   int            `gorm:"comment:排序值"`
}
