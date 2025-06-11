package database

import (
	"fmt"

	"temp-mailbox-service/internal/domain/user"
)

// Migrate 执行数据库迁移
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}
	
	// 自动迁移所有模型
	err := DB.AutoMigrate(
		&user.User{},
		// 在这里添加其他需要迁移的模型
	)
	
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}
	
	fmt.Println("数据库迁移成功完成")
	return nil
}

// DropTables 删除所有表（危险操作，仅用于开发环境）
func DropTables() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}
	
	// 删除所有表
	err := DB.Migrator().DropTable(
		&user.User{},
		// 在这里添加其他需要删除的表
	)
	
	if err != nil {
		return fmt.Errorf("删除表失败: %w", err)
	}
	
	fmt.Println("所有表已删除")
	return nil
}

// CreateIndexes 创建索引
func CreateIndexes() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}
	
	// 为用户表创建额外索引（如果需要）
	if !DB.Migrator().HasIndex(&user.User{}, "email") {
		if err := DB.Migrator().CreateIndex(&user.User{}, "email"); err != nil {
			return fmt.Errorf("创建邮箱索引失败: %w", err)
		}
	}
	
	if !DB.Migrator().HasIndex(&user.User{}, "username") {
		if err := DB.Migrator().CreateIndex(&user.User{}, "username"); err != nil {
			return fmt.Errorf("创建用户名索引失败: %w", err)
		}
	}
	
	fmt.Println("索引创建完成")
	return nil
} 