package db

import (
	"bic-cd/pkg/config"
	"bic-cd/pkg/log"
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

var db *gorm.DB

func DB() *gorm.DB {
	return db
}

func WithContext(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

func Setup() {
	// 连接数据库，若数据库不存在则会自动创建
	sqlitePath := config.GlobalConf.App.SqlitePath
	// 获取文件的目录路径
	dir := filepath.Dir(sqlitePath)
	// 创建所有必要的目录（包括中间目录）
	// 0755 是目录权限：所有者可读可写可执行，组和其他可读可执行
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic("创建目录失败: " + err.Error())
	}
	var err error
	db, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{
		Logger:          log.ExportGormLogger(time.Second*3, 100*1000),
		CreateBatchSize: 100,
	})
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}
}
