package db

import (
	"bic-cd/pkg/config"
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	var err error
	db, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Warn),
		CreateBatchSize: 100,
	})
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}
}
