package database

import (
	"fmt"
	"time"

	"temp-mailbox-service/internal/infrastructure/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *config.DatabaseConfig) error {
	var dialector gorm.Dialector
	
	// 根据驱动类型选择对应的dialector
	switch cfg.Driver {
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	default:
		return fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
	}
	
	// GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: nil, // 使用默认命名策略
	}
	
	// 连接数据库
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	
	// 获取底层sql.DB以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", err)
	}
	
	// 配置连接池
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Minute)
	
	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}
	
	// 设置全局数据库实例
	DB = db
	
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if DB == nil {
		return nil
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", err)
	}
	
	return sqlDB.Close()
} 