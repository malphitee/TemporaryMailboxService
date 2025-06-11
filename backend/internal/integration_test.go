package internal

import (
	"context"
	"os"
	"testing"

	"temp-mailbox-service/internal/application"
	"temp-mailbox-service/internal/domain/user"
	"temp-mailbox-service/internal/infrastructure/auth"
	"temp-mailbox-service/internal/infrastructure/config"
	"temp-mailbox-service/internal/infrastructure/database"
	"temp-mailbox-service/internal/infrastructure/persistence"
)

// TestUserSystemIntegration 用户系统集成测试
func TestUserSystemIntegration(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("CGO_ENABLED", "1")
	
	// 1. 加载配置
	t.Log("🔧 加载测试配置...")
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("配置加载失败: %v", err)
	}
	
	// 使用测试数据库配置
	cfg.Database.Driver = "sqlite"
	cfg.Database.DSN = "./integration_test.db"
	cfg.JWT.Secret = "test-jwt-secret-key-for-testing-only-2024"
	
	// 清理测试数据库
	os.Remove("./integration_test.db")
	defer os.Remove("./integration_test.db")

	// 2. 测试数据库连接
	t.Log("🗄️  测试数据库连接...")
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
	}
	defer database.Close(db)
	t.Log("✅ 数据库连接成功")

	// 3. 测试数据库迁移
	t.Log("🔄 测试数据库迁移...")
	if err := database.Migrate(db); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}
	t.Log("✅ 数据库迁移完成")

	// 4. 初始化服务
	t.Log("⚙️  初始化服务...")
	userRepo := persistence.NewUserRepository(db)
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.Issuer,
	)
	userService := application.NewUserService(userRepo, jwtService)
	t.Log("✅ 服务初始化完成")

	ctx := context.Background()

	// 5. 测试用户注册
	t.Log("👤 测试用户注册...")
	registerReq := &user.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	registeredUser, err := userService.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("用户注册失败: %v", err)
	}
	
	if registeredUser.ID == 0 {
		t.Error("注册用户ID不应该为0")
	}
	if registeredUser.Username != "testuser" {
		t.Errorf("期望用户名为 'testuser'，得到 '%s'", registeredUser.Username)
	}
	if registeredUser.Email != "test@example.com" {
		t.Errorf("期望邮箱为 'test@example.com'，得到 '%s'", registeredUser.Email)
	}
	t.Logf("✅ 用户注册成功，ID: %d", registeredUser.ID)

	// 6. 测试重复注册（应该失败）
	t.Log("🔄 测试重复注册...")
	_, err = userService.Register(ctx, registerReq)
	if err == nil {
		t.Error("重复注册应该失败")
	}
	t.Logf("✅ 重复注册正确被拒绝: %v", err)

	// 7. 测试用户登录
	t.Log("🔐 测试用户登录...")
	loginReq := &user.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResp, err := userService.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("用户登录失败: %v", err)
	}
	
	if loginResp.Tokens.AccessToken == "" {
		t.Error("访问令牌不应该为空")
	}
	if loginResp.Tokens.RefreshToken == "" {
		t.Error("刷新令牌不应该为空")
	}
	t.Log("✅ 用户登录成功")

	// 8. 测试错误密码登录（应该失败）
	t.Log("🚫 测试错误密码登录...")
	wrongLoginReq := &user.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	_, err = userService.Login(ctx, wrongLoginReq)
	if err == nil {
		t.Error("错误密码登录应该失败")
	}
	t.Logf("✅ 错误密码登录正确被拒绝: %v", err)

	// 9. 测试JWT令牌验证
	t.Log("🔑 测试JWT令牌验证...")
	claims, err := jwtService.ValidateAccessToken(loginResp.Tokens.AccessToken)
	if err != nil {
		t.Fatalf("JWT令牌验证失败: %v", err)
	}
	
	if claims.UserID != registeredUser.ID {
		t.Errorf("期望用户ID为 %d，得到 %d", registeredUser.ID, claims.UserID)
	}
	if claims.Email != "test@example.com" {
		t.Errorf("期望邮箱为 'test@example.com'，得到 '%s'", claims.Email)
	}
	t.Logf("✅ JWT令牌验证成功，用户ID: %d", claims.UserID)

	// 10. 测试获取用户资料
	t.Log("📋 测试获取用户资料...")
	profile, err := userService.GetProfile(ctx, registeredUser.ID)
	if err != nil {
		t.Fatalf("获取用户资料失败: %v", err)
	}
	
	if profile.FirstName != "Test" {
		t.Errorf("期望名字为 'Test'，得到 '%s'", profile.FirstName)
	}
	if profile.LastName != "User" {
		t.Errorf("期望姓氏为 'User'，得到 '%s'", profile.LastName)
	}
	t.Logf("✅ 获取用户资料成功，全名: %s %s", profile.FirstName, profile.LastName)

	// 11. 测试更新用户资料
	t.Log("✏️  测试更新用户资料...")
	updateReq := &user.UpdateUserRequest{
		FirstName: "Updated",
		LastName:  "Name",
		TimeZone:  "Asia/Shanghai",
		Language:  "zh-CN",
	}

	updatedProfile, err := userService.UpdateProfile(ctx, registeredUser.ID, updateReq)
	if err != nil {
		t.Fatalf("更新用户资料失败: %v", err)
	}
	
	if updatedProfile.FirstName != "Updated" {
		t.Errorf("期望更新后名字为 'Updated'，得到 '%s'", updatedProfile.FirstName)
	}
	if updatedProfile.LastName != "Name" {
		t.Errorf("期望更新后姓氏为 'Name'，得到 '%s'", updatedProfile.LastName)
	}
	t.Logf("✅ 更新用户资料成功，新全名: %s %s", updatedProfile.FirstName, updatedProfile.LastName)

	// 12. 测试令牌刷新
	t.Log("🔄 测试令牌刷新...")
	newTokens, err := userService.RefreshToken(ctx, loginResp.Tokens.RefreshToken)
	if err != nil {
		t.Fatalf("令牌刷新失败: %v", err)
	}
	
	if newTokens.AccessToken == "" {
		t.Error("新访问令牌不应该为空")
	}
	if newTokens.RefreshToken == "" {
		t.Error("新刷新令牌不应该为空")
	}
	t.Log("✅ 令牌刷新成功")

	// 13. 测试修改密码
	t.Log("🔑 测试修改密码...")
	changePasswordReq := &user.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "newpassword456",
	}
	err = userService.ChangePassword(ctx, registeredUser.ID, changePasswordReq)
	if err != nil {
		t.Fatalf("修改密码失败: %v", err)
	}
	t.Log("✅ 修改密码成功")

	// 14. 测试使用新密码登录
	t.Log("🔐 测试使用新密码登录...")
	newLoginReq := &user.LoginRequest{
		Email:    "test@example.com",
		Password: "newpassword456",
	}
	_, err = userService.Login(ctx, newLoginReq)
	if err != nil {
		t.Fatalf("新密码登录失败: %v", err)
	}
	t.Log("✅ 新密码登录成功")

	// 15. 测试用户列表
	t.Log("📋 测试用户列表...")
	users, total, err := userService.ListUsers(ctx, 0, 10)
	if err != nil {
		t.Fatalf("获取用户列表失败: %v", err)
	}
	
	if total != 1 {
		t.Errorf("期望总用户数为 1，得到 %d", total)
	}
	if len(users) != 1 {
		t.Errorf("期望返回用户数为 1，得到 %d", len(users))
	}
	t.Logf("✅ 获取用户列表成功，总数: %d", total)

	// 16. 验证数据库文件创建
	t.Log("📁 验证数据库文件...")
	if _, err := os.Stat("./integration_test.db"); err != nil {
		t.Errorf("数据库文件不存在: %v", err)
	} else {
		t.Log("✅ 数据库文件已创建")
	}

	t.Log("🎉 所有集成测试通过！")
} 