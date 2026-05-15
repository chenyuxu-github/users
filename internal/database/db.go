package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"user/internal/model"
)

// Init 根据传入的 MySQL DSN 初始化 GORM 数据库连接，并自动迁移数据表。
func Init(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate 会在 users 表不存在时自动创建，并补齐模型中新增的字段。
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
