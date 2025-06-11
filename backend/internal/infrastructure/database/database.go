package database

import (
	"fmt"
	"time"

	"temp-mailbox-service/internal/domain/user"
	"temp-mailbox-service/internal/infrastructure/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect 连接数据库
func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector
	
	switch cfg.Driver {
	case "sqlite":
		// 使用SQLite驱动，这里GORM会自动选择可用的驱动
		dialector = sqlite.Open(cfg.DSN)
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
	}
	
	// GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// 禁用外键约束检查以避免一些兼容性问题
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	
	// 连接数据库
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}
	
	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}
	
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Minute)
	
	return db, nil
}

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) error {
	// 自动迁移用户表
	if err := db.AutoMigrate(&user.User{}); err != nil {
		return fmt.Errorf("用户表迁移失败: %w", err)
	}
	
	// TODO: 添加其他表的迁移
	// 例如：临时邮箱、域名、邮件等表
	
	return nil
}

// Close 关闭数据库连接
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
} 