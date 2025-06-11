package main

import (
	"context"
	"fmt"
	"os"

	"temp-mailbox-service/internal/application"
	"temp-mailbox-service/internal/domain/user"
	"temp-mailbox-service/internal/infrastructure/auth"
	"temp-mailbox-service/internal/infrastructure/config"
	"temp-mailbox-service/internal/infrastructure/database"
	"temp-mailbox-service/internal/infrastructure/persistence"
)

func main() {
	fmt.Println("🧪 测试真实SQLite数据库连接...")

	// 1. 加载配置（如果没有配置文件则使用默认值）
	fmt.Println("📝 加载配置...")
	cfg, err := config.Load("")
	if err != nil {
		fmt.Printf("❌ 配置加载失败: %v\n", err)
		os.Exit(1)
	}
	
	// 覆盖数据库配置为测试配置 
	cfg.Database.Driver = "sqlite"
	cfg.Database.DSN = "./test.db"
	cfg.JWT.Secret = "test-jwt-secret-key-for-testing-only-2024"

	// 2. 尝试连接数据库
	fmt.Println("🗄️  连接SQLite数据库...")
	fmt.Printf("数据库配置: Driver=%s, DSN=%s\n", cfg.Database.Driver, cfg.Database.DSN)
	
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		fmt.Printf("❌ 数据库连接失败: %v\n", err)
		fmt.Printf("错误类型: %T\n", err)
		os.Exit(1)
	}
	defer database.Close(db)
	fmt.Println("✅ SQLite数据库连接成功")

	// 3. 执行数据库迁移
	fmt.Println("🔄 执行数据库迁移...")
	if err := database.Migrate(db); err != nil {
		fmt.Printf("❌ 数据库迁移失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 数据库迁移完成")

	// 4. 初始化服务
	fmt.Println("⚙️  初始化服务...")
	userRepo := persistence.NewUserRepository(db)
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.Issuer,
	)
	userService := application.NewUserService(userRepo, jwtService)
	fmt.Println("✅ 服务初始化完成")

	// 5. 测试数据库操作
	fmt.Println("👤 测试用户注册（写入数据库）...")
	ctx := context.Background()
	registerReq := &user.CreateUserRequest{
		Username:  "dbtest",
		Email:     "dbtest@example.com",
		Password:  "password123",
		FirstName: "Database",
		LastName:  "Test",
	}

	registeredUser, err := userService.Register(ctx, registerReq)
	if err != nil {
		fmt.Printf("❌ 用户注册失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 用户注册成功，ID: %d（已保存到SQLite）\n", registeredUser.ID)

	// 6. 测试登录（从数据库读取）
	fmt.Println("🔐 测试用户登录（从数据库验证）...")
	loginReq := &user.LoginRequest{
		Email:    "dbtest@example.com",
		Password: "password123",
	}

	_, err = userService.Login(ctx, loginReq)
	if err != nil {
		fmt.Printf("❌ 用户登录失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 用户登录成功（从SQLite数据库验证）\n")

	// 7. 检查数据库文件是否创建
	fmt.Println("📁 检查数据库文件...")
	if _, err := os.Stat("./test.db"); err == nil {
		fmt.Println("✅ 数据库文件 'test.db' 已创建")
	} else {
		fmt.Printf("❌ 数据库文件不存在: %v\n", err)
	}

	fmt.Println("\n🎉 SQLite数据库测试完成！")
	fmt.Println("📊 验证项目:")
	fmt.Println("   ✅ SQLite数据库连接")
	fmt.Println("   ✅ 数据库表创建（迁移）")
	fmt.Println("   ✅ 用户数据写入")
	fmt.Println("   ✅ 用户数据读取")
	fmt.Println("   ✅ 数据库文件持久化")
} 