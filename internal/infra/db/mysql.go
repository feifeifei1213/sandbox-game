package db

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	appconfig "sandbox-game/internal/config"
)

// NewMySQL 使用项目配置创建首版 MySQL 连接。
func NewMySQL(cfg *appconfig.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db from gorm: %w", err)
	}

	applyMySQLPool(sqlDB, cfg)

	return db, nil
}

func applyMySQLPool(sqlDB *sql.DB, cfg *appconfig.Config) {
	if cfg.MySQL.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	}
	if cfg.MySQL.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	}
}
