package main

import (
	"context"
	"fmt"
	"log"

	"temp-mailbox-service/internal/application"
	"temp-mailbox-service/internal/domain/user"
	"temp-mailbox-service/internal/infrastructure/auth"
	"temp-mailbox-service/internal/infrastructure/config"
	"temp-mailbox-service/internal/infrastructure/database"
	"temp-mailbox-service/internal/infrastructure/persistence"
)

func main() {
	fmt.Println("🧪 开始测试用户系统核心组件...")

	// 1. 使用测试配置
	fmt.Println("📝 使用测试配置...")
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver:       "sqlite",
			DSN:          "./test.db",
			MaxOpenConns: 25,
			MaxIdleConns: 10,
			MaxLifetime:  30,
		},
		JWT: config.JWTConfig{
			Secret:          "test-jwt-secret-key-for-testing-only-2024",
			AccessTokenTTL:  60,
			RefreshTokenTTL: 10080,
			Issuer:          "temp-mailbox-test",
		},
	}
	fmt.Printf("✅ 配置加载成功，数据库: %s\n", cfg.Database.Driver)

	// 2. 连接数据库
	fmt.Println("🗄️  连接数据库...")
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}
	defer database.Close(db)
	fmt.Println("✅ 数据库连接成功")

	// 3. 执行数据库迁移
	fmt.Println("🔄 执行数据库迁移...")
	if err := database.Migrate(db); err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err)
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

	// 5. 测试用户注册
	fmt.Println("👤 测试用户注册...")
	ctx := context.Background()
	registerReq := &user.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	registeredUser, err := userService.Register(ctx, registerReq)
	if err != nil {
		log.Fatalf("❌ 用户注册失败: %v", err)
	}
	fmt.Printf("✅ 用户注册成功，ID: %d, 用户名: %s\n", registeredUser.ID, registeredUser.Username)

	// 6. 测试用户登录
	fmt.Println("🔐 测试用户登录...")
	loginReq := &user.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResp, err := userService.Login(ctx, loginReq)
	if err != nil {
		log.Fatalf("❌ 用户登录失败: %v", err)
	}
	fmt.Printf("✅ 用户登录成功，访问令牌长度: %d\n", len(loginResp.Tokens.AccessToken))

	// 7. 测试JWT令牌验证
	fmt.Println("🔑 测试JWT令牌验证...")
	claims, err := jwtService.ValidateAccessToken(loginResp.Tokens.AccessToken)
	if err != nil {
		log.Fatalf("❌ JWT令牌验证失败: %v", err)
	}
	fmt.Printf("✅ JWT令牌验证成功，用户ID: %d, 邮箱: %s\n", claims.UserID, claims.Email)

	// 8. 测试获取用户资料
	fmt.Println("📋 测试获取用户资料...")
	profile, err := userService.GetProfile(ctx, registeredUser.ID)
	if err != nil {
		log.Fatalf("❌ 获取用户资料失败: %v", err)
	}
	fmt.Printf("✅ 获取用户资料成功，全名: %s %s\n", profile.FirstName, profile.LastName)

	// 9. 测试更新用户资料
	fmt.Println("✏️  测试更新用户资料...")
	updateReq := &user.UpdateUserRequest{
		FirstName: "Updated",
		LastName:  "Name",
		TimeZone:  "Asia/Shanghai",
		Language:  "zh-CN",
	}

	updatedProfile, err := userService.UpdateProfile(ctx, registeredUser.ID, updateReq)
	if err != nil {
		log.Fatalf("❌ 更新用户资料失败: %v", err)
	}
	fmt.Printf("✅ 更新用户资料成功，新全名: %s %s\n", updatedProfile.FirstName, updatedProfile.LastName)

	// 10. 测试令牌刷新
	fmt.Println("🔄 测试令牌刷新...")
	newTokens, err := userService.RefreshToken(ctx, loginResp.Tokens.RefreshToken)
	if err != nil {
		log.Fatalf("❌ 令牌刷新失败: %v", err)
	}
	fmt.Printf("✅ 令牌刷新成功，新访问令牌长度: %d\n", len(newTokens.AccessToken))

	fmt.Println("\n🎉 所有测试通过！用户系统核心组件运行正常！")
	fmt.Println("📊 测试总结:")
	fmt.Println("   ✅ 配置管理")
	fmt.Println("   ✅ 数据库连接和迁移")
	fmt.Println("   ✅ 用户注册")
	fmt.Println("   ✅ 用户登录")
	fmt.Println("   ✅ JWT令牌生成和验证")
	fmt.Println("   ✅ 用户资料管理")
	fmt.Println("   ✅ 令牌刷新")
} 