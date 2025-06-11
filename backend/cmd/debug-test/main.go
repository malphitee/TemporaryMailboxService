package main

import (
	"fmt"
	"os"

	"temp-mailbox-service/internal/infrastructure/config"
)

func main() {
	fmt.Println("🔍 调试配置加载...")

	// 测试配置加载
	fmt.Println("正在加载配置...")
	cfg, err := config.Load("")
	if err != nil {
		fmt.Printf("配置加载错误: %v\n", err)
		fmt.Printf("错误类型: %T\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 配置加载成功\n")
	fmt.Printf("数据库驱动: %s\n", cfg.Database.Driver)
	fmt.Printf("数据库DSN: %s\n", cfg.Database.DSN)
	fmt.Printf("JWT密钥: %s\n", cfg.JWT.Secret[:20]+"...")
} 